package models

import "github.com/jinzhu/gorm"

type Vstream struct {
	gorm.Model

	UserID uint
	Name   string
	Video  []Video `gorm:"ForeignKey:VstreamID"`
}

type Video struct {
	gorm.Model

	UserID     uint
	Public     bool
	VstreamID  uint `gorm:"DEFAULT:NULL"`
	StrID      int
	FileName   string
	State      string
	VideoCodec string
	Width      int
	Height     int
	FrameRate  float64
	AudioT     []Audio `gorm:"ForeignKey:VideoID"`
	SubtitleT  []Sub   `gorm:"ForeignKey:VideoID"`
}

type Audio struct {
	ID       int `json:"id,primary_key"`
	VideoID  uint
	StreamID int
	AtCodec  string
	Language string
	Channels int
}

type Sub struct {
	ID       int `json:"id,primary_key"`
	VideoID  uint
	StreamID int
	Language string
}

type Pdata struct {
	FileName string
	Save     bool
	Streams  []Stream
}

type Stream struct {
	VtId      int
	VidPreset string
	AudPreset string
	AudioT    []Audio
	SubtitleT []Sub
}

type Preset struct {
	ID         uint
	Name       string
	Type       int `gorm:"tinyint(1)`
	Resolution string
	Codec      string
	Bitrate    string
}

type VideoData struct {
	VideoStream []VideoStream
}

type VideoStream struct {
	Stream     bool
	StreamName string
	State      string
	Video      []Video
}

type VfNPrd struct {
	Pdata Pdata
	Video Video
	Err   error
}

type Data struct {
	Vidinfo    Vidinfo
	Vidpresets []Preset
	Audpresets []Preset
}
