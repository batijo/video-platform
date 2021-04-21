package models

import (
	"time"
)

type Encodedata struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`

	EncData Encode   `gorm:"ForeignKey:QueueID" json:"enc_data"`
	Presets []Stream `gorm:"ForeignKey:QueueID" json:"presets"`

	Video Video `gorm:"ForeignKey:QueueID" json:"video"`
}

type Vstream struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	UserID uint    `json:"user_id"`
	Name   string  `json:"name"`
	Video  []Video `gorm:"ForeignKey:VstreamID" json:"video"`
}

type Video struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	QueueID   uint      `gorm:"DEFAULT:NULL" json:"queue_id"`
	UserID    uint      `json:"user_id"`
	VstreamID uint      `gorm:"DEFAULT:NULL" json:"vstream_id"`

	Title       string  `json:"title"`
	Description string  `json:"description"`
	Public      bool    `json:"public"`
	StrID       int     `json:"str_id"`
	FileName    string  `json:"file_name"`
	State       string  `json:"state"` // Possible states: "not_transcoded", "transcoding", "transcoded"
	VideoCodec  string  `json:"video_codec"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	FrameRate   float64 `json:"frame_rate"`
	Duration    float64 `json:"duration"`
	AudioT      []Audio `gorm:"ForeignKey:VideoID" json:"audio_t"`
	SubtitleT   []Sub   `gorm:"ForeignKey:VideoID" json:"subtitle_t"`
}

type Audio struct {
	ID      int  `json:"id,primary_key"`
	VideoID uint `gorm:"DEFAULT:NULL" json:"video_id"`
	EncID   uint `gorm:"DEFAULT:NULL" json:"enc_id"`
	StrID   uint `gorm:"DEFAULT:NULL" json:"str_id"`

	StreamID int    `json:"stream_id"`
	AtCodec  string `json:"at_codec"`
	Language string `json:"language"`
	Channels int    `json:"channels"`
}

type Sub struct {
	ID      int  `json:"id,primary_key"`
	VideoID uint `gorm:"DEFAULT:NULL" json:"video_id"`
	EncID   uint `gorm:"DEFAULT:NULL" json:"enc_id"`
	StrID   uint `gorm:"DEFAUT:NULL" json:"str_id"`

	StreamID int    `json:"stream_id"`
	Language string `json:"language"`
}

type Encode struct {
	ID        uint      `gorm:"primary_key" json:"id"`
	CreatedAt time.Time `json:"created_at"`
	QueueID   uint      `json:"queue_id"`

	StrID      int     `json:"str_id"`
	FileName   string  `json:"file_name"`
	VideoCodec string  `json:"video_codec"`
	Width      int     `json:"width"`
	Height     int     `json:"height"`
	FrameRate  float64 `json:"frame_rate"`
	AudioT     []Audio `gorm:"ForeignKey:EncID" json:"audio_t"`
	SubtitleT  []Sub   `gorm:"ForeignKey:EncID" json:"subtitle_t"`
}

type Stream struct {
	ID      uint `json:"id"`
	QueueID uint `json:"queue_id"`

	VidPreset string  `json:"vid_preset"`
	AudPreset string  `json:"aud_preset"`
	VtId      int     `json:"vt_id"`
	AudioT    []Audio `gorm:"ForeignKey:StrID" json:"audio_t"`
	SubtitleT []Sub   `gorm:"ForeignKey:StrID" json:"subtitle_t"`
}

type Preset struct {
	ID uint `json:"id"`

	Name       string `json:"name"`
	Type       int    `gorm:"tinyint(1)" json:"type"`
	Resolution string `json:"resolution"`
	Codec      string `json:"codec"`
	Bitrate    string `json:"bitrate"`
}

func (v *Video) ParseWithPreset(vp, ap Preset, frameRate float64, vtID, atID int, atLang string) {
	v.StrID = vtID
	v.VideoCodec = vp.Codec
	v.Width = GetPresetWidth(vp.Resolution)
	v.Height = GetPresetHeight(vp.Resolution)
	v.FrameRate = frameRate

	at := Audio{
		StreamID: atID,
		AtCodec:  ap.Codec,
		Language: atLang,
	}
	v.AudioT = append(v.AudioT, at)
}

func (v *Video) ParseWithVidinfo(i Vidinfo) {
	v.FileName = i.FileName
	for _, t := range i.Videotrack {
		v.StrID = t.Index
		v.VideoCodec = t.CodecName
		v.Width = t.Width
		v.Height = t.Height
		v.FrameRate = t.FrameRate
		v.Duration = t.Duration
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

func (v *Video) ParseWithEncode(e Encode, state string) {
	v.StrID = e.StrID
	v.FileName = e.FileName
	v.VideoCodec = e.VideoCodec
	v.Width = e.Width
	v.Height = e.Height
	v.FrameRate = e.FrameRate
	v.State = state
	for _, at := range e.AudioT {
		v.AudioT = append(v.AudioT, at)
	}
	for _, st := range e.SubtitleT {
		v.SubtitleT = append(v.SubtitleT, st)
	}
}

type ByCreateDate []Encodedata

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

// ====================================================== //

type VideoData struct {
	VideoStream []VideoStream
}

type VideoStream struct {
	Stream     bool
	StreamName string
	State      string
	Video      []Video
}

type Presets struct {
	Video      Video
	Vidpresets []Preset
	Audpresets []Preset
}
