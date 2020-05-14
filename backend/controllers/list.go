package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Dzionys/video-platform/backend/models"
	"github.com/Dzionys/video-platform/backend/utils"
)

// ListHandler ...
func ListHandler(w http.ResponseWriter, r *http.Request) {

	data := putVideosToJSON()

	dt, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}

	w.Write(dt)
}

func putVideosToJSON() models.VideoData {
	var streams []models.Vstream
	var videos []models.Video
	var videodata models.VideoData

	utils.DB.Find(&streams)
	utils.DB.Where("vstream_id IS NULL").Find(&videos)

	for _, s := range streams {
		var tempvideo models.VideoStream

		for _, v := range s.Video {
			tempvideo.Video = append(tempvideo.Video, v)
		}

		tempvideo.StreamName = s.Name
		tempvideo.State = s.Video[0].State
		tempvideo.Stream = true

		videodata.VideoStream = append(videodata.VideoStream, tempvideo)
	}

	for _, v := range videos {
		var tempvideo models.VideoStream

		tempvideo.Video = append(tempvideo.Video, v)
		tempvideo.StreamName = v.FileName
		tempvideo.State = v.State
		tempvideo.Stream = false

		videodata.VideoStream = append(videodata.VideoStream, tempvideo)
	}

	return videodata
}