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

func (q *Queue) Put(ED []Encodedata, userID uint) {
	sort.Sort(ByCreateDate(ED))

	for i, e := range ED {
		(*q).Elements = append((*q).Elements,
			queueElement{
				Position:   uint(i + 1),
				VideoTitle: getTile(e.Video, userID),
				Owns:       btoi(e.ID),
			},
		)
	}
}

func getTile(v Video, userID uint) string {
	if !v.Public || v.UserID != userID {
		return "unknown"
	}
	return v.Title
}

func btoi(i uint) bool {
	if i > 0 {
		return true
	}
	return false
}
