package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Dzionys/video-platform/backend/models"
	tc "github.com/Dzionys/video-platform/backend/transcode"
	"github.com/Dzionys/video-platform/backend/utils"
	"github.com/Dzionys/video-platform/backend/utils/auth"
)

var (
	vfnprd = make(chan models.VfNPrd)
)

// VideoUpload upload handler which only allows to upload video
func VideoUpload(w http.ResponseWriter, r *http.Request) {

	//Starts readig file by chuncking
	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("file")
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Failed to upload file", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		utils.WLog("Error: failed to upload file", r.RemoteAddr)
		return
	}
	defer file.Close()

	// Check if video file format is allowed
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Failed to open conf.toml", "error": nil}
		json.NewEncoder(w).Encode(resp)
		return
	}
	allowed := false
	for _, ave := range utils.Conf.FileTypes {
		if filepath.Ext(handler.Filename) == ave {
			allowed = true
		}
	}
	if !allowed {
		var resp = map[string]interface{}{"status": false, "message": "This file format is not allowed " + filepath.Ext(handler.Filename), "error": nil}
		json.NewEncoder(w).Encode(resp)
		utils.WLog("Error: this file format is not allowed: "+filepath.Ext(handler.Filename), r.RemoteAddr)
		return
	}

	// Checks if uploaded file with the same name already exists
	if _, err := os.Stat(utils.Conf.SD + handler.Filename); err == nil {
		var resp = map[string]interface{}{"status": false, "message": fmt.Sprintf("File \"%v\" already exists", handler.Filename), "error": nil}
		json.NewEncoder(w).Encode(resp)
		utils.WLog(fmt.Sprintf("Error: file \"%v\" already exists", handler.Filename), r.RemoteAddr)
		return
	}

	//Create empty file in /videos folder
	utils.WLog("Creating file", r.RemoteAddr)
	dst, err := os.OpenFile(utils.Conf.SD+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Could not create file", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		utils.WLog("Error: could not create file", r.RemoteAddr)
		return
	}

	//Copies a temporary file to empty file in /videos folder
	utils.WLog("Writing to file", r.RemoteAddr)
	if _, err := io.Copy(dst, file); err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Failed to write file", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		utils.WLog("Error: failed to write file", r.RemoteAddr)
		return
	}

	var resp = map[string]interface{}{"status": true, "message": "Upload successful", "error": nil}
	json.NewEncoder(w).Encode(resp)

	utils.WLog("Upload successful", r.RemoteAddr)

	data, err := writeJSONResponse(w, handler.Filename, r.RemoteAddr)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Error getting video info", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		utils.WLog("Error: failed send video data to client", r.RemoteAddr)
		return
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Could not verify user", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		return
	}

	err = utils.InsertVideo(data, handler.Filename, "not_transcoded", userID, -1)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Sql error", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeVideo(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		utils.WLog("Error: failed to insert video data in database", r.RemoteAddr)
		return
	}

	utils.UpdateMessage(handler.Filename)

	go func() {
		dat := <-vfnprd
		if dat.Err == nil {
			vf := dat.Video
			prd := dat.Pdata
			go tc.ProcessVodFile(handler.Filename, data, vf, prd, r.RemoteAddr, userID)
		}
	}()
}

// Send json response after file upload
func writeJSONResponse(w http.ResponseWriter, filename string, ClientID string) (models.Vidinfo, error) {
	var (
		data    models.Data
		vidinfo models.Vidinfo
		err     error
	)

	if err != nil {
		return vidinfo, err
	}

	vidinfo, err = tc.GetVidInfo(utils.Conf.SD, filename, utils.Conf.TempJson, utils.Conf.DataGen, utils.Conf.TempTxt, ClientID)
	if err != nil {
		return vidinfo, err
	}

	if utils.Conf.Presets {
		data = utils.AddPresetsToJSON(vidinfo)
		json.NewEncoder(w).Encode(&data)
	} else {
		json.NewEncoder(w).Encode(&vidinfo)
	}

	return vidinfo, nil
}

func removeVideo(path string, filename string, ClientID string) error {

	var err error
	if err = os.Remove(path + filename); err != nil {
		log.Println(err)
		utils.WLog("Error: failed removing source file", ClientID)
	}
	err = utils.DeleteVideo(filename)
	if err != nil {
		return err
	}
	return err
}
