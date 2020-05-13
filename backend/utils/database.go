package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Dzionys/video-platform/backend/models"

	"github.com/BurntSushi/toml"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres" //Gorm postgres dialect interface
	"github.com/joho/godotenv"
)

// DB database variable
var DB *gorm.DB

//ConnectDB function: Make database connection
func ConnectDB() *gorm.DB {

	//Load environmenatal variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	username := os.Getenv("databaseUser")
	password := os.Getenv("databasePassword")
	databaseName := os.Getenv("databaseName")
	databaseHost := os.Getenv("databaseHost")

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
		models.Video{},
		models.Audio{},
		models.Sub{},
	)

	db.Model(&models.Video{}).AddForeignKey("user_id", "users(id)", "SET NULL", "NO ACTION")
	db.Model(&models.Audio{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Sub{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")

	return db
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
