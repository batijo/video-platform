package main

import (
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"net/http"

	"github.com/Dzionys/video-platform/backend/controllers"
	"github.com/Dzionys/video-platform/backend/utils"
	"github.com/Dzionys/video-platform/backend/utils/auth"
)

var (
	config utils.Config
)

func handlers() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	r := router.PathPrefix("/api").Subrouter()
	r.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")
	r.HandleFunc("/ngx/mapping/{name}", controllers.NginxMappingHandler).Methods("GET")

	// Auth route
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)
	s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	s.HandleFunc("/user/{id}", controllers.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/{id}", controllers.DeleteUser).Methods("DELETE")
	s.HandleFunc("/video", controllers.FetchVideos).Methods("GET")
	s.HandleFunc("/video/{id}", controllers.DeleteVideo).Methods("DELETE")
	s.HandleFunc("/video/{id}", controllers.UpdateVideo).Methods("PUT")
	s.HandleFunc("/video/{id}", controllers.GetVideo).Methods("GET")
	s.HandleFunc("/upload", controllers.VideoUpload).Methods("POST")
	s.HandleFunc("/upload", controllers.TranscodeHandler).Methods("PUT")
	s.HandleFunc("/tc", controllers.TcTypeHandler).Methods("POST")
	s.HandleFunc("/list", controllers.ListHandler).Methods("GET")

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}

func main() {

	// Make a new Broker instance
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
	config, err = utils.GetConf()
	if err != nil {
		log.Println("Error: failed to load config file")
		log.Println(err)
		return
	}
	utils.Conf = config

	// Write all logs to file
	err = utils.OpenLogFile(config.LogP)
	if err != nil {
		log.Println("Error: failed open log file")
		log.Panicln(err)
		return
	}
	//defer utils.LogFile.Close()

	// Connect to database
	utils.DB = utils.ConnectDB()
	defer utils.DB.Close()

	// Insert presets to database
	err = utils.InsertPresets()
	if err != nil {
		log.Println(err)
	}

	// Load .env file
	err = godotenv.Load()
	if err != nil {
		log.Println("Error: failed to load .env file")
		log.Println(err)
		return
	}
	port := os.Getenv("PORT")

	// Handle routes
	http.Handle("/", handlers())

	// serve
	log.Printf("Server up on port '%s'", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
