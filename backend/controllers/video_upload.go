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
	config, err := utils.GetConf()
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Failed to open conf.toml", "error": nil}
		json.NewEncoder(w).Encode(resp)
		return
	}
	allowed := false
	for _, ave := range config.FileTypes {
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
	if _, err := os.Stat("videos/" + handler.Filename); err == nil {
		var resp = map[string]interface{}{"status": false, "message": fmt.Sprintf("File \"%v\" already exists", handler.Filename), "error": nil}
		json.NewEncoder(w).Encode(resp)
		return

		// utils.WLog(fmt.Sprintf("Error: file \"%v\" already exists", handler.Filename), r.RemoteAddr)
		// w.WriteHeader(403)
	}

	//Create empty file in /videos folder
	//utils.WLog("Creating file", r.RemoteAddr)
	dst, err := os.OpenFile("videos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
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
		removeFile("videos/", handler.Filename, r.RemoteAddr)
		return

		// w.WriteHeader(500)
		// utils.WLog("Error: failed to write file", r.RemoteAddr)
	}

	var resp = map[string]interface{}{"status": true, "message": "Upload successful", "error": nil}
	json.NewEncoder(w).Encode(resp)

	//utils.WLog("Upload successful", r.RemoteAddr)

	_, err = writeJSONResponse(w, handler.Filename, r.RemoteAddr)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Error getting video info", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile("videos/", handler.Filename, r.RemoteAddr)

		// w.WriteHeader(500)
		// utils.WLog("Error: failed send video data to client", r.RemoteAddr)
	}

	// err = db.InsertVideo(data, handler.Filename, "Not transcoded", -1)
	// if err != nil {
	// 	utils.WLog("Error: failed to insert video data in database", r.RemoteAddr)
	// 	log.Println(err)
	// 	w.WriteHeader(500)
	// 	removeFile("./videos/", handler.Filename, r.RemoteAddr)
	// 	return
	// }

	// utils.UpdateMessage(handler.Filename)

	// if CONF.Advanced {
	// 	go func() {
	// 		dat := <-vfnprd
	// 		if dat.Err == nil {
	// 			vf := dat.Video
	// 			prd := dat.PData
	// 			go tc.ProcessVodFile(handler.Filename, data, vf, prd, CONF, r.RemoteAddr)
	// 		}
	// 	}()
	// } else {
	// 	//go tc.ProcessVodFile(handler.Filename, data, vf, prd, CONF, r.RemoteAddr)
	// }
}

// Send json response after file upload
func writeJSONResponse(w http.ResponseWriter, filename string, clid string) (models.Vidinfo, error) {
	var (
		//data    models.Data
		vidinfo models.Vidinfo
		err     error
		//info    []byte
	)

	config, err := utils.GetConf()

	vidinfo, err = tc.GetVidInfo("videos/", filename, config.TempJson, config.DataGen, config.TempTxt, clid)
	if err != nil {
		log.Println(err)
		return vidinfo, err
	}

	// if config.Presets {
	// 	data, err = db.AddPresetsToJson(vidinfo)
	// 	if err != nil {
	// 		return vidinfo, err
	// 	}

	// 	info, err = json.Marshal(data)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return vidinfo, err
	// 	}
	// } else {
	// 	info, err = json.Marshal(vidinfo)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return vidinfo, err
	// 	}
	// }

	// info, err = json.Marshal(vidinfo)
	// if err != nil {
	// 	log.Println(err)
	// 	return vidinfo, err
	// }

	//w.WriteHeader(200)
	//w.Write(info)

	json.NewEncoder(w).Encode(&vidinfo)

	return vidinfo, nil
}

func removeFile(path string, filename string, clid string) error {

	var err error
	if err = os.Remove(path + filename); err != nil {
		utils.WLog("Error: failed removing source file", clid)
	}
	//db.RemoveRowByName(filename, "Video")
	return err
}
