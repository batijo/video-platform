package transcode

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
)

var (
	allRes  = ""
	lastPer = -1
)

// ProcessVodFile is the shitiest function I've ever writen. Hate it
func processVodFile(ED models.Encodedata) {
	utils.PrintStruct(ED, "EncodeData model")
	var (
		wg                        sync.WaitGroup
		err                       error
		data                      models.Vidinfo
		cmd                       string
		dfs                       []string
		dur                       int
		fullSourceFilePathAndName string
		sourceFileNameWithoutExt  string
		destinationFile           string
		tempfile                  string
		// All temporary destination files of preset
		tempdfs []string
		// All temporary destination files of preset put in one string
		dfsline      string
		video        = models.Video{ID: ED.Video.ID}
		newVideoData models.Video
		path         = map[string]string{
			"not_transcoded": utils.Conf.SD,
			"transcoding":    utils.Conf.TD,
			"transcoded":     utils.Conf.DD,
		}
		clientID = ED.Video.UserID
		vidId    = int(ED.Video.ID)
	)
	utils.WLog("Starting VOD Processor..", clientID)

	// Check if video is already in transcoding
	if ED.Video.State == "transcoding" {
		log.Println(err)
		utils.WLog("Error: video already in transcoding", clientID)
		finished <- false
		return
	} else if ED.Video.State != "not_transcoded" {
		if ED.Video.State != "transcoded" {
			log.Println(err)
			utils.WLog(fmt.Sprintf("Error: unkown video state: '%v'", ED.Video.State), clientID)
			finished <- false
			return
		} else {
			oldName := ED.Video.FileName
			ED.Video.FileName = utils.ReturnDifNameIfDublicate(ED.Video.FileName, utils.Conf.SD)
			if err := utils.MoveFile(
				utils.Conf.DD+oldName,
				utils.Conf.SD+ED.Video.FileName,
			); err != nil {

				log.Println(err)
				utils.WLog("Error: moving file to source directory", clientID)
				finished <- false
				return
			} else {
				ED.Video.State = "not_transcoded"
			}
		}
	}

	// Gather video data
	data, err = GetVidInfo(path[ED.Video.State], ED.Video.FileName, clientID, vidId)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full information about video", clientID)
		finished <- false
		return
	}
	utils.PrintStruct(data, "Source video data")

	// Checks if source file exists
	if data.FileName != "" {
		if _, err := os.Stat(path[ED.Video.State] + data.FileName); err == nil {
			utils.WLog("File found", clientID)
		} else if os.IsNotExist(err) {
			utils.WLog("Error: file does not exist", clientID)
			finished <- false
			return
		} else {
			log.Println(err)
			utils.WLog("Error: file may or may not exist", clientID)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
	} else {
		utils.WLog("Error: Missing file name in data", clientID)
		finished <- false
		return
	}

	// Full source file name
	fullSourceFilePathAndName, err = filepath.EvalSymlinks(path[ED.Video.State] + data.FileName)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full file path", clientID)
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	}
	utils.DebugLog("fullSourceFilePathAndName: " + fullSourceFilePathAndName)

	// Source file name without extension
	sourceFileNameWithoutExt, err = fileNameWithoutExt(data.FileName, fullSourceFilePathAndName)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get filename without extention", clientID)
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	}
	utils.DebugLog("sourceFileNameWithoutExt: " + sourceFileNameWithoutExt)

	// If transcoding directory does not exist creat it
	if _, err = os.Stat(utils.Conf.TD); os.IsNotExist(err) {
		err = os.Mkdir(utils.Conf.TD, 0777)
		if err != nil {
			log.Println(err)
			utils.WLog("Error: failed to create transcoding directory", clientID)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
	}

	if _, err = os.Stat(utils.Conf.DD); os.IsNotExist(err) {
		err = os.Mkdir(utils.Conf.DD, 0777)
		if err != nil {
			log.Println(err)
			utils.WLog("Error: failed to create destination directory", clientID)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
	}

	// File name after transcoding
	tempfile = fmt.Sprintf("%v%v.mp4", utils.Conf.TD, sourceFileNameWithoutExt)
	tempfile = utils.ReturnDifNameIfDublicate(tempfile, "")
	utils.DebugLog("tempfile: " + tempfile)

	// Create new name for file
	destinationFile = fmt.Sprintf("%v%v.mp4", utils.Conf.DD, sourceFileNameWithoutExt)
	destinationFile = utils.ReturnDifNameIfDublicate(destinationFile, "")
	utils.DebugLog("destinationFile: " + destinationFile)

	// Checks if transcoded file with the same name already exists
	if _, err := os.Stat(tempfile); err == nil {
		utils.WLog(
			fmt.Sprintf("Error: file \"%v\" already transcoding", sourceFileNameWithoutExt+".mp4"),
			clientID,
		)
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	} else if _, err := os.Stat(destinationFile); err == nil {
		utils.WLog(
			fmt.Sprintf(
				"Error: file \"%v\" already exist in transcoded folder", sourceFileNameWithoutExt+".mp4"),
			clientID,
		)
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	}

	utils.WLog(fmt.Sprintf("Starting to process %s", data.FileName), clientID)

	// Generate one video if only one preset is given
	if len(ED.Presets) == 1 {
		utils.DebugLog("Creating nonpreset cmd line of preset")
		var (
			vidEndoceData models.Video
			vtPreset      models.Preset
			atPreset      models.Preset
		)
		vtPreset, err = utils.GetPreset(ED.Presets[0].VidPreset)
		if err != nil {
			log.Println(err)
			utils.WLog("Error: failed to get preset", clientID)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
		atPreset, err = utils.GetPreset(ED.Presets[0].AudPreset)
		if err != nil {
			log.Println(err)
			utils.WLog("Error: failed to get preset", clientID)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
		vidEndoceData.ParseWithPreset(
			vtPreset,
			atPreset,
			ED.Video.FrameRate,
			ED.Presets[0].VtId,
			ED.Presets[0].AudioT[0].StreamID,
			ED.Presets[0].AudioT[0].Language,
		)
		cmd = generateClientCmdLine(
			vidEndoceData,
			data,
			path[ED.Video.State]+data.FileName,
			fullSourceFilePathAndName,
			tempfile,
		)
		// Generate CMD based on presets
	} else if len(ED.Presets) > 1 {
		utils.DebugLog("Creating preset cmd line")
		cmd, dfs, err = generatePresetCmdLine(
			ED.Presets,
			data,
			path[ED.Video.State]+data.FileName,
			fullSourceFilePathAndName,
			//fmt.Sprintf("%v%v", utils.Conf.TD, sourceFileNameWithoutExt),
			sourceFileNameWithoutExt,
		)
		if err != nil || len(dfs) < 2 {
			utils.WLog("Error: failed to generate cmd line", clientID)
			log.Println(err)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
		for _, d := range dfs {
			tempdfs = append(tempdfs, utils.Conf.TD+d)
		}
		tempfile = tempdfs[0]
		for i, d := range tempdfs {
			if i != len(tempdfs)-1 {
				dfsline += d + " "
			} else {
				dfsline += d
			}
		}

	} else {
		utils.DebugLog("Creating nonpreset cmd line")
		var vidEndoceData models.Video
		vidEndoceData.ParseWithEncode(ED.EncData, ED.Video.State)
		cmd = generateClientCmdLine(
			vidEndoceData,
			data,
			path[ED.Video.State]+data.FileName,
			fullSourceFilePathAndName,
			tempfile)
	}
	utils.DebugLog("CMD line: " + cmd)

	// Run generated command line
	utils.WLog("Starting to transcode", clientID)
	if utils.Conf.Debug {
		dur = durToSec(utils.Conf.DebugEnd)
	} else {
		dur = int(data.Videotrack[0].Duration)
	}
	utils.DebugLog("dur: " + fmt.Sprint(dur))

	err = utils.DB.Model(&video).Update("state", "transcoding").Error
	if err != nil {
		utils.WLog("Error: failed to update state in database", clientID)
		log.Println(err)
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	}

	wg.Add(1)
	err = runCmdCommand(cmd, dur, &wg, clientID)
	wg.Wait()

	if err != nil {
		log.Println(err)
		utils.WLog("Error: could not start trancoding", clientID)
		log.Printf("Error cmd line: %v", cmd)
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	} else if out, err := os.Stat(tempfile); os.IsNotExist(err) || out == nil {
		log.Println(err)
		utils.WLog("Error: transcoder failed", clientID)
		log.Printf("Error cmd line: %v", cmd)
		log.Println("==================================FFMPEG_OUTPUT====================================")
		log.Println(allRes)
		log.Println("==================================FFMPEG_OUTPUT====================================")
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	} else {

		// Removes source file and moves transcoded file to /videos/transcoded
		if len(ED.Presets) > 1 {
			var (
				ndata []models.Vidinfo
			)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			for i := range tempdfs {

				// if newName := utils.ReturnDifNameIfDublicate(dfs[i], utils.Conf.DD); newName != data.FileName {
				// 	if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.TD+newName); err != nil {
				// 		utils.WLog("Error: failed while moving files", clientID)
				// 		log.Println(err)
				// 		removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, clientID)
				// 		finished <- false
				// 		return
				// 	}
				// 	dfs[i] = newName
				// }

				if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.DD+dfs[i]); err != nil {
					utils.WLog("Error: failed while moving files", clientID)
					log.Println(err)
					removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, clientID)
					finished <- false
					return
				}

				nd, err := GetVidInfo(utils.Conf.DD, dfs[i], clientID, vidId)
				if err != nil {
					utils.WLog("Error: failed getting video data", clientID)
					log.Println(err)
					removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, clientID)
					finished <- false
					return
				}
				ndata = append(ndata, nd)
			}
			utils.PrintStruct(ndata, "Preset videos data")

			if err = utils.InsertStream(ndata, dfs, "transcoded", sourceFileNameWithoutExt, clientID, ED.Video.Public); err != nil {
				utils.WLog("Error: failed to insert stream data in database", clientID)
				log.Println(err)
				removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, clientID)
				return
			}

			utils.WLog(
				fmt.Sprintf(
					"Transcoding complete, stream name: %v",
					sourceFileNameWithoutExt),
				clientID)

		} else {

			dfn := sourceFileNameWithoutExt + ".mp4"
			if newName := utils.ReturnDifNameIfDublicate(dfn, utils.Conf.DD); newName != data.FileName {
				if err := utils.MoveFile(utils.Conf.TD+dfn, utils.Conf.TD+newName); err != nil {
					utils.WLog("Error: failed while moving file", clientID)
					log.Println(err)
					removeVideo(utils.Conf.DD+dfn, vidId, clientID)
					finished <- false
					return
				}
				dfn = newName
			}

			removeVideo(path[ED.Video.State]+data.FileName, -1, clientID)
			if err := utils.MoveFile(utils.Conf.TD+dfn, destinationFile); err != nil {
				log.Println(err)
				removeVideo(utils.Conf.TD+dfn, vidId, clientID)
				finished <- false
				return
			}

			newData, err := GetVidInfo(utils.Conf.DD, dfn, clientID, vidId)
			if err != nil {
				utils.WLog("Error: failed getting video data", clientID)
				log.Println(err)
				removeVideo(utils.Conf.DD+dfn, vidId, clientID)
				finished <- false
				return
			}
			utils.PrintStruct(newData, "New video data")

			newVideoData.ParseWithVidinfo(newData)
			newVideoData.State = "transcoded"
			utils.UpdateVideo(uint(vidId), newVideoData)
			if err != nil {
				utils.WLog("Error: failed to update video data in database", clientID)
				log.Println(err)
				removeVideo(utils.Conf.DD+dfn, vidId, clientID)
				finished <- false
				return
			}

			utils.WLog(
				fmt.Sprintf(
					"Transcoding complete, file name: %v",
					filepath.Base(tempfile)),
				clientID)
		}
	}
	finished <- true
}
