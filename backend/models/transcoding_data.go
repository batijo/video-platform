package models

type PData struct {
	FileName string
	Save     bool
	Streams  []Streams
}

type Streams struct {
	VtId      int
	VidPreset string
	AudPreset string
	AudioT    []AudT
	SubtitleT []SubT
}

type AudT struct {
	AtId int
	Lang string
}

type SubT struct {
	StId int
	Lang string
}

type Preset struct {
	ID         uint
	Name       string
	Type       int `gorm:"tinyint(1)`
	Resolution string
	Codec      string
	Bitrate    string
}
