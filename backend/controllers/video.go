package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"

	"github.com/gorilla/mux"
)

// FetchVideos return all
func FetchVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	result := utils.DB.Preload("AudioT").Preload("SubtitleT").Find(&videos)

	if result.Error != nil {
		json.NewEncoder(w).Encode(result.Error.Error())
		return
	}

	json.NewEncoder(w).Encode(videos)
}

// UpdateVideo ...
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	video := &models.Video{}
	var id = mux.Vars(r)["id"]
	result := utils.DB.First(&video, id)

	if result.Error != nil {
		json.NewEncoder(w).Encode(result.Error.Error())
		return
	}

	json.NewDecoder(r.Body).Decode(video)
	utils.DB.Save(&video)
	json.NewEncoder(w).Encode(video)
}

// DeleteVideo ...
func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	var video models.Video
	result := utils.DB.First(&video, id)

	if result.Error != nil {
		json.NewEncoder(w).Encode(result.Error.Error())
		return
	}

	var filePath string
	switch video.State {
	case "not_transcoded":
		filePath = utils.Conf.SD
		break
	case "transcoding":
		filePath = utils.Conf.TD
		break
	case "transcoded":
		filePath = utils.Conf.DD
		break
	default:
		json.NewEncoder(w).Encode(fmt.Sprintf("Unkown video file state: %s", video.State))
		filePath = utils.Conf.SD
	}

	if err := os.Remove(filePath + video.FileName); err != nil {
		log.Println(err)
		return
	}
	utils.DB.Delete(&video)

	json.NewEncoder(w).Encode("Video Deleted")
}

// GetVideo ...
func GetVideo(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	var video models.Video
	result := utils.DB.Preload("AudioT").Preload("SubtitleT").First(&video, id)

	if result.Error != nil {
		json.NewEncoder(w).Encode(result.Error.Error())
		return
	}

	json.NewEncoder(w).Encode(&video)
}
