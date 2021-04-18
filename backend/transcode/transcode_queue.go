package transcode

import (
	"log"
	"sort"

	"github.com/batijo/video-platform/backend/models"
	"github.com/batijo/video-platform/backend/utils"
	"github.com/jinzhu/gorm"
)

var (
	finished = make(chan bool)
	active   bool
)

func AddToQueue(enc models.Encode, prData []models.Stream, videoID uint) error {
	var (
		encData models.Encodedata
		video   models.Video
		resp    *gorm.DB
	)

	if len(encData.Presets) > 0 {
		resp = utils.DB.Where("id = ?", videoID).First(&video)
		if resp.Error != nil {
			return resp.Error
		}

		encData.Presets = append(encData.Presets, prData...)
		encData.Video = video
		resp = utils.DB.Save(&encData)
		if resp.Error != nil {
			return resp.Error
		}

	} else {
		resp = utils.DB.Where("id = ?", videoID).First(&video)
		if resp.Error != nil {
			return resp.Error
		}

		encData.EncData = enc
		encData.Video = video
		resp = utils.DB.Save(&encData)
		if resp.Error != nil {
			return resp.Error
		}
	}

	if !active {
		go startTranscoder(encData)
	}

	return nil
}

func Active() bool {
	return active
}

func startTranscoder(ED models.Encodedata) {
	active = true
	for {
		// Start procesing file
		go processVodFile(ED)
		<-finished

		// Remove transcoded video from queue
		err := removeFromQueue(ED.ID)
		if err != nil {
			active = false
			log.Panicln(err)
		}

		// Get new video id for transcoding if there is none, stop transcoder
		newEdID, err := nextInQueue()
		if err != nil {
			active = false
			log.Panicln(err)
		} else if newEdID < 0 {
			break
		}

		// Get new video for transcoding
		ED, err = getEncData(uint(newEdID))
		if err != nil {
			active = false
			log.Panicln(err)
		}
	}
	active = false
}

// Returns full video structure
func getEncData(edID uint) (models.Encodedata, error) {
	var ED models.Encodedata

	err := utils.DB.Preload("Video").Preload("EncData").Preload("Presets").Where("id = ?", edID).First(&ED).Error
	if err != nil {
		return ED, err
	}

	return ED, nil
}

// Removes already transcoded video from queue
func removeFromQueue(edID uint) error {
	err := utils.DB.Where("id = ?", edID).Delete(&models.Encodedata{}).Error
	if err != nil {
		return err
	}

	return nil
}

// Retruns id of next video in queue, or -1 if no videos in queue
func nextInQueue() (int, error) {
	var encData []models.Encodedata

	resp := utils.DB.Find(&encData)
	if resp.Error != nil || len(encData) < 1 {
		return -1, resp.Error
	}
	sort.Sort(models.ByCreateDate(encData))

	return int(encData[0].ID), nil
}
