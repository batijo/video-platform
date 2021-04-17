package routes

import (
	"net/http"

	"github.com/batijo/video-platform/backend/controllers"
	"github.com/batijo/video-platform/backend/utils"
	"github.com/batijo/video-platform/backend/utils/auth"
	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)
	router.Use(commonMiddleware)

	r := router.PathPrefix("/api").Subrouter()
	r.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	r.Handle("/sse/dashboard", utils.B)
	// r.HandleFunc("/ngx/mapping/{name}", controllers.NginxMappingHandler).Methods("GET")

	// Auth route
	s := r.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)

	s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")
	s.HandleFunc("/user/{id}", controllers.GetUser).Methods("GET")
	s.HandleFunc("/user/update/{id}", controllers.UpdateUser).Methods("POST")
	s.HandleFunc("/user/delete/{id}", controllers.DeleteUser).Methods("POST")

	s.HandleFunc("/video", controllers.FetchVideos).Methods("GET")
	s.HandleFunc("/video/{id}", controllers.GetVideo).Methods("GET")
	s.HandleFunc("/video/update/{id}", controllers.UpdateVideo).Methods("POST")
	s.HandleFunc("/video/delete/{id}", controllers.DeleteVideo).Methods("POST")

	s.HandleFunc("/upload", controllers.VideoUpload).Methods("POST")
	s.HandleFunc("/transcode/{id}", controllers.TranscodeHandler).Methods("POST")

	a := r.PathPrefix("/admin").Subrouter()
	a.Use(auth.AdminVerify)

	// s.HandleFunc("/tc", controllers.TcTypeHandler).Methods("POST")
	// s.HandleFunc("/list", controllers.ListHandler).Methods("GET")

	return router
}

func commonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		// Allow CORS
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
		w.Header().Set(
			"Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})
}