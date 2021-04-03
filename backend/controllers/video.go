package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"

	"github.com/gorilla/mux"
)

// FetchVideos return all
func FetchVideos(w http.ResponseWriter, r *http.Request) {
	var videos []models.Video
	utils.DB.Preload("AudioT").Preload("SubtitleT").Find(&videos)

	json.NewEncoder(w).Encode(videos)
}

// UpdateVideo ...
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	video := &models.Video{}
	var id = mux.Vars(r)["id"]
	utils.DB.First(&video, id)
	json.NewDecoder(r.Body).Decode(video)
	utils.DB.Save(&video)
	json.NewEncoder(w).Encode(video)
}

// DeleteVideo ...
func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	var video models.Video
	utils.DB.First(&video, id)
	utils.DB.Delete(&video)
	json.NewEncoder(w).Encode("Video Deleted")
}

// GetVideo ...
func GetVideo(w http.ResponseWriter, r *http.Request) {
	var id = mux.Vars(r)["id"]
	var video models.Video
	utils.DB.Preload("AudioT").Preload("SubtitleT").First(&video, id)
	json.NewEncoder(w).Encode(&video)
}
