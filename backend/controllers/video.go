package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
	"github.com/batijo/video-platform/backend/utils/auth"
	"github.com/jinzhu/gorm"

	"github.com/gorilla/mux"
)

// FetchVideos return all
func FetchVideos(w http.ResponseWriter, r *http.Request) {
	var (
		videos []models.Video
		res    *gorm.DB
	)
	userId, admin, err := auth.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if admin {
		res = utils.DB.Preload("AudioT").Preload("SubtitleT").Find(&videos)
	} else {
		res = utils.DB.Preload("AudioT").Preload("SubtitleT").Where("user_id = ? OR public = ?", userId, true).Find(&videos)
	}
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not fetch videos", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: videos}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// UpdateVideo ...
func UpdateVideo(w http.ResponseWriter, r *http.Request) {
	video := &models.Video{}
	var id = mux.Vars(r)["id"]

	userId, admin, err := auth.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res := utils.DB.First(&video, id)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Video not found", Error: res.Error.Error()}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch true {
	case userId == video.UserID:
		break
	case admin:
		break
	default:
		resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&video)
	if err != nil {
		resp := models.Response{Status: false, Message: "Invalid request", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res = utils.DB.Save(&video)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not save video", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Video Updated", Data: video}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// DeleteVideo ...
func DeleteVideo(w http.ResponseWriter, r *http.Request) {
	var (
		id    = mux.Vars(r)["id"]
		video models.Video
	)

	userId, admin, err := auth.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res := utils.DB.First(&video, id)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Video not found", Error: res.Error.Error()}
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch true {
	case userId == video.UserID:
		break
	case admin:
		break
	default:
		resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
		w.WriteHeader(http.StatusForbidden)
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
		filePath = utils.Conf.SD
	}

	if err := os.Remove(filePath + video.FileName); err != nil {
		log.Println(err)
		resp := models.Response{Status: false, Message: "Could not delete video", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	res = utils.DB.Delete(&video)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not delete video", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Video Deleted"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

// GetVideo ...
func GetVideo(w http.ResponseWriter, r *http.Request) {
	var (
		id    = mux.Vars(r)["id"]
		video models.Video
	)

	userId, admin, err := auth.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	res := utils.DB.Preload("AudioT").Preload("SubtitleT").First(&video, id)
	if res.Error != nil {
		resp := models.Response{Status: false, Message: "Could not get video", Error: res.Error.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	switch true {
	case video.Public:
		break
	case userId == video.UserID:
		break
	case admin:
		break
	default:
		resp := models.Response{Status: false, Message: "You have no privilage to perform this action"}
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := models.Response{Status: true, Message: "Success", Data: video}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
