package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/transcode"
	"github.com/batijo/video-platform/backend/utils"
)

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
	var (
		encData    models.Encode
		presetData models.Pdata
		err        error
	)

	// If you didn't add comments at first, later you don't know what this shit does
	err = r.ParseForm()
	if err != nil {
		resp := models.Response{Status: false, Message: "Some ParseForm do not like you", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return
	}

	// Decode json file
	if utils.Conf.Presets {
		err = json.NewDecoder(r.Body).Decode(&presetData)
		if err != nil {
			resp := models.Response{Status: false, Message: "Cannot decode json", Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			log.Println(err)
			return
		}
	} else {
		err = json.NewDecoder(r.Body).Decode(&encData)
		if err != nil {
			resp := models.Response{Status: false, Message: "Cannot decode json", Error: err.Error()}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(resp)
			log.Println(err)
			return
		}
	}
	// data := models.VfNPrd {
	// 	Pdata: prd,
	// 	Enc:   vf,
	// 	Err:   err,
	// }
	// vfnprd <- data

	err = transcode.AddToQueue(encData)
	if err != nil {
		resp := models.Response{Status: false, Message: "Error adding to queue", Error: err.Error()}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return
	}

	resp := models.Response{Status: true, Message: "Adding video to transcoder queue"}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
}
