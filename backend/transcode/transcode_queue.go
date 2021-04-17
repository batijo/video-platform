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

func AddToQueue(encData models.Encode) error {
	var (
		video models.Video
		resp  *gorm.DB
	)

	resp = utils.DB.Where("id = ?", encData.VideoID).First(&video)
	if resp.Error != nil {
		return resp.Error
	}

	video.EncData = encData
	resp = utils.DB.Save(&video)
	if resp.Error != nil {
		return resp.Error
	}

	if !active {
		go startTranscoder(video)
	}

	return nil
}

func startTranscoder(video models.Video) {
	active = true
	for {
		var clientData models.Video
		clientData.ParseWithEncode(video.EncData)
		clientData.State = video.State

		// Start procesing file
		go processVodFile(clientData, models.Pdata{}, "", video.ID, video.UserID)
		<-finished

		// Remove transcoded video from queue
		err := removeFromQueue(video.ID)
		if err != nil {
			active = false
			log.Panicln(err)
		}

		// Get new video id for transcoding if there is none, stop transcoder
		newVidId, err := nextInQueue()
		if err != nil {
			active = false
			log.Panicln(err)
		} else if newVidId < 0 {
			break
		}

		// Get new video for transcoding
		video, err = getVideo(uint(newVidId))
		if err != nil {
			active = false
			log.Panicln(err)
		}
	}
	active = false
}

// Returns full video structure
func getVideo(vidId uint) (models.Video, error) {
	var video models.Video

	err := utils.DB.Preload("AudioT").Preload("SubtitleT").Preload("EncData").Where("id = ?", vidId).First(&video).Error
	if err != nil {
		return video, err
	}

	return video, nil
}

// Removes already transcoded video from queue
func removeFromQueue(vidId uint) error {

	video, err := getVideo(vidId)
	if err != nil {
		return err
	}

	err = utils.DB.Delete(&models.Encode{}, video.EncData.ID).Error
	if err != nil {
		return err
	}

	return nil
}

// Retruns id of next video in queue, or -1 if no videos in queue
func nextInQueue() (int, error) {
	var encVideos []models.Encode

	resp := utils.DB.Find(&encVideos)
	if resp.Error != nil || len(encVideos) < 1 {
		return -1, resp.Error
	}
	sort.Sort(models.ByCreateDate(encVideos))

	return int(encVideos[0].VideoID), nil
}
