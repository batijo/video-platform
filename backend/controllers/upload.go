package controllers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/batijo/video-platform/backend/models"
	tc "github.com/batijo/video-platform/backend/transcode"
	"github.com/batijo/video-platform/backend/utils"
)

// VideoUpload upload handler which only allows to upload video
func VideoUpload(w http.ResponseWriter, r *http.Request) {
	// Gets user ID
	userID, _, err := utils.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not verify user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return
	}

	//Starts reading file by chuncking <- NOPE
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		resp := models.Response{Status: false, Message: "Failed to upload file", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		utils.WLog("Error: failed to upload file", userID)
		return
	}
	defer file.Close()

	// Check if video file format is allowed
	allowed := false
	for _, ave := range utils.Conf.FileTypes {
		if filepath.Ext(handler.Filename) == ave {
			allowed = true
		}
	}
	if !allowed {
		resp := models.Response{Status: false, Message: "This file format is not allowed " + filepath.Ext(handler.Filename)}
		w.WriteHeader(http.StatusUnsupportedMediaType)
		json.NewEncoder(w).Encode(resp)
		utils.WLog("Error: this file format is not allowed: "+filepath.Ext(handler.Filename), userID)
		return
	}

	// Checks if uploaded file with the same name already exists
	fileName := utils.ReturnDifNameIfDublicate(handler.Filename, utils.Conf.SD)
	if fileName != handler.Filename {
		utils.WLog("File with the same name already exist so it has been changed", userID)
	}

	//Create empty file in /videos folder
	utils.WLog("Creating file", userID)
	dst, err := os.OpenFile(utils.Conf.SD+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not create file", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		utils.WLog("Error: could not create file", userID)
		return
	}

	//Copies a temporary file to empty file in /videos folder
	utils.WLog("Writing to file", userID)
	if _, err := io.Copy(dst, file); err != nil {
		resp := models.Response{Status: false, Message: "Failed to write file", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, -1, userID)
		utils.WLog("Error: failed to write file", userID)
		return
	}

	utils.WLog("Upload successful", userID)

	data, err := tc.GetVidInfo(utils.Conf.SD, fileName, userID, -1)
	if err != nil {
		resp := models.Response{Status: false, Message: "Error getting video info", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, -1, userID)
		utils.WLog("Error: could not get video info", userID)
		return
	}

	vidId, err := utils.InsertVideo(data, "not_transcoded", userID, -1)
	if err != nil {
		resp := models.Response{Status: false, Message: "Database error", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, int(vidId), userID)
		utils.WLog("Error: failed to insert video data in database", userID)
		return
	}

	var videoData models.Video
	if res := utils.DB.Preload("AudioT").Preload("SubtitleT").Where("id = ?", vidId).First(&videoData); res.Error != nil {
		resp := models.Response{Status: false, Message: "Database error", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, int(vidId), userID)
		utils.WLog("Error: failed to retrieve video data from database", userID)
		return
	}

	utils.UpdateUserMessage(fileName, userID)
	resp := models.Response{Status: true, Message: fileName, Data: videoData}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}

func removeVideo(path string, vidId int, ClientID uint) error {
	var err error

	if err = os.Remove(path); err != nil {
		log.Println(err)
		utils.WLog("Error: failed removing source file", ClientID)
	}
	if vidId >= 0 {
		err = utils.DeleteVideo(uint(vidId))
		if err != nil {
			return err
		}
	}

	return err
}
