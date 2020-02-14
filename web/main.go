package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.com/systemz/gotag/config"
	"log"
	"net/http"
	"os"
)

func StartWebInterface() {
	// create multiple routes
	r := mux.NewRouter()

	// JSON API for JS frontend
	//r.HandleFunc("/api/v1/login", ApiLogin)
	r.HandleFunc("/api/v1/files", FilePaginate).Methods("POST")

	// no-JSON zone
	r.HandleFunc("/img/thumbs/{w}/{h}/{sha256}", Thumb).Methods("GET")
	r.HandleFunc("/img/full/{sha256}", FullImg).Methods("GET")

	// start internal http server with logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("HTTP server started on :4000")

	if config.DEV_MODE {
		log.Println("DEV env detected, CORS wildcard")
		loggedRouter = cors.AllowAll().Handler(loggedRouter)
	}

	log.Fatal(http.ListenAndServe(":4000", loggedRouter))
}
