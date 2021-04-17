package transcode

import (
	"encoding/json"
	"errors"
	"os/exec"
	"strings"
	"sync"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
)

func getMediaInfoJSON(source string, wg *sync.WaitGroup) ([]byte, error) {
	defer wg.Done()

	if source == "" {
		return []byte(""), errors.New("no video file provided")
	}

	cmd := "ffprobe -v quiet -print_format json -show_streams -show_format "

	//Splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd + source)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()

	if strings.Replace(string(out), "\n", "", -1) == `{}` {
		return out, errors.New("json data is empty")
	} else if err != nil {
		return out, err
	}

	//checks if data is json file
	if !json.Valid(out) {
		return out, errors.New("data is not valid json file")
	}

	return out, nil
}

// GetVidInfo retruns struct with information about video file
func GetVidInfo(path string, filename string, ClientID string, vidId int) (models.Vidinfo, error) {
	var (
		wg sync.WaitGroup
		vi models.Vidinfo
	)

	// Geting data about video
	wg.Add(1)
	infob, err := getMediaInfoJSON(path+filename, &wg)
	if err != nil {
		utils.WLog("Error: could not get json data from file", ClientID)
		removeVideo(path+filename, vidId, ClientID)
		return vi, err
	}
	wg.Wait()

	// Unmarshal data into Ffprobe struct
	var metadata models.Ffprobe
	err = json.Unmarshal(infob, &metadata)
	if err != nil {
		utils.WLog("Error: failed to unmarshal json file", ClientID)
		removeVideo(path+filename, vidId, ClientID)
		return vi, err
	}

	vi.ParseFFprobeData(metadata, filename)

	return vi, nil
}
