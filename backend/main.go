package main

import (
	"log"
	"os"

	"net/http"

	"github.com/batijo/video-platform/backend/routes"
	"github.com/batijo/video-platform/backend/utils"
)

func main() {

	// Make a new Broker instance for SSE
	utils.B = &utils.Broker{
		Clients:        make(map[chan string]string),
		NewClients:     make(chan utils.Client),
		DefunctClients: make(chan (chan string)),
		Messages:       make(chan utils.Message),
	}

	// Start processing events
	utils.B.Start()

	// Load config file
	var err error
	config, err := utils.GetConf()
	if err != nil {
		log.Println("Error: failed to load config file")
		log.Println(err)
		return
	}
	utils.Conf = config

	// Connect to database
	utils.DB = utils.ConnectDB()
	defer utils.DB.Close()

	// Insert presets to database
	err = utils.InsertPresets()
	if err != nil {
		log.Println(err)
	}

	port := os.Getenv("PORT")

	// Handle routes
	http.Handle("/", routes.Handlers())

	// serve
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
