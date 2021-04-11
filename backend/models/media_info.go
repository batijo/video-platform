package models

import (
	"reflect"
	"strconv"
	"strings"
)

type videotrack struct {
	Index       int     `json:"index"`
	Duration    string  `json:"duration"`
	Width       int     `json:"width"`
	Height      int     `json:"height"`
	FrameRate   float64 `json:"frameRate"`
	CodecName   string  `json:"codecName"`
	AspectRatio string  `json:"aspectRatio"`
	FieldOrder  string  `json:"fieldOrder"`
}

type audiotrack struct {
	Index      int    `json:"index"`
	Channels   int    `json:"channels"`
	SampleRate int    `json:"sampleRate"`
	Language   string `json:"language"`
	BitRate    int    `json:"bitRate"`
	CodecName  string `json:"codecName"`
}

type subtitle struct {
	Index    int    `json:"index"`
	Language string `json:"language"`
}

// Vidinfo json struct with information about video file
type Vidinfo struct {
	Videotracks int          `json:"videotracks"`
	Audiotracks int          `json:"audiotracks"`
	Subtitles   int          `json:"subtitles"`
	Videotrack  []videotrack `json:"videotrack"`
	Audiotrack  []audiotrack `json:"audiotrack"`
	Subtitle    []subtitle   `json:"subtitle "`
}

// IsEmpty method which checks if Vidinfo is empty struct
func (s Vidinfo) IsEmpty() bool {
	return reflect.DeepEqual(s, Vidinfo{})
}

func (v *Vidinfo) ParseFFprobeData(out Ffprobe) {
	var (
		vc = 0
		ac = 0
		sc = 0
	)
	for _, s := range out.Streams {
		if s.CodecType == "video" {
			v.Videotrack = append(v.Videotrack, videotrack{})
			v.Videotrack[vc].Index = s.Index
			v.Videotrack[vc].CodecName = s.CodecName
			v.Videotrack[vc].Duration = s.Tags.Duration
			v.Videotrack[vc].Width = s.Width
			v.Videotrack[vc].Height = s.Height
			if s.RFrameRrate != "" {
				split := strings.Split(s.RFrameRrate, "/")
				fr, _ := strconv.ParseFloat(split[0], 64)
				sk, _ := strconv.ParseFloat(split[1], 64)
				v.Videotrack[vc].FrameRate = fr / sk
			} else {
				v.Videotrack[vc].FrameRate = 0
			}
			v.Videotrack[vc].AspectRatio = s.DisplayAspectRatio
			v.Videotrack[vc].FieldOrder = ""

			vc = vc + 1
		} else if s.CodecType == "audio" {
			v.Audiotrack = append(v.Audiotrack, audiotrack{})
			v.Audiotrack[ac].Index = s.Index
			if s.Tags.Language == "" || s.Tags.Language == "und" {
				v.Audiotrack[ac].Language = "undefined"
			} else {
				v.Audiotrack[ac].Language = s.Tags.Language
			}
			v.Audiotrack[ac].CodecName = s.CodecName
			v.Audiotrack[ac].Channels = s.Channels
			v.Audiotrack[ac].SampleRate, _ = strconv.Atoi(s.SampleRate)
			v.Audiotrack[ac].BitRate, _ = strconv.Atoi(s.BitRate)

			ac = ac + 1
		} else if s.CodecType == "srt" || s.CodecType == "subrip" {
			v.Subtitle = append(v.Subtitle, subtitle{})
			v.Subtitle[sc].Index = s.Index
			v.Subtitle[sc].Language = s.Tags.Language

			sc = sc + 1
		}
	}
	v.Videotracks = vc
	v.Audiotracks = ac
	v.Subtitles = sc
}
