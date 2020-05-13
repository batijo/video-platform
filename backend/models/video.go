package models

import "github.com/jinzhu/gorm"

type Video struct {
	gorm.Model

	//ID         int `json:"id,primary_key"`
	UserID     uint `gorm:"TYPE:integer REFERENCES User"`
	StreamID   int
	FileName   string
	State      string
	VideoCodec string
	Width      int
	Height     int
	FrameRate  float64
	//VtRes     string
	//Save      bool
	AudioT    []Audio `gorm:"ForeignKey:VideoID"`
	SubtitleT []Sub   `gorm:"ForeignKey:VideoID"`
}

type Audio struct {
	ID       int  `json:"id,primary_key"`
	VideoID  uint `gorm:"TYPE:integer REFERENCES Videos"`
	StreamID int
	AtCodec  string
	Language string
	Channels int
}

type Sub struct {
	ID       int  `json:"id,primary_key"`
	VideoID  uint `gorm:"TYPE:integer REFERENCES Videos"`
	StreamID int
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
