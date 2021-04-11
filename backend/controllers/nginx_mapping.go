package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
	"github.com/gorilla/mux"
)

// NginxMappingHandler ...
func NginxMappingHandler(w http.ResponseWriter, r *http.Request) {

	tcvidpath := utils.Conf.APTGD + "/go/src/github.com/batijo/video-platform/videos/transcoded/%v"
	var sqncs models.Sequences

	vars := mux.Vars(r)

	if filepath.Ext(vars["name"]) == ".mp4" {
		//err := db.IsExist("Video", vars["name"])
		if err := utils.DB.Where("file_name = ?", vars["name"]).First(&models.Video{}).Error; err != nil {
			log.Println(err)
			w.WriteHeader(404)
			return
		}

		temp := models.Clip{
			Type: "source",
			Path: fmt.Sprintf(tcvidpath, vars["name"]),
		}
		var tempclip models.Clips
		tempclip.Clips = append(tempclip.Clips, temp)
		sqncs.Sequences = append(sqncs.Sequences, tempclip)

	} else {
		names, err := getAllStreamNames(vars["name"])
		if err != nil {
			log.Println(err)
			w.WriteHeader(404)
			return
		}

		for _, n := range names {
			temp := models.Clip{
				Type: "source",
				Path: fmt.Sprintf(tcvidpath, n),
			}
			var tempclip models.Clips
			tempclip.Clips = append(tempclip.Clips, temp)
			sqncs.Sequences = append(sqncs.Sequences, tempclip)
		}
	}

	j, err := json.Marshal(sqncs)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(j)
}

func getAllStreamNames(sname string) ([]string, error) {
	var stream models.Vstream
	var names []string

	if err := utils.DB.Where("Name = ?", sname).First(&stream).Error; err != nil {
		return names, err
	}

	for _, vid := range stream.Video {
		names = append(names, vid.FileName)
	}

	return names, nil
}
