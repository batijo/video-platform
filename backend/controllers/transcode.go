package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"
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
			resp := models.Response{Status: false, Message: "cannot decode json", Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			log.Println(err)
			return
		}
	} else {
		decoder := json.NewDecoder(r.Body)
		err = decoder.Decode(&vf)
		if err != nil {
			resp := models.Response{Status: false, Message: "cannot decode json", Error: err.Error()}
			json.NewEncoder(w).Encode(resp)
			log.Println(err)
			return
		}
	}

	data := models.VfNPrd{
		Pdata: prd,
		Video: vf,
		Err:   err,
	}

	vfnprd <- data

	resp := models.Response{Status: true, Message: "transcode starting", Error: err.Error()}
	json.NewEncoder(w).Encode(resp)
}
