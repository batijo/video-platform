package utils

import (
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/batijo/video-platform/backend/models"
	"golang.org/x/crypto/bcrypt"

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
	dbURI := fmt.Sprintf(
		"host=%s user=%s dbname=%s sslmode=disable password=%s",
		databaseHost,
		username,
		databaseName,
		password,
	)

	//connect to db URI
	db, err := gorm.Open("postgres", dbURI)

	if err != nil {
		log.Panicln(err)
	}

	db.DropTableIfExists(&models.Preset{})

	// Migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.Preset{},
		&models.Encodedata{},
		models.Vstream{},
		models.Video{},
		models.Audio{},
		models.Sub{},
		models.Encode{},
		models.Stream{},
	)

	// Add foreign keys
	db.Model(&models.Video{}).AddForeignKey("user_id", "users(id)", "SET NULL", "NO ACTION")
	db.Model(&models.Video{}).AddForeignKey("vstream_id", "vstreams(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Vstream{}).AddForeignKey("user_id", "users(id)", "SET NULL", "NO ACTION")
	db.Model(&models.Audio{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Sub{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")
	// Encode data
	db.Model(&models.Encode{}).AddForeignKey("queue_id", "encodedata(id)", "CASCADE", "CASCADE")
	db.Model(&models.Stream{}).AddForeignKey("queue_id", "encodedata(id)", "CASCADE", "CASCADE")
	db.Model(&models.Video{}).AddForeignKey("queue_id", "encodedata(id)", "SET NULL", "SET NULL")
	db.Model(&models.Audio{}).AddForeignKey("enc_id", "encodes(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Audio{}).AddForeignKey("str_id", "streams(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Sub{}).AddForeignKey("enc_id", "encodes(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Sub{}).AddForeignKey("str_id", "streams(id)", "CASCADE", "NO ACTION")

	return db
}

// Create admin account
func CreateSuperUser(email, pass, username string) error {
	re := regexp.MustCompile(
		"^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$",
	)
	if re.MatchString(email) {
		return errors.New(fmt.Sprint("email adress is not valid: ", email))
	} else if len(pass) < 6 {
		return errors.New("password must be at least 5 characters")
	} else if username == "" {
		return errors.New("username must not be empty")
	}

	hpass, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	var admin = models.User{
		Username: username,
		Email:    email,
		Password: string(hpass),
		Admin:    true,
	}

	if err := DB.Create(&admin).Error; err != nil {
		return err
	}

	return nil
}

// InsertVideo adds video to database
func InsertVideo(vidinfo models.Vidinfo, state string, userID uint, streamID int) (uint, error) {
	var (
		user     models.User
		audio    []models.Audio
		subtitle []models.Sub
		video    models.Video
	)

	if err := DB.First(&user, userID).Error; err != nil {
		return 0, err
	}

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

	if len(vidinfo.Videotrack) < 1 {
		return 0, errors.New("InserVideo: length of video track cannot be zero")
	}
	video = models.Video{
		StrID:      vidinfo.Videotrack[0].Index,
		FileName:   vidinfo.FileName,
		State:      state,
		VideoCodec: vidinfo.Videotrack[0].CodecName,
		Width:      vidinfo.Videotrack[0].Width,
		Height:     vidinfo.Videotrack[0].Height,
		FrameRate:  vidinfo.Videotrack[0].FrameRate,
		Duration:   vidinfo.Videotrack[0].Duration,
		AudioT:     audio,
		SubtitleT:  subtitle,
	}

	if streamID > 0 {
		video.VstreamID = uint(streamID)
	}

	user.Video = append(user.Video, video)

	if err := DB.Save(&user).Error; err != nil {
		return 0, err
	}

	return user.Video[0].ID, nil
}

// Unfinished
func UpdateVideo(id uint, updatedVideo models.Video) error {
	var video = models.Video{ID: id}
	if err := DB.Model(&video).Update(updatedVideo).Error; err != nil {
		return err
	}
	return nil
}

// DeleteVideo deletes video based on its name
func DeleteVideo(id uint) error {
	var video models.Video
	if err := DB.Where("id = ?", id).Delete(&video).Error; err != nil {
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
func InsertStream(ndata []models.Vidinfo, names []string, state string, sname string, userID uint) error {

	var (
		user   models.User
		stream models.Vstream
	)

	if err := DB.First(&user, userID).Error; err != nil {
		return err
	}

	stream = models.Vstream{
		Name: sname,
	}
	user.Stream = append(user.Stream, stream)
	if err := DB.Save(&user).Error; err != nil {
		return err
	}

	if err := DB.Where("name = ?", sname).First(&stream).Error; err != nil {
		return err
	}

	for i, video := range ndata {
		video.FileName = names[i]
		_, err := InsertVideo(video, state, userID, int(stream.ID))
		if err != nil {
			return err
		}
	}

	return nil
}

// GetPresetsWithData ...
func GetPresetsWithData(vid models.Video) (models.Presets, error) {
	var (
		presets []models.Preset
		data    models.Presets
	)

	if err := DB.Find(&presets).Error; err != nil {
		return data, err
	}
	for _, p := range presets {
		if p.Type == 0 {
			data.Vidpresets = append(data.Vidpresets, p)
		} else {
			data.Audpresets = append(data.Audpresets, p)
		}
	}
	data.Video = vid

	return data, nil
}

// GetPreset returns Preset data from database with given name
func GetPreset(name string) (models.Preset, error) {
	var preset models.Preset

	if err := DB.Where("name = ?", name).First(&preset).Error; err != nil {
		return preset, err
	}

	return preset, nil
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
