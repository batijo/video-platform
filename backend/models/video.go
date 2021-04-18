package models

import (
	"time"
)

type Vstream struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint
	Name   string
	Video  []Video `gorm:"ForeignKey:VstreamID"`
}

type Video struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	Title       string  `json:"title"`
	Description string  `json:"description"`
	UserID      uint    `json:"user_id"`
	Public      bool    `json:"public"`
	VstreamID   uint    `gorm:"DEFAULT:NULL" json:"vstream_id"`
	StrID       int     `json:"str_id"`
	FileName    string  `json:"file_name"`
	State       string  `json:"state"` // Three possible states: not_transcoded, transcoding, transcoded
	VideoCodec  string  `json:"video_codec"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	FrameRate   float64 `json:"frame_rate"`
	AudioT      []Audio `gorm:"ForeignKey:VideoID" json:"audio_t"`
	SubtitleT   []Sub   `gorm:"ForeignKey:VideoID" json:"subtitle_t"`

	EncData Encode `gorm:"ForeignKey:VideoID" json:"enc_data"`
}

func (v *Video) ParseWithVidinfo(i Vidinfo) {
	v.FileName = i.FileName
	for _, t := range i.Videotrack {
		v.StrID = t.Index
		v.VideoCodec = t.CodecName
		v.Width = t.Width
		v.Height = t.Height
		v.FrameRate = t.FrameRate
	}
	for _, t := range i.Audiotrack {
		var at Audio
		at.StreamID = t.Index
		at.AtCodec = t.CodecName
		at.Language = t.Language
		at.Channels = t.Channels
		v.AudioT = append(v.AudioT, at)
	}
	for _, t := range i.Subtitle {
		var st Sub
		st.StreamID = t.Index
		st.Language = t.Language
		v.SubtitleT = append(v.SubtitleT, st)
	}
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
	ID      int  `json:"id,primary_key"`
	VideoID uint `json:"video_id"`
	EncID   uint `gorm:"DEFAULT:NULL" json:"enc_id"`

	StreamID int    `json:"stream_id"`
	AtCodec  string `json:"at_codec"`
	Language string `json:"language"`
	Channels int    `json:"channels"`
}

type Sub struct {
	ID      int  `json:"id,primary_key"`
	VideoID uint `json:"video_id"`
	EncID   uint `gorm:"DEFAULT:NULL" json:"enc_id"`

	StreamID int    `json:"stream_id"`
	Language string `json:"language"`
}

type Encode struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	VideoID   uint      `json:"video_id"`

	StrID      int     `json:"str_id"`
	FileName   string  `json:"file_name"`
	VideoCodec string  `json:"video_codec"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	FrameRate  float64 `json:"frame_rate"`
	AudioT     []Audio `gorm:"ForeignKey:EncID" json:"audio_t"`
	SubtitleT  []Sub   `gorm:"ForeignKey:EncID" json:"subtitle_t"`
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
	VideoID  uint
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
	Enc   Encode
	Err   error
}

type Data struct {
	Vidinfo    Vidinfo
	Vidpresets []Preset
	Audpresets []Preset
}
