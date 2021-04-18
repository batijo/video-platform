package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/transcode"
	"github.com/batijo/video-platform/backend/utils"
)

func ReturnQueue(w http.ResponseWriter, r *http.Request) {
	var (
		videos []models.Video
		queue  models.Queue
	)

	if transcode.Active() {
		resp := models.Response{Status: true, Message: "No videos are transcoding"}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
		return
	}

	userId, _, err := utils.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not authorise user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err := utils.DB.Find(&videos).Error; err != nil {
		resp := models.Response{Status: false, Message: "Sql error", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	queue.Put(videos, userId)

	resp := models.Response{Status: true, Message: "Success", Data: queue}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
