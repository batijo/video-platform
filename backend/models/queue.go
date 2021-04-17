package models

import (
	"sort"
)

type Queue struct {
	Elements []queueElement `json:"elements"`
}

type queueElement struct {
	Position   uint   `json:"position"`
	VideoTitle string `json:"video_title"`
	Owns       bool   `json:"owns"`
}

func (q *Queue) Put(vids []Video, userID uint) {
	var (
		enc []Encode
	)
	for _, v := range vids {
		if v.EncData.VideoID == v.ID {
			enc = append(enc, v.EncData)
			if !v.Public {
				enc[len(enc)-1].FileName = "unknown"
			} else {
				enc[len(enc)-1].FileName = v.Title
			}
			if v.UserID == userID {
				enc[len(enc)-1].ID = 1
			} else {
				enc[len(enc)-1].ID = 0
			}
		}
	}
	sort.Sort(ByCreateDate(enc))

	for i, e := range enc {
		(*q).Elements = append((*q).Elements,
			queueElement{
				Position:   uint(i + 1),
				VideoTitle: e.FileName,
				Owns:       btoi(e.ID),
			},
		)
	}

}

func btoi(i uint) bool {
	if i > 0 {
		return true
	}
	return false
}
