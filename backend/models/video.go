package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type Vstream struct {
	gorm.Model

	UserID uint
	Name   string
	Video  []Video `gorm:"ForeignKey:VstreamID"`
}

type Video struct {
	gorm.Model

	Title       string
	Description string
	UserID      uint
	Public      bool
	VstreamID   uint `gorm:"DEFAULT:NULL"`
	StrID       int
	FileName    string
	State       string
	VideoCodec  string
	Width       int
	Height      int
	FrameRate   float64
	AudioT      []Audio `gorm:"ForeignKey:VideoID"`
	SubtitleT   []Sub   `gorm:"ForeignKey:VideoID"`

	EncData Encode `gorm:"ForeignKey:VideoID"`
}

func (v *Video) ParseWithEncode(e Encode) {
	v.StrID = e.StrID
	v.FileName = e.FileName
	v.VideoCodec = e.VideoCodec
	v.Width = e.Width
	v.Height = e.Height
	v.FrameRate = e.FrameRate
	for _, at := range e.AudioT {
		v.AudioT = append(v.AudioT, at)
	}
	for _, st := range e.SubtitleT {
		v.SubtitleT = append(v.SubtitleT, st)
	}
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

type Encode struct {
	ID        uint `gorm:"primary_key"`
	CreatedAt time.Time
	VideoID   uint

	StrID      int
	FileName   string
	VideoCodec string
	Width      int
	Height     int
	FrameRate  float64
	AudioT     []Audio `gorm:"ForeignKey:VideoID"`
	SubtitleT  []Sub   `gorm:"ForeignKey:VideoID"`
}

type ByCreateDate []Encode

// Forward request for length
func (p ByCreateDate) Len() int {
	return len(p)
}

// Define compare
func (p ByCreateDate) Less(i, j int) bool {
	return p[i].CreatedAt.Before(p[j].CreatedAt)
}

// Define swap over an array
func (p ByCreateDate) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
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
	Type       int `gorm:"tinyint(1)"`
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
