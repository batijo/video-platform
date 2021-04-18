package transcode

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
)

var (
	allRes  = ""
	lastPer = -1
)

// ProcessVodFile ...
func processVodFile(ED models.EncodeData) {
	var (
		wg                        sync.WaitGroup
		err                       error
		data                      models.Vidinfo
		cmd                       string
		dfs                       []string
		dur                       string
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
			"transcoded":     utils.Conf.DD,
		}
		clientID = ED.Video.UserID
		vidId    = int(ED.Video.ID)
	)
	utils.WLog("Starting VOD Processor..", clientID)

	// Gather video data
	data, err = GetVidInfo(path[ED.Video.State], ED.Video.FileName, clientID, vidId)
	if err != nil {
		log.Println(err)
		finished <- false
		return
	}

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

	// Source file name without extension
	sourceFileNameWithoutExt = strings.Split(data.FileName, filepath.Ext(fullSourceFilePathAndName))[0]

	// If transcoding directory does not exist creat it
	if _, err = os.Stat(utils.Conf.TD); os.IsNotExist(err) {
		os.Mkdir(utils.Conf.TD, 0777)
	}

	// File name after transcoding
	tempfile = fmt.Sprintf("%v%v.mp4", utils.Conf.TD, sourceFileNameWithoutExt)

	// f
	destinationFile = fmt.Sprintf("%v%v.mp4", utils.Conf.DD, sourceFileNameWithoutExt)
	destinationFile = utils.ReturnDifNameIfDublicate(destinationFile, "")

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

	// If data is empty get video info
	if data.IsEmpty() {
		data, err = GetVidInfo(path[ED.Video.State], data.FileName, clientID, vidId)
		if err != nil {
			log.Println(err)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
	}

	utils.WLog(
		fmt.Sprintf(
			"%v video track(s), %v audio track(s) and %v subtitle(s) found",
			data.Videotracks,
			data.Audiotracks,
			data.Subtitles),
		clientID)

	// Generate command line
	if len(ED.Presets) > 0 {
		cmd, tempdfs, err = generatePresetCmdLine(
			ED.Presets,
			data,
			path[ED.Video.State]+data.FileName,
			fullSourceFilePathAndName,
			fmt.Sprintf("%v%v", utils.Conf.TD, sourceFileNameWithoutExt))

		tempfile = tempdfs[0]
		if err != nil {
			utils.WLog("Error: failed to generate cmd line", clientID)
			log.Println(err)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			finished <- false
			return
		}
	} else {
		var vidEndoceData models.Video
		vidEndoceData.ParseWithEncode(ED.EncData, ED.Video.State)
		cmd = generateClientCmdLine(
			vidEndoceData,
			data,
			path[ED.Video.State]+data.FileName,
			fullSourceFilePathAndName,
			tempfile)
	}

	for i, d := range tempdfs {
		if i != len(tempdfs)-1 {
			dfsline += d + " "
		} else {
			dfsline += d
		}
	}

	// Removes path from stream files names
	for _, d := range tempdfs {
		df := strings.SplitAfterN(d, "/", 3)[2]
		dfs = append(dfs, df)
	}

	// Run generated command line

	utils.WLog("Starting to transcode", clientID)
	if utils.Conf.Debug {
		dur = utils.Conf.DebugEnd
	} else {
		dur = data.Videotrack[0].Duration
	}

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
		removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
		finished <- false
		return
	} else {

		if _, err = os.Stat(utils.Conf.DD); os.IsNotExist(err) {
			os.Mkdir(utils.Conf.DD, 0777)
		}
		// Removes source file and moves transcoded file to /videos/transcoded
		if len(ED.Presets) > 0 {
			var (
				ndata []models.Vidinfo
			)
			removeVideo(path[ED.Video.State]+data.FileName, vidId, clientID)
			for i := range tempdfs {

				if newName := utils.ReturnDifNameIfDublicate(dfs[i], utils.Conf.DD); newName != data.FileName {
					if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.TD+newName); err != nil {
						utils.WLog("Error: failed while moving files", clientID)
						log.Println(err)
						removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, clientID)
						finished <- false
						return
					}
					dfs[i] = newName
				}

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

			utils.InsertStream(ndata, dfs, "transcoded", sourceFileNameWithoutExt, clientID)
			// err = db.InsertStream(ndata, dfs, "Transcoded", sourceFileNameWithoutExt)
			// if err != nil {
			// 	utils.WLog("Error: failed to insert stream data in database", clientID)
			// 	log.Println(err)
			// 	removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, clientID)
			// 	return
			// }

			utils.WLog(
				fmt.Sprintf(
					"Transcoding coplete, stream name: %v",
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
					"Transcoding coplete, file name: %v",
					filepath.Base(tempfile)),
				clientID)
		}
	}
	finished <- true
}
