package transcode

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"
)

var (
	wg      sync.WaitGroup
	allRes  = ""
	lastPer = -1
)

// ProcessVodFile ...
func ProcessVodFile(source string, data models.Vidinfo, cldata models.Video, prdata models.Pdata, conf utils.Config, ClientID string, userID uint) {
	utils.WLog("Starting VOD Processor..", ClientID)
	var (
		err error
		cmd string
	)

	// Path to source file
	sfpath := utils.Conf.SD + source

	// Checks if source file exists
	if source != "" {
		if _, err := os.Stat(sfpath); err == nil {
			utils.WLog("File found", ClientID)
		} else if os.IsNotExist(err) {
			utils.WLog("Error: file does not exist", ClientID)
			return
		} else {
			log.Println(err)
			utils.WLog("Error: file may or may not exist", ClientID)
			removeFile(utils.Conf.SD, source, ClientID)
			return
		}
	} else {
		removeFile(utils.Conf.SD, source, ClientID)
		return
	}

	// Full source file name
	fullsfname, err := filepath.EvalSymlinks(sfpath)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full file name", ClientID)
		removeFile(utils.Conf.SD, source, ClientID)
		return
	}

	// Source file name without extension
	sfnamewe := strings.Split(source, filepath.Ext(fullsfname))[0]

	// If transcoding directory does not exist creat it
	if _, err = os.Stat(utils.Conf.TD); os.IsNotExist(err) {
		os.Mkdir(utils.Conf.TD, 0777)
	}

	// File name after transcoding
	tempfile := fmt.Sprintf("%v%v.mp4", utils.Conf.TD, sfnamewe)

	// f
	destinationfile := fmt.Sprintf("%v%v.mp4", utils.Conf.DD, sfnamewe)

	// Checks if transcoded file with the same name already exists
	if _, err := os.Stat(tempfile); err == nil {
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already transcoding", sfnamewe+".mp4"), ClientID)
		removeFile(utils.Conf.SD, source, ClientID)
		return
	} else if _, err := os.Stat(destinationfile); err == nil {
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already exist in transcoded folder", sfnamewe+".mp4"), ClientID)
		removeFile(utils.Conf.SD, source, ClientID)
		return
	}

	utils.WLog(fmt.Sprintf("Starting to process %s", source), ClientID)

	// If data is empty get video info
	if data.IsEmpty() {
		data, err = GetVidInfo(utils.Conf.SD, source, utils.Conf.TempJson, utils.Conf.DataGen, utils.Conf.TempTxt, ClientID)
		if err != nil {
			log.Println(err)
			removeFile(utils.Conf.SD, source, ClientID)
			return
		}
	}

	// Generate thumbnails
	// utils.WLog("Generating thumbnail", ClientID)
	// wg.Add(1)
	// err = generateThumbnail(&wg, fullsfname, sfnamewe, data)
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
		cmd, tempdfs, err = generatePresetCmdLine(prdata, data, sfpath, fullsfname, fmt.Sprintf("%v%v", utils.Conf.TD, sfnamewe))
		tempfile = tempdfs[0]
		if err != nil {
			utils.WLog("Error: failed to generate cmd line", ClientID)
			log.Println(err)
			removeFile(utils.Conf.SD, source, ClientID)
			return
		}
	} else {
		cmd = generateClientCmdLine(cldata, data, sfpath, fullsfname, tempfile)
	}

	// check if client wants to save cmd line
	// if save {
	// 	err := db.AddCmdLine(source, cmd, tempdfs)
	// 	if err != nil {
	// 		utils.WLog("Error: failed to insert command line in database", ClientID)
	// 		log.Println(err)
	// 		removeFile(utils.Conf.SD, source, ClientID)
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
	go StartTranscode(source, utils.Conf, cmd, dfsl, ClientID, userID)
}

// StartTranscode ...
func StartTranscode(source string, conf utils.Config, cmdg string, dfsl string, ClientID string, userID uint) {
	var (
		err     error
		cmd     string
		dfsline string
		dfs     []string
		dur     string
		data    models.Vidinfo
	)

	utils.Conf = conf

	// Path to source file
	sfpath := utils.Conf.SD + source

	// Checks if source file exists
	if source != "" {
		if _, err := os.Stat(sfpath); err == nil {
			utils.WLog("File found", ClientID)
		} else if os.IsNotExist(err) {
			utils.WLog("Error: file does not exist", ClientID)
			return
		} else {
			log.Println(err)
			utils.WLog("Error: file may or may not exist", ClientID)
			removeFile(utils.Conf.SD, source, ClientID)
			return
		}
	} else {
		removeFile(utils.Conf.SD, source, ClientID)
		return
	}

	// Full source file name
	fullsfname, err := filepath.EvalSymlinks(sfpath)
	if err != nil {
		log.Println(err)
		utils.WLog("Error: failed to get full file name", ClientID)
		removeFile(utils.Conf.SD, source, ClientID)
		return
	}

	// Source file name without extension
	sfnamewe := strings.Split(source, filepath.Ext(fullsfname))[0]

	// If transcoding directory does not exist creat it
	if _, err = os.Stat(utils.Conf.TD); os.IsNotExist(err) {
		os.Mkdir(utils.Conf.TD, 0777)
	}

	// File name after transcoding
	tempfile := fmt.Sprintf("%v%v.mp4", utils.Conf.TD, sfnamewe)

	// f
	destinationfile := fmt.Sprintf("%v%v.mp4", utils.Conf.DD, sfnamewe)

	data, err = GetVidInfo(utils.Conf.SD, source, utils.Conf.TempJson, utils.Conf.DataGen, utils.Conf.TempTxt, ClientID)

	// ===============================================================

	if cmdg != "" {
		cmd = cmdg
		dfsline = dfsl
	} else {
		//cmd, dfsline, err = db.GetTranscodingInfo(source)
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

	// err = db.UpdateState(source, "Transcoding")
	// if err != nil {
	// 	utils.WLog("Error: failed to update state in database", ClientID)
	// 	log.Println(err)
	// 	removeFile(utils.Conf.SD, source, ClientID)
	// 	return
	// }

	wg.Add(1)
	err = runCmdCommand(cmd, dur, &wg, ClientID)
	wg.Wait()
	if err != nil {
		log.Println(err)
		utils.WLog("Error: could not start trancoding", ClientID)
		log.Printf("Error cmd line: %v", cmd)
		removeFile(utils.Conf.SD, source, ClientID)
		return
	} else if out, err := os.Stat(tempfile); os.IsNotExist(err) || out == nil {
		log.Println(err)
		utils.WLog("Error: transcoder failed", ClientID)
		log.Printf("Error cmd line: %v", cmd)
		removeFile(utils.Conf.SD, source, ClientID)
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
			removeFile(utils.Conf.SD, source, ClientID)
			for i := range tempdfs {
				os.Rename(utils.Conf.TD+dfs[i], utils.Conf.DD+dfs[i])
				nd, err := GetVidInfo(utils.Conf.DD, dfs[i], utils.Conf.TempJson, utils.Conf.DataGen, utils.Conf.TempTxt, ClientID)
				if err != nil {
					utils.WLog("Error: failed getting video data", ClientID)
					log.Println(err)
					removeStreamFiles(utils.Conf.DD, dfs, sfnamewe, ClientID)
					return
				}
				ndata = append(ndata, nd)
			}

			utils.InsertStream(ndata, dfs, "Transcoded", sfnamewe, userID)
			// err = db.InsertStream(ndata, dfs, "Transcoded", sfnamewe)
			// if err != nil {
			// 	utils.WLog("Error: failed to insert stream data in database", ClientID)
			// 	log.Println(err)
			// 	removeStreamFiles(utils.Conf.DD, dfs, sfnamewe, ClientID)
			// 	return
			// }

			msg := fmt.Sprintf("Transcoding coplete, stream name: %v", sfnamewe)
			utils.WLog(msg, ClientID)

		} else {
			removeFile(utils.Conf.SD, source, ClientID)
			os.Rename(tempfile, destinationfile)
			dfn := sfnamewe + ".mp4"
			ndata, err := GetVidInfo(utils.Conf.DD, dfn, utils.Conf.TempJson, utils.Conf.DataGen, utils.Conf.TempTxt, ClientID)
			if err != nil {
				utils.WLog("Error: failed getting video data", ClientID)
				log.Println(err)
				removeFile(utils.Conf.DD, dfn, ClientID)
				return
			}
			utils.InsertVideo(ndata, dfn, "Transcoded", userID, -1)

			// err = db.InsertVideo(ndata, dfn, "Transcoded", -1)
			// if err != nil {
			// 	utils.WLog("Error: failed to insert video data in database", ClientID)
			// 	log.Println(err)
			// 	removeFile(utils.Conf.DD, dfn, ClientID)
			// 	return
			// }

			msg := fmt.Sprintf("Transcoding coplete, file name: %v", filepath.Base(tempfile))
			utils.WLog(msg, ClientID)
		}
	}

}

func durToSec(dur string) (sec int) {
	durAry := strings.Split(dur, ":")
	if len(durAry) != 3 {
		return
	}
	hr, _ := strconv.Atoi(durAry[0])
	sec = hr * (60 * 60)
	min, _ := strconv.Atoi(durAry[1])
	sec += min * (60)
	second, _ := strconv.Atoi(durAry[2])
	sec += second
	return
}

func getRatio(res string, duration int, ClientID string) {
	i := strings.Index(res, "time=")
	if i >= 0 {
		time := res[i+5:]
		if len(time) > 8 {
			time = time[0:8]
			sec := durToSec(time)
			per := (sec * 100) / duration
			if lastPer != per {
				lastPer = per
				utils.UpdateLogMessage(fmt.Sprintf("Progress: %v %%", per), ClientID)
			}
			allRes = ""
		}
	}
}

func runCmdCommand(cmdl string, dur string, wg *sync.WaitGroup, ClientID string) error {
	defer wg.Done()

	if cmdl == "" {
		return errors.New("Error: cmd line is empty")
	}
	// Splits cmd command
	parts := strings.Fields(cmdl)
	head := parts[0]
	parts = parts[1:]

	cmd := exec.Command(head, parts...)

	// Creates pipe to listen to output
	stdout, err := cmd.StderrPipe()
	if err != nil {
		log.Println(err)
		return err
	}

	// Run commad
	if err := cmd.Start(); err != nil {
		log.Println(err)
		return err
	}
	oneByte := make([]byte, 8)

	// If duration is not provided dont sent progress bar
	if dur == "" {
		utils.WLog("Progress bar unavailable", ClientID)
	} else {
		duration := durToSec(dur)
		for {
			_, err := stdout.Read(oneByte)
			if err != nil {
				log.Println(err)
				break
			}
			allRes += string(oneByte)
			getRatio(allRes, duration, ClientID)
		}
	}

	return nil
}

func generateThumbnail(wg *sync.WaitGroup, source string, sourcewe string, data models.Vidinfo) error {
	defer wg.Done()

	//var timeStamp = (durToSec(data.Videotrack[0].Duration)) / utils.Conf.TNNum
	//var timeStamp = int((float64(durToSec(data.Videotrack[0].Duration)) * data.Videotrack[0].FrameRate) / 10)

	//var baseCmd = "ffmpeg -i %v -vf fps=1/%v %v%v%%03d.jpg"
	//var baseCmd = "ffmpeg -i %v -vf thumbnail=%v,setpts=N/TB -r 1 -vframes %v %v%v%%03d.jpg"

	var baseCmd = "ffmpeg -i %v -ss %v -vframes 1 %v%v.jpg"
	cmdl := fmt.Sprintf(baseCmd, source, utils.Conf.TNTS, utils.Conf.TND, sourcewe)

	parts := strings.Fields(cmdl)
	head := parts[0]
	parts = parts[1:]

	cmd := exec.Command(head, parts...)
	err := cmd.Run()

	return err
}

func removeFile(path string, filename string, ClientID string) {
	if _, err := os.Stat(path + filename); os.Remove(path+filename) != nil && !os.IsNotExist(err) {
		utils.WLog("Error: failed removing file", ClientID)
	}
	if err := utils.DeleteVideo(filename); err != nil {
		utils.WLog("Error: failed to remove file from database", ClientID)
	}

	return
}

func removeStreamFiles(path string, filenames []string, sname string, ClientID string) {
	for _, filename := range filenames {
		if os.Remove(path+filename) != nil {
			utils.WLog("Error: failed removing stream file(s)", ClientID)
		}
	}
	if err := utils.DeleteStream(sname); err != nil {
		utils.WLog("Error: failed to remove file from database", ClientID)
	}

	return
}
