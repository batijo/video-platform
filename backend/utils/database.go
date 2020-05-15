package utils

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Dzionys/video-platform/backend/models"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres dialect interface
)

// DB database variable
var DB *gorm.DB

//ConnectDB function: Make database connection
func ConnectDB() *gorm.DB {
	username := os.Getenv("DATABASE_USER")
	password := os.Getenv("DATABASE_PASSWORD")
	databaseName := os.Getenv("DATABASE_NAME")
	databaseHost := os.Getenv("DATABASE_HOST")

	//Define DB connection string
	dbURI := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s", databaseHost, username, databaseName, password)

	//connect to db URI
	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		fmt.Println("error", err)
		panic(err)
	}

	db.DropTableIfExists(&models.Preset{})

	// Migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.Preset{},
		models.Vstream{},
		models.Video{},
		models.Audio{},
		models.Sub{},
	)

	db.Model(&models.Video{}).AddForeignKey("user_id", "users(id)", "SET NULL", "NO ACTION")
	db.Model(&models.Video{}).AddForeignKey("vstream_id", "vstreams(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Vstream{}).AddForeignKey("user_id", "users(id)", "SET NULL", "NO ACTION")
	db.Model(&models.Audio{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Sub{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")

	return db
}

// InsertVideo adds video to database
func InsertVideo(vidinfo models.Vidinfo, name string, state string, userID uint, streamID int) error {

	var (
		user     models.User
		audio    []models.Audio
		subtitle []models.Sub
		video    models.Video
	)

	DB.First(&user, userID)

	for _, a := range vidinfo.Audiotrack {
		at := models.Audio{
			StreamID: a.Index,
			AtCodec:  a.CodecName,
			Language: a.Language,
			Channels: a.Channels,
		}
		audio = append(audio, at)
	}

	for _, s := range vidinfo.Subtitle {
		st := models.Sub{
			StreamID: s.Index,
			Language: s.Language,
		}
		subtitle = append(subtitle, st)
	}

	if streamID < 0 {
		video = models.Video{
			StrID:      vidinfo.Videotrack[0].Index,
			FileName:   name,
			State:      state,
			VideoCodec: vidinfo.Videotrack[0].CodecName,
			Width:      vidinfo.Videotrack[0].Width,
			Height:     vidinfo.Videotrack[0].Height,
			FrameRate:  vidinfo.Videotrack[0].FrameRate,
			AudioT:     audio,
			SubtitleT:  subtitle,
		}
	} else {
		video = models.Video{
			VstreamID:  uint(streamID),
			StrID:      vidinfo.Videotrack[0].Index,
			FileName:   name,
			State:      state,
			VideoCodec: vidinfo.Videotrack[0].CodecName,
			Width:      vidinfo.Videotrack[0].Width,
			Height:     vidinfo.Videotrack[0].Height,
			FrameRate:  vidinfo.Videotrack[0].FrameRate,
			AudioT:     audio,
			SubtitleT:  subtitle,
		}
	}

	user.Video = append(user.Video, video)

	createdVideo := DB.Save(&user)

	if createdVideo.Error != nil {
		return createdVideo.Error
	}

	return nil
}

// DeleteVideo deletes video based on its name
func DeleteVideo(name string) error {
	var video models.Video
	if err := DB.Where("file_name = ?", name).Delete(&video).Error; err != nil {
		return err
	}
	return nil
}

// DeleteStream deletes stream based on its name
func DeleteStream(name string) error {
	var stream models.Vstream
	if err := DB.Where("name = ?", name).Delete(&stream).Error; err != nil {
		return err
	}
	return nil
}

// InsertStream ...
func InsertStream(ndata []models.Vidinfo, names []string, state string, sname string, userID uint) {

	var (
		user   models.User
		stream models.Vstream
	)

	DB.First(&user, userID)

	stream = models.Vstream{
		Name: sname,
	}
	user.Stream = append(user.Stream, stream)
	DB.Save(&user)

	DB.Where("name = ?", sname).First(&stream)

	for i, video := range ndata {
		InsertVideo(video, names[i], state, userID, int(stream.ID))
	}

}

// AddPresetsToJSON ...
func AddPresetsToJSON(vid models.Vidinfo) models.Data {

	presets := getPresets()
	var data models.Data

	for _, p := range presets {
		if p.Type == 0 {
			data.Vidpresets = append(data.Vidpresets, p)
		} else {
			data.Audpresets = append(data.Audpresets, p)
		}
	}
	data.Vidinfo = vid

	return data
}

func getPresets() []models.Preset {
	var presets []models.Preset
	DB.Find(&presets)

	return presets
}

// GetPreset returns Preset data from database with given name
func GetPreset(name string) (models.Preset, error) {
	presets := getPresets()

	for _, p := range presets {
		if p.Name == name {
			return p, nil
		}
	}

	return models.Preset{}, fmt.Errorf("error: no preset found with given name")
}

// InsertPresets ...
func InsertPresets() error {
	var presets struct {
		PresetValues [][]string
	}
	if _, err := toml.DecodeFile("utils/preset_values.toml", &presets); err != nil {
		return err
	}

	for _, p := range presets.PresetValues {

		tp, err := strconv.Atoi(p[1])
		if err != nil {
			return err
		}
		preset := &models.Preset{
			Name:       p[0],
			Type:       tp,
			Resolution: p[2],
			Codec:      p[3],
			Bitrate:    p[4],
		}

		prst := DB.Create(preset)
		if prst.Error != nil {
			return prst.Error
		}
	}

	return nil
}
