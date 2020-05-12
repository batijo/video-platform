package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"

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

	data, err := writeJSONResponse(w, handler.Filename, r.RemoteAddr)
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Error getting video info", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile("videos/", handler.Filename, r.RemoteAddr)

		// w.WriteHeader(500)
		// utils.WLog("Error: failed send video data to client", r.RemoteAddr)
	}

	err = InsertVideo(data, handler.Filename, "Not transcoded")
	if err != nil {
		var resp = map[string]interface{}{"status": false, "message": "Sql error", "error": err}
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		removeFile("videos/", handler.Filename, r.RemoteAddr)
		return

		// utils.WLog("Error: failed to insert video data in database", r.RemoteAddr)
		// w.WriteHeader(500)
	}

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

// InsertVideo adds video to database
func InsertVideo(vidinfo models.Vidinfo, name string, state string) error {

	var audio []models.Audio
	var subtitle []models.Sub

	for _, a := range vidinfo.Audiotrack {
		at := models.Audio{
			StreamID: a.Index,
			AtCodec:  a.CodecName,
			Language: a.Language,
			Channels: a.Channels,
		}
		audio = append(audio, at)
	}

	for _, s := range vidinfo.Subtitle {
		st := models.Sub{
			StreamID: s.Index,
			Language: s.Language,
		}
		subtitle = append(subtitle, st)
	}

	video := models.Video{
		StreamID:   vidinfo.Videotrack[0].Index,
		FileName:   name,
		State:      state,
		VideoCodec: vidinfo.Videotrack[0].CodecName,
		Width:      vidinfo.Videotrack[0].Width,
		Height:     vidinfo.Videotrack[0].Height,
		FrameRate:  vidinfo.Videotrack[0].FrameRate,
		AudioT:     audio,
		SubtitleT:  subtitle,
	}

	createdVideo := utils.DB.Create(&video)

	if createdVideo.Error != nil {
		return createdVideo.Error
	}

	return nil
}

func removeFile(path string, filename string, clid string) error {

	var err error
	if err = os.Remove(path + filename); err != nil {
		utils.WLog("Error: failed removing source file", clid)
	}
	//db.RemoveRowByName(filename, "Video")
	return err
}
