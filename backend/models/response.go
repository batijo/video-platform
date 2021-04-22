package models

import (
	"errors"
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

func (vr *VideoResponse) SerializeWithVideo(v Video) {
	vr.Video = v
}

func (vr *VideoResponse) SerializeWithStream(s Vstream) error {
	if len(s.Video) < 1 {
		return errors.New("stream must have at least one video element")
	}
	vr.Video = s.Video[0]
	vr.Video.Title = s.Name
	for i := 0; i < len(s.Video); i++ {
		vr.Resolution = append(vr.Resolution, fmt.Sprintf("%vp", s.Video[i].Width))
	}

	return nil
}
