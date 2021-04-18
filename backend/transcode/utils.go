package transcode

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"

	"github.com/batijo/video-platform/backend/models"
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
			}
			allRes = ""
		}
	}
}

func runCmdCommand(cmdl string, dur string, wg *sync.WaitGroup, ClientID uint) error {
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

func generateThumbnail(wg *sync.WaitGroup, fileName string, sourcewe string, data models.Vidinfo) error {
	defer wg.Done()

	//var timeStamp = (durToSec(data.Videotrack[0].Duration)) / utils.Conf.TNNum
	//var timeStamp = int((float64(durToSec(data.Videotrack[0].Duration)) * data.Videotrack[0].FrameRate) / 10)

	//var baseCmd = "ffmpeg -i %v -vf fps=1/%v %v%v%%03d.jpg"
	//var baseCmd = "ffmpeg -i %v -vf thumbnail=%v,setpts=N/TB -r 1 -vframes %v %v%v%%03d.jpg"

	var baseCmd = "ffmpeg -i %v -ss %v -vframes 1 %v%v.jpg"
	cmdl := fmt.Sprintf(baseCmd, fileName, utils.Conf.TNTS, utils.Conf.TND, sourcewe)

	parts := strings.Fields(cmdl)
	head := parts[0]
	parts = parts[1:]

	cmd := exec.Command(head, parts...)
	err := cmd.Run()

	return err
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
