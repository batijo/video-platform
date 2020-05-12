package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/Dzionys/video-platform/backend/models"

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

	// Migrate the schema
	db.AutoMigrate(
		&models.User{},
		&models.Video{},
		models.Audio{},
		models.Sub{},
	)

	db.Model(&models.Audio{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")
	db.Model(&models.Sub{}).AddForeignKey("video_id", "videos(id)", "CASCADE", "NO ACTION")

	return db
}
