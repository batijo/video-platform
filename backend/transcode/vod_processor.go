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
func ProcessVodFile(fileName string, data models.Vidinfo, clientData models.Video, presetData models.Pdata, ClientID string, vidId, userID uint) {
	utils.WLog("Starting VOD Processor..", ClientID)
	var (
		err error
		cmd string
	)

	// Path to source file
	sourceFileWithPath := utils.Conf.SD + fileName

	// Checks if source file exists
	if fileName != "" {
		if _, err := os.Stat(sourceFileWithPath); err == nil {
			utils.WLog("File found", ClientID)
		} else if os.IsNotExist(err) {
			utils.WLog("Error: file does not exist", ClientID)
			return
		} else {
			log.Println(err)
			utils.WLog("Error: file may or may not exist", ClientID)
			removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
			return
		}
	} else {
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	}

	// Full source file name
	fullSourceFilePathAndName, err := filepath.EvalSymlinks(sourceFileWithPath)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full file name", ClientID)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	}

	// source file name without extension
	sourceFileNameWithoutExt := strings.Split(fileName, filepath.Ext(fullSourceFilePathAndName))[0]

	// If transcoding directory does not exist creat it
	if _, err = os.Stat(utils.Conf.TD); os.IsNotExist(err) {
		os.Mkdir(utils.Conf.TD, 0777)
	}

	// File name after transcoding
	tempfile := fmt.Sprintf("%v%v.mp4", utils.Conf.TD, sourceFileNameWithoutExt)

	// f
	destinationFile := fmt.Sprintf("%v%v.mp4", utils.Conf.DD, sourceFileNameWithoutExt)

	// Checks if transcoded file with the same name already exists
	if _, err := os.Stat(tempfile); err == nil {
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already transcoding", sourceFileNameWithoutExt+".mp4"), ClientID)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	} else if _, err := os.Stat(destinationFile); err == nil {
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already exist in transcoded folder", sourceFileNameWithoutExt+".mp4"), ClientID)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	}

	utils.WLog(fmt.Sprintf("Starting to process %s", fileName), ClientID)

	// If data is empty get video info
	if data.IsEmpty() {
		data, err = GetVidInfo(utils.Conf.SD, fileName, ClientID, int(vidId))
		if err != nil {
			log.Println(err)
			removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
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

	msg := "%v video track(s), %v audio track(s) and %v subtitle(s) found"
	frmt := fmt.Sprintf(msg, data.Videotracks, data.Audiotracks, data.Subtitles)
	utils.WLog(frmt, ClientID)

	// Generate command line
	var tempdfs []string

	if utils.Conf.Presets {
		cmd, tempdfs, err = generatePresetCmdLine(presetData, data, sourceFileWithPath, fullSourceFilePathAndName,
			fmt.Sprintf("%v%v", utils.Conf.TD, sourceFileNameWithoutExt))
		tempfile = tempdfs[0]
		if err != nil {
			utils.WLog("Error: failed to generate cmd line", ClientID)
			log.Println(err)
			removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
			return
		}
	} else {
		cmd = generateClientCmdLine(clientData, data, sourceFileWithPath, fullSourceFilePathAndName, tempfile)
	}

	// check if client wants to save cmd line
	// if save {
	// 	err := db.AddCmdLine(fileName, cmd, tempdfs)
	// 	if err != nil {
	// 		utils.WLog("Error: failed to insert command line in database", ClientID)
	// 		log.Println(err)
	// 		removeVideo(utils.Conf.SD, fileName, ClientID)
	// 	} else {
	// 		utils.WLog("Transcoding parameters saved", ClientID)
	// 	}
	// } else {
	var dfsl string
	for i, d := range tempdfs {
		if i != len(tempdfs)-1 {
			dfsl += d + " "
		} else {
			dfsl += d
		}
	}
	go StartTranscode(fileName, cmd, dfsl, ClientID, vidId, userID)
}

// StartTranscode ...
func StartTranscode(fileName, cmdg, dfsl, ClientID string, vidId, userID uint) {
	var (
		err     error
		cmd     string
		dfsline string
		dfs     []string
		dur     string
		data    models.Vidinfo
	)

	// Path to source file
	sourceFileWithPath := utils.Conf.SD + fileName

	// Checks if source file exists
	if fileName != "" {
		if _, err := os.Stat(sourceFileWithPath); err == nil {
			utils.WLog("File found", ClientID)
		} else if os.IsNotExist(err) {
			utils.WLog("Error: file does not exist", ClientID)
			return
		} else {
			log.Println(err)
			utils.WLog("Error: file may or may not exist", ClientID)
			removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
			return
		}
	} else {
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	}

	// Full source file name
	fullSourceFilePathAndName, err := filepath.EvalSymlinks(sourceFileWithPath)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full file name", ClientID)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	}

	// source file name without extension
	sourceFileNameWithoutExt := strings.Split(fileName, filepath.Ext(fullSourceFilePathAndName))[0]

	// If transcoding directory does not exist creat it
	if _, err = os.Stat(utils.Conf.TD); os.IsNotExist(err) {
		os.Mkdir(utils.Conf.TD, 0777)
	}

	// File name after transcoding
	tempfile := fmt.Sprintf("%v%v.mp4", utils.Conf.TD, sourceFileNameWithoutExt)

	// f
	destinationFile := fmt.Sprintf("%v%v.mp4", utils.Conf.DD, sourceFileNameWithoutExt)
	destinationFile = utils.ReturnDifNameIfDublicate(destinationFile, "")

	data, err = GetVidInfo(utils.Conf.SD, fileName, ClientID, int(vidId))

	// ===============================================================

	if cmdg != "" {
		cmd = cmdg
		dfsline = dfsl
	} else {
		//cmd, dfsline, err = db.GetTranscodingInfo(fileName)
	}

	tempdfs := strings.Split(dfsline, " ")

	if dfsline != "" {
		tempfile = tempdfs[0]

		// removes path from stream files names
		for _, d := range tempdfs {
			df := strings.SplitAfterN(d, "/", 3)[2]
			dfs = append(dfs, df)
		}
	}

	// Run generated command line
	utils.WLog("Starting to transcode", ClientID)
	if utils.Conf.Debug {
		dur = utils.Conf.DebugEnd
	} else {
		dur = data.Videotrack[0].Duration
	}

	video := models.Video{Model: gorm.Model{ID: vidId}}
	err = utils.DB.Model(&video).Update("state", "transcoding").Error
	if err != nil {
		utils.WLog("Error: failed to update state in database", ClientID)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	}

	wg.Add(1)
	err = runCmdCommand(cmd, dur, &wg, ClientID)
	wg.Wait()
	if err != nil {
		log.Println(err)
		utils.WLog("Error: could not start trancoding", ClientID)
		log.Printf("Error cmd line: %v", cmd)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
		return
	} else if out, err := os.Stat(tempfile); os.IsNotExist(err) || out == nil {
		log.Println(err)
		utils.WLog("Error: transcoder failed", ClientID)
		log.Printf("Error cmd line: %v", cmd)
		removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
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
			removeVideo(utils.Conf.SD+fileName, int(vidId), ClientID)
			for i := range tempdfs {

				if newName := utils.ReturnDifNameIfDublicate(dfs[i], utils.Conf.DD); newName != fileName {
					if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.TD+newName); err != nil {
						utils.WLog("Error: failed while moving files", ClientID)
						log.Println(err)
						removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
						return
					}
					dfs[i] = newName
				}

				if err := utils.MoveFile(utils.Conf.TD+dfs[i], utils.Conf.DD+dfs[i]); err != nil {
					utils.WLog("Error: failed while moving files", ClientID)
					log.Println(err)
					removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
					return
				}

				nd, err := GetVidInfo(utils.Conf.DD, dfs[i], ClientID, int(vidId))
				if err != nil {
					utils.WLog("Error: failed getting video data", ClientID)
					log.Println(err)
					removeStreamVideos(utils.Conf.DD, dfs, sourceFileNameWithoutExt, ClientID)
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

			msg := fmt.Sprintf("Transcoding coplete, stream name: %v", sourceFileNameWithoutExt)
			utils.WLog(msg, ClientID)

		} else {

			dfn := sourceFileNameWithoutExt + ".mp4"
			if newName := utils.ReturnDifNameIfDublicate(dfn, utils.Conf.DD); newName != fileName {
				if err := utils.MoveFile(utils.Conf.TD+dfn, utils.Conf.TD+newName); err != nil {
					utils.WLog("Error: failed while moving file", ClientID)
					log.Println(err)
					removeVideo(utils.Conf.DD+dfn, int(vidId), ClientID)
					return
				}
				dfn = newName
			}

			removeVideo(utils.Conf.SD+fileName, -1, ClientID)
			if err := utils.MoveFile(tempfile, destinationFile); err != nil {
				log.Println(err)
				removeVideo(utils.Conf.TD+dfn, int(vidId), ClientID)
				return
			}

			ndata, err := GetVidInfo(utils.Conf.DD, dfn, ClientID, int(vidId))
			if err != nil {
				utils.WLog("Error: failed getting video data", ClientID)
				log.Println(err)
				removeVideo(utils.Conf.DD+dfn, int(vidId), ClientID)
				return
			}
			_, err = utils.InsertVideo(ndata, dfn, "transcoded", userID, -1)
			if err != nil {
				utils.WLog("Error: failed to insert video data in database", ClientID)
				log.Println(err)
				removeVideo(utils.Conf.DD+dfn, int(vidId), ClientID)
				return
			}

			msg := fmt.Sprintf("Transcoding coplete, file name: %v", filepath.Base(tempfile))
			utils.WLog(msg, ClientID)
		}
	}

}
