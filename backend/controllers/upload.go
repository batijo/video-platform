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
	"github.com/batijo/video-platform/backend/utils/auth"
)

// VideoUpload upload handler which only allows to upload video
func VideoUpload(w http.ResponseWriter, r *http.Request) {
	//Starts reading file by chuncking <- NOPE
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		resp := models.Response{Status: false, Message: "Failed to upload file", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		utils.WLog("Error: failed to upload file", r.RemoteAddr)
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
		utils.WLog("Error: this file format is not allowed: "+filepath.Ext(handler.Filename), r.RemoteAddr)
		return
	}

	// Checks if uploaded file with the same name already exists
	fileName := utils.ReturnDifNameIfDublicate(handler.Filename, utils.Conf.SD)
	if fileName != handler.Filename {
		utils.WLog("File with the same name already exist so it has been changed", r.RemoteAddr)
	}

	//Create empty file in /videos folder
	utils.WLog("Creating file", r.RemoteAddr)
	dst, err := os.OpenFile(utils.Conf.SD+fileName, os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not create file", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		utils.WLog("Error: could not create file", r.RemoteAddr)
		return
	}

	//Copies a temporary file to empty file in /videos folder
	utils.WLog("Writing to file", r.RemoteAddr)
	if _, err := io.Copy(dst, file); err != nil {
		resp := models.Response{Status: false, Message: "Failed to write file", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, -1, r.RemoteAddr)
		utils.WLog("Error: failed to write file", r.RemoteAddr)
		return
	}

	utils.WLog("Upload successful", r.RemoteAddr)

	data, err := tc.GetVidInfo(utils.Conf.SD, fileName, r.RemoteAddr, -1)
	if err != nil {
		resp := models.Response{Status: false, Message: "Error getting video info", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, -1, r.RemoteAddr)
		utils.WLog("Error: could not get video info", r.RemoteAddr)
		return
	}

	userID, _, err := auth.GetUserID(r)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not verify user", Error: err.Error()}
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, -1, r.RemoteAddr)
		return
	}

	vidId, err := utils.InsertVideo(data, "not_transcoded", userID, -1)
	if err != nil {
		resp := models.Response{Status: false, Message: "Sql error", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD+fileName, int(vidId), r.RemoteAddr)
		utils.WLog("Error: failed to insert video data in database", r.RemoteAddr)
		return
	}

	utils.UpdateMessage(fileName)
	if utils.Conf.Presets {
		dataWP := utils.AddPresetsToJSON(data)
		resp := models.Response{Status: true, Message: fileName, Data: dataWP}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	} else {
		var videoData models.Video
		videoData.ParseWithVidinfo(data)
		videoData.ID = vidId
		resp := models.Response{Status: true, Message: fileName, Data: videoData}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func removeVideo(path string, vidId int, ClientID string) error {
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
