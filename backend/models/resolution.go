package models

import (
	"errors"
	"strconv"
	"strings"
)

var prRes = map[string]string{
	"240p":  "352x240",
	"576p":  "720x576",
	"720p":  "1280x720",
	"360p":  "480x360",
	"1080p": "1920x1080",
}

func PresetResolution(res string) string {
	return prRes[res]
}

func GetPresetWidth(res string) int {
	hw, err := getWidthHeight(res)
	if err != nil {
		return 0
	}
	return hw[1]
}

func GetPresetHeight(res string) int {
	hw, err := getWidthHeight(res)
	if err != nil {
		return 0
	}
	return hw[0]
}

func getWidthHeight(res string) ([]int, error) {
	var hwInt []int
	hw := strings.Split(prRes[res], "x")
	if len(hw) != 2 {
		return hwInt, errors.New("resolution not found")
	}
	for _, e := range hw {
		eint, err := strconv.Atoi(e)
		if err != nil {
			return hwInt, err
		}
		hwInt = append(hwInt, eint)
	}
	return hwInt, nil
}
