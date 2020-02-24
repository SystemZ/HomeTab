package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/systemz/gotag/config"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func StartWebInterface() {
	// create multiple routes
	r := mux.NewRouter()

	// JSON API for JS frontend
	r.HandleFunc("/api/v1/login", Login).Methods("POST")
	r.HandleFunc("/api/v1/scan", Scan).Methods("POST")
	r.HandleFunc("/api/v1/files", FilePaginate).Methods("POST")
	r.HandleFunc("/api/v1/file/{sha256}/similar", FileSimilar).Methods("GET")
	r.HandleFunc("/api/v1/file/{sha256}/tag/delete", TagDelete).Methods("POST")
	r.HandleFunc("/api/v1/file/{sha256}/tag/add", TagAdd).Methods("POST")
	r.HandleFunc("/api/v1/tags", TagList).Methods("GET")
	r.HandleFunc("/api/v1/file/tags", TagListForFiles).Methods("POST")

	// no-JSON zone
	r.HandleFunc("/img/thumbs/{w}/{h}/{sha256}", Thumb).Methods("GET")
	r.HandleFunc("/img/full/{sha256}", FullImg).Methods("GET")

	// serve frontend
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	dir += "/frontend"
	log.Printf("Serving static content from %v", dir)
	// TODO check security of this
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, dir+"/index.html")
	})
	r.PathPrefix("/css/").Handler(
		http.StripPrefix("/css/",
			http.FileServer(
				http.Dir(dir+"/css"),
			),
		),
	)
	r.PathPrefix("/js/").Handler(
		http.StripPrefix("/js/",
			http.FileServer(
				http.Dir(dir+"/js"),
			),
		),
	)
	r.PathPrefix("/fonts/").Handler(
		http.StripPrefix("/fonts/",
			http.FileServer(
				http.Dir(dir+"/fonts"),
			),
		),
	)

	// start internal http server with logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("HTTP server started on :4000")

	if config.DEV_MODE {
		log.Println("DEV env detected, CORS wildcard")
		loggedRouter = cors.AllowAll().Handler(loggedRouter)
	}

	log.Fatal(http.ListenAndServe(":4000", loggedRouter))
}
