package controllers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/transcode"
	"github.com/gorilla/mux"
)

// TranscodeHandler ...
func TranscodeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		encData models.Encode
		err     error
		vidID   uint
	)

	vidID, err = getId(w, r)
	if err != nil {
		return
	}
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
	err = json.NewDecoder(r.Body).Decode(&encData)
	if err != nil {
		resp := models.Response{Status: false, Message: "Cannot decode json", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return
	}

	err = transcode.AddToQueue(encData, []models.Stream{}, vidID)
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

func PresetTranscodeHandler(w http.ResponseWriter, r *http.Request) {
	var (
		presetData []models.Stream
		err        error
		vidID      uint
	)

	vidID, err = getId(w, r)
	if err != nil {
		return
	}
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
	err = json.NewDecoder(r.Body).Decode(&presetData)
	if err != nil {
		resp := models.Response{Status: false, Message: "Cannot decode json", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return
	}

	err = transcode.AddToQueue(models.Encode{}, presetData, vidID)
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

func getId(w http.ResponseWriter, r *http.Request) (uint, error) {
	id := mux.Vars(r)["id"]
	i, err := strconv.Atoi(id)
	if err != nil {
		resp := models.Response{Status: false, Message: "Could not get id from URL", Error: err.Error()}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		log.Println(err)
		return uint(i), err
	}
	return uint(i), nil
}
