package models

import (
	"fmt"
)

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Error   string      `json:"error"`
	Data    interface{} `json:"data"`
}

type VideoResponse struct {
	Video
	Resolution []string `json:"resolution"`
}

func SerializeWithVideo(v Video) VideoResponse {
	var vr VideoResponse
	vr.Video = v
	return vr
}

func SerializeWithStream(s Vstream) VideoResponse {
	var vr VideoResponse
	if len(s.Video) < 1 {
		return VideoResponse{}
	}
	vr.Video = s.Video[0]
	vr.Video.FileName = s.Name
	for i := 0; i < len(s.Video); i++ {
		vr.Resolution = append(vr.Resolution, fmt.Sprintf("%vp", s.Video[i].Height))
	}

	return vr
}
