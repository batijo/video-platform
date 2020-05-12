package models

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model

	FileName  string
	VtId      int
	VtCodec   string
	FrameRate float64
	VtRes     string
	Save      bool
	AudioT    []Audio
	SubtitleT []Sub
}

type Audio struct {
	gorm.Model

	VideoID  int
	AtCodec  string
	Language string
	Channels int
}

type Sub struct {
	gorm.Model

	VideoID  int
	Language string
}

type VideoData struct {
	VideoStream []VideoStream
}

type VideoStream struct {
	Stream     bool
	StreamName string
	State      string
	Thumbnail  string
	Video      []Video
}
