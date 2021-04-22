package main

import (
	"log"
	"os"

	"net/http"

	"github.com/batijo/video-platform/backend/routes"
	"github.com/batijo/video-platform/backend/utils"
)

func main() {

	// Start processing events
	utils.NewSseServer()

	// Load config file
	if err := utils.Conf.Load(); err != nil {
		log.Println("Error: failed to load config file")
		log.Println(err)
		return
	}

	// Connect to database
	utils.DB = utils.ConnectDB()
	defer utils.DB.Close()

	// Insert presets to database
	if err := utils.InsertPresets(); err != nil {
		log.Println(err)
		return
	}

	// Initialize redis connection
	utils.InitRedisClient(os.Getenv("REDIS_PORT"), os.Getenv("REDIS_PASSWORD"))

	if err := utils.CreateSuperUser(os.Getenv("SU_EMAIL"), os.Getenv("SU_PASS"), os.Getenv("SU_USERNAME")); err != nil {
		log.Panicln(err)
	}

	// Handle routes
	http.Handle("/", routes.SetupRoutes())

	// serve
	log.Printf("Server up on port '%s'", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":"+os.Getenv("PORT"), nil))
}
