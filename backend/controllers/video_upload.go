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

	"github.com/gorilla/mux"
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
		return

		// w.WriteHeader(500)
		// utils.WLog("Error: failed to upload file", r.RemoteAddr)
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
		return

		// utils.WLog("Error: this file format is not allowed "+filepath.Ext(handler.Filename), r.RemoteAddr)
		// w.WriteHeader(403)
	}

	// Checks if uploaded file with the same name already exists
	if _, err := os.Stat(utils.Conf.SD + handler.Filename); err == nil {
		var resp = map[string]interface{}{"status": false, "message": fmt.Sprintf("File \"%v\" already exists", handler.Filename), "error": nil}
		json.NewEncoder(w).Encode(resp)
		return

		// utils.WLog(fmt.Sprintf("Error: file \"%v\" already exists", handler.Filename), r.RemoteAddr)
		// w.WriteHeader(403)
	}

	//Create empty file in /videos folder
	//utils.WLog("Creating file", r.RemoteAddr)
	dst, err := os.OpenFile(utils.Conf.SD+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	defer dst.Close()
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Could not create file", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return

		// w.WriteHeader(500)
		// utils.WLog("Error: could not create file", r.RemoteAddr)
	}

	//Copies a temporary file to empty file in /videos folder
	utils.WLog("Writing to file", r.RemoteAddr)
	if _, err := io.Copy(dst, file); err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Failed to write file", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		return

		// w.WriteHeader(500)
		// utils.WLog("Error: failed to write file", r.RemoteAddr)
	}

	var resp = map[string]interface{}{"status": true, "message": "Upload successful", "error": nil}
	json.NewEncoder(w).Encode(resp)

	//utils.WLog("Upload successful", r.RemoteAddr)

	data, err := writeJSONResponse(w, handler.Filename, r.RemoteAddr)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Error getting video info", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		return

		// w.WriteHeader(500)
		// utils.WLog("Error: failed send video data to client", r.RemoteAddr)
	}

	userID, err := auth.GetUserID(r)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Could not verify user", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		return
	}

	err = utils.InsertVideo(data, handler.Filename, "Not transcoded", userID, -1)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Sql error", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile(utils.Conf.SD, handler.Filename, r.RemoteAddr)
		return

		// utils.WLog("Error: failed to insert video data in database", r.RemoteAddr)
		// w.WriteHeader(500)
	}

	// utils.UpdateMessage(handler.Filename)

	go func() {
		dat := <-vfnprd
		if dat.Err == nil {
			vf := dat.Video
			prd := dat.Pdata
			go tc.ProcessVodFile(handler.Filename, data, vf, prd, utils.Conf, r.RemoteAddr, userID)
		}
	}()
}

// TcTypeHandler ...
func TcTypeHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		w.WriteHeader(422)
		return
	}

	type response struct {
		Typechange string `json:"Tc"`
	}
	var rsp response
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&rsp)
	if err != nil {
		log.Println(err)
		w.WriteHeader(415)
		return
	}

	if rsp.Typechange == "true" {
		utils.Conf.Presets = false
	} else if rsp.Typechange == "false" {
		utils.Conf.Presets = true
	} else {
		log.Println(fmt.Errorf("uknown change type: '%v', expected 'true' or 'false'", rsp.Typechange))
		w.WriteHeader(415)

	}

	w.WriteHeader(200)
}

// TranscodeHandler ...
func TranscodeHandler(w http.ResponseWriter, r *http.Request) {

	var err error

	err = r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	var (
		vf  models.Video
		prd models.Pdata
	)

	// Decode json file
	if utils.Conf.Presets {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&prd)
		if err != nil {
			var resp = map[string]interface{}{"status": false, "message": "cannot decode json", "error": err}
			json.NewEncoder(w).Encode(resp)
			log.Println(err)
			return
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&vf)
		if err != nil {
			var resp = map[string]interface{}{"status": false, "message": "cannot decode json", "error": err}
			json.NewEncoder(w).Encode(resp)
			log.Println(err)
			return
		}
	}

	data := models.VfNPrd{
		prd,
		vf,
		err,
	}

	vfnprd <- data

	var resp = map[string]interface{}{"status": true, "message": "transcode starting", "error": nil}
	json.NewEncoder(w).Encode(resp)
}

// Send json response after file upload
func writeJSONResponse(w http.ResponseWriter, filename string, clid string) (models.Vidinfo, error) {
	var (
		data    models.Data
		vidinfo models.Vidinfo
		err     error
		//info    []byte
	)

	if err != nil {
		return vidinfo, err
	}

	vidinfo, err = tc.GetVidInfo(utils.Conf.SD, filename, utils.Conf.TempJson, utils.Conf.DataGen, utils.Conf.TempTxt, clid)
	if err != nil {
		return vidinfo, err
	}

	if utils.Conf.Presets {
		data = utils.AddPresetsToJSON(vidinfo)
		json.NewEncoder(w).Encode(&data)

		// info, err = json.Marshal(data)
		// if err != nil {
		// 	return vidinfo, err
		// }
	} else {
		json.NewEncoder(w).Encode(&vidinfo)
		// info, err = json.Marshal(vidinfo)
		// if err != nil {
		// 	return vidinfo, err
		// }
	}

	//w.WriteHeader(200)
	//w.Write(info)

	//json.NewEncoder(w).Encode(&vidinfo)

	return vidinfo, nil
}

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

func removeFile(path string, filename string, clid string) error {

	var err error
	if err = os.Remove(path + filename); err != nil {
		utils.WLog("Error: failed removing source file", clid)
	}
	//db.RemoveRowByName(filename, "Video")
	return err
}
