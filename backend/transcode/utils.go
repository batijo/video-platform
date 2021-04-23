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

	"github.com/batijo/video-platform/backend/utils"
)

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

func getRatio(res string, duration int, ClientID uint) {
	i := strings.Index(res, "time=")
	if i >= 0 {
		time := res[i+5:]
		if len(time) > 8 {
			time = time[0:8]
			sec := durToSec(time)
			per := (sec * 100) / duration
			if lastPer != per {
				lastPer = per
				utils.UpdateUserMessage(fmt.Sprintf("Progress: %v %%", per), ClientID)
				utils.UpdateAllUsersMessage(fmt.Sprintf("Progress: %v %%", per))
				utils.DebugLog(fmt.Sprintf("==%v%%", per))
			}
			allRes = ""
		}
	}
}

func runCmdCommand(cmdl string, dur int, wg *sync.WaitGroup, ClientID uint) error {
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
	if dur == 0 {
		utils.WLog("Progress bar unavailable", ClientID)
	} else {
		duration := dur
		for {
			_, err := stdout.Read(oneByte)
			if err != nil {
				if err.Error() == "EOF" {
					break
				}
				log.Println(err)
				break
			}
			allRes += string(oneByte)
			getRatio(allRes, duration, ClientID)
		}
	}

	return nil
}

func removeVideo(path string, vidId int, ClientID uint) {
	if _, err := os.Stat(path); os.Remove(path) != nil && !os.IsNotExist(err) {
		utils.WLog("Error: failed removing file", ClientID)
	}
	if vidId >= 0 {
		if err := utils.DeleteVideo(uint(vidId)); err != nil {
			utils.WLog("Error: failed to remove file from database", ClientID)
		}
	}

	return
}

func removeStreamVideos(path string, filenames []string, sname string, ClientID uint) {
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

func fileNameWithoutExt(filename, nameWithPath string) (string, error) {
	split := strings.Split(filename, filepath.Ext(nameWithPath))
	if len(split) < 1 {
		return "", errors.New("split file name is empty array")
	}
	return split[0], nil
}
