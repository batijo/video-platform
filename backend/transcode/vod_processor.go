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
	"github.com/jinzhu/gorm"
)

var (
	wg      sync.WaitGroup
	allRes  = ""
	lastPer = -1
)

// ProcessVodFile ...
func ProcessVodFile(clientData models.Video, presetData models.Pdata, ClientID string, vidId, userID uint) {
	utils.WLog("Starting VOD Processor..", ClientID)
	var (
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
		dfsline string
		video   = models.Video{Model: gorm.Model{ID: vidId}}
	)

	// Checks if source file exists
	if data.FileName != "" {
		if _, err := os.Stat(utils.Conf.SD + data.FileName); err == nil {
			utils.WLog("File found", ClientID)
		} else if os.IsNotExist(err) {
			utils.WLog("Error: file does not exist", ClientID)
			finished <- false
			return
		} else {
			log.Println(err)
			utils.WLog("Error: file may or may not exist", ClientID)
			removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
			finished <- false
			return
		}
	} else {
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
		finished <- false
		return
	}

	// Gather video data
	data, err = GetVidInfo(utils.Conf.SD, clientData.FileName, ClientID, int(vidId))
	if err != nil {
		log.Println(err)
		finished <- false
		return
	}

	// Full source file name
	fullSourceFilePathAndName, err = filepath.EvalSymlinks(utils.Conf.SD + data.FileName)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full file path", ClientID)
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
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
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already transcoding", sourceFileNameWithoutExt+".mp4"), ClientID)
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
		finished <- false
		return
	} else if _, err := os.Stat(destinationFile); err == nil {
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already exist in transcoded folder", sourceFileNameWithoutExt+".mp4"), ClientID)
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
		finished <- false
		return
	}

	utils.WLog(fmt.Sprintf("Starting to process %s", data.FileName), ClientID)

	// If data is empty get video info
	if data.IsEmpty() {
		data, err = GetVidInfo(utils.Conf.SD, data.FileName, ClientID, int(vidId))
		if err != nil {
			log.Println(err)
			removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
			finished <- false
			return
		}
	}

	// Generate thumbnails
	// utils.WLog("Generating thumbnail", ClientID)
	// wg.Add(1)
	// err = generateThumbnail(&wg, fullSourceFilePathAndName, sourceFileNameWithoutExt, data)
	// if err != nil {
	// 	log.Printf("Generate thumbnail exited with error: %v", err)
	// }
	// wg.Wait()

	utils.WLog(
		fmt.Sprintf(
			"%v video track(s), %v audio track(s) and %v subtitle(s) found",
			data.Videotracks,
			data.Audiotracks,
			data.Subtitles),
		ClientID)

	// Generate command line
	if utils.Conf.Presets {
		cmd, tempdfs, err = generatePresetCmdLine(
			presetData,
			data,
			utils.Conf.SD+data.FileName,
			fullSourceFilePathAndName,
			fmt.Sprintf("%v%v", utils.Conf.TD, sourceFileNameWithoutExt))

		tempfile = tempdfs[0]
		if err != nil {
			utils.WLog("Error: failed to generate cmd line", ClientID)
			log.Println(err)
			removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
			finished <- false
			return
		}
	} else {
		cmd = generateClientCmdLine(
			clientData,
			data,
			utils.Conf.SD+data.FileName,
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

	utils.WLog("Starting to transcode", ClientID)
	if utils.Conf.Debug {
		dur = utils.Conf.DebugEnd
	} else {
		dur = data.Videotrack[0].Duration
	}

	err = utils.DB.Model(&video).Update("state", "transcoding").Error
	if err != nil {
		utils.WLog("Error: failed to update state in database", ClientID)
		log.Println(err)
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
		finished <- false
		return
	}

	wg.Add(1)
	err = runCmdCommand(cmd, dur, &wg, ClientID)
	wg.Wait()

	if err != nil {
		log.Println(err)
		utils.WLog("Error: could not start trancoding", ClientID)
		log.Printf("Error cmd line: %v", cmd)
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
		finished <- false
		return
	} else if out, err := os.Stat(tempfile); os.IsNotExist(err) || out == nil {
		log.Println(err)
		utils.WLog("Error: transcoder failed", ClientID)
		log.Printf("Error cmd line: %v", cmd)
		removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
		finished <- false
		return
	} else {

		if _, err = os.Stat(utils.Conf.DD); os.IsNotExist(err) {
			os.Mkdir(utils.Conf.DD, 0777)
		}
		// Removes source file and moves transcoded file to /videos/transcoded
		if utils.Conf.Presets {
			var (
				ndata []models.Vidinfo
			)
			removeVideo(utils.Conf.SD+data.FileName, int(vidId), ClientID)
			for i := range tempdfs {

				if newName := utils.ReturnDifNameIfDublicate(dfs[i], utils.Conf.DD); newName != data.FileName {
					if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.TD+newName); err != nil {
						utils.WLog("Error: failed while moving files", ClientID)
						log.Println(err)
						removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
						finished <- false
						return
					}
					dfs[i] = newName
				}

				if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.DD+dfs[i]); err != nil {
					utils.WLog("Error: failed while moving files", ClientID)
					log.Println(err)
					removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
					finished <- false
					return
				}

				nd, err := GetVidInfo(utils.Conf.DD, dfs[i], ClientID, int(vidId))
				if err != nil {
					utils.WLog("Error: failed getting video data", ClientID)
					log.Println(err)
					removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
					finished <- false
					return
				}
				ndata = append(ndata, nd)
			}

			utils.InsertStream(ndata, dfs, "transcoded", sourceFileNameWithoutExt, userID)
			// err = db.InsertStream(ndata, dfs, "Transcoded", sourceFileNameWithoutExt)
			// if err != nil {
			// 	utils.WLog("Error: failed to insert stream data in database", ClientID)
			// 	log.Println(err)
			// 	removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
			// 	return
			// }

			utils.WLog(
				fmt.Sprintf(
					"Transcoding coplete, stream name: %v",
					sourceFileNameWithoutExt),
				ClientID)

		} else {

			dfn := sourceFileNameWithoutExt + ".mp4"
			if newName := utils.ReturnDifNameIfDublicate(dfn, utils.Conf.DD); newName != data.FileName {
				if err := utils.MoveFile(utils.Conf.TD+dfn, utils.Conf.TD+newName); err != nil {
					utils.WLog("Error: failed while moving file", ClientID)
					log.Println(err)
					removeVideo(utils.Conf.DD+dfn, int(vidId), ClientID)
					finished <- false
					return
				}
				dfn = newName
			}

			removeVideo(utils.Conf.SD+data.FileName, -1, ClientID)
			if err := utils.MoveFile(tempfile, destinationFile); err != nil {
				log.Println(err)
				removeVideo(utils.Conf.TD+dfn, int(vidId), ClientID)
				finished <- false
				return
			}

			ndata, err := GetVidInfo(utils.Conf.DD, dfn, ClientID, int(vidId))
			if err != nil {
				utils.WLog("Error: failed getting video data", ClientID)
				log.Println(err)
				removeVideo(utils.Conf.DD+dfn, int(vidId), ClientID)
				finished <- false
				return
			}
			_, err = utils.InsertVideo(ndata, "transcoded", userID, -1)
			if err != nil {
				utils.WLog("Error: failed to insert video data in database", ClientID)
				log.Println(err)
				removeVideo(utils.Conf.DD+dfn, int(vidId), ClientID)
				finished <- false
				return
			}

			utils.WLog(
				fmt.Sprintf(
					"Transcoding coplete, file name: %v",
					filepath.Base(tempfile)),
				ClientID)
		}
	}

	finished <- true
}
