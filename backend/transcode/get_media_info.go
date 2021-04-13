package transcode

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
)

func getMediaInfoJSON(source string, wg *sync.WaitGroup) ([]byte, error) {
	defer wg.Done()

	cmd := "ffprobe -v quiet -print_format json -show_streams -show_format"
	cmd += " " + source

	//Splitting head => g++ parts => rest of the command
	parts := strings.Fields(cmd)
	head := parts[0]
	parts = parts[1:]

	out, err := exec.Command(head, parts...).Output()

	ej := `{}`
	if string(out) == ej {
		return out, errors.New("json data is empty")
	} else if err != nil {
		return out, err
	}

	//checks if data is json file
	if !json.Valid([]byte(out)) {
		return out, errors.New("data is not valid json file")
	}

	return out, nil
}

func generateDataFile(wg *sync.WaitGroup, gpath string, prefix string) error {
	defer wg.Done()

	parts := []string{gpath, prefix}
	out, err := exec.Command("python3", parts...).Output()
	if err != nil {
		log.Println(err)
		log.Println("boi")
		out, err = exec.Command("python", parts...).Output()
		if err != nil {
			log.Println(err)
			if err := os.Remove(fmt.Sprintf(utils.Conf.SourceJson, prefix)); err != nil {
				log.Println(err)
			}
			return err
		}
	}

	if string(out) == "False\n" {
		return errors.New("generate_data.py output False")
	} else if string(out) != "True\n" {
		return fmt.Errorf("generate_data.py uknown output: %v", string(out))
	}

	return nil
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
	//var raw map[string]interface{}
	var metadata models.Ffprobe
	err = json.Unmarshal(infob, &metadata)
	//info, err := json.Marshal(raw)
	if err != nil {
		utils.WLog("Error: failed to unmarshal json file", ClientID)
		removeVideo(path+filename, vidId, ClientID)
		return vi, err
	}

	vi.ParseFFprobeData(metadata, filename)

	// // random number
	// rand.Seed(time.Now().UnixNano())
	// prefix := fmt.Sprint(rand.Intn(10000))
	// tempSoureFileName := fmt.Sprintf(utils.Conf.SourceJson, prefix)

	// err = ioutil.WriteFile(tempSoureFileName, info, 0666)
	// if err != nil {
	// 	utils.WLog("Error: could not create json file", ClientID)
	// 	removeVideo(path, filename, ClientID)
	// 	return vi, err
	// }

	// // Run python script to get nesessary data from json file
	// gpath, err := filepath.Abs(utils.Conf.DataGen)
	// wg.Add(1)
	// err = generateDataFile(&wg, gpath, prefix)
	// wg.Wait()
	// if err != nil {
	// 	utils.WLog("Error: failed to generate video data", ClientID)
	// 	removeVideo(path, filename, ClientID)
	// 	return vi, err
	// }

	// // Write data to Vidinfo struct
	// tempDataFileName := fmt.Sprintf(utils.Conf.DataJson, prefix)
	// vi, err = parseFile(tempDataFileName)
	// if err != nil || vi.IsEmpty() {
	// 	utils.WLog("Error: failed parsing data file", ClientID)
	// 	removeVideo(path, filename, ClientID)
	// 	return vi, err
	// }

	return vi, nil
}

// func parseFile(f string) (models.Vidinfo, error) {
// 	var (
// 		vi models.Vidinfo
// 	)
// 	file, err := os.Open(f)
// 	if err != nil {
// 		return vi, err
// 	}
// 	defer file.Close()
// 	defer os.Remove(f)

// 	byteValue, err := ioutil.ReadAll(file)
// 	if err != nil {
// 		return vi, err
// 	}

// 	err = json.Unmarshal(byteValue, &vi)
// 	if err != nil {
// 		return vi, err
// 	}

// 	return vi, nil
// }
