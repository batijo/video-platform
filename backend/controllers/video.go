package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"

	"github.com/gorilla/mux"
)

// FetchVideos return all
func FetchVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	res := utils.DB.Preload("AudioT").Preload("SubtitleT").Find(&videos)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not fetch videos", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: videos}
	json.NewEncoder(w).Encode(resp)
}

// UpdateVideo ...
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	video := &models.Video{}
	var id = mux.Vars(r)["id"]
	res := utils.DB.First(&video, id)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Video not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	json.NewDecoder(r.Body).Decode(video)
	res = utils.DB.Save(&video)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not save video", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Video Updated", Data: video}
	json.NewEncoder(w).Encode(resp)
}

// DeleteVideo ...
func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	var video models.Video
	res := utils.DB.First(&video, id)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Video not found", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
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
		resp := models.Response{Message: fmt.Sprintf("Unkown video file state: %s", video.State)}
		json.NewEncoder(w).Encode(resp)
		filePath = utils.Conf.SD
	}

	if err := os.Remove(filePath + video.FileName); err != nil {
		log.Println(err)
		resp := models.Response{Status: false, Message: "Could not delete video", Error: err.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}
	utils.DB.Delete(&video)

	resp := models.Response{Status: true, Message: "Video Deleted"}
	json.NewEncoder(w).Encode(resp)
}

// GetVideo ...
func GetVideo(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	var video models.Video
	res := utils.DB.Preload("AudioT").Preload("SubtitleT").First(&video, id)

	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not get video", Error: res.Error.Error()}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: video}
	json.NewEncoder(w).Encode(resp)
}
