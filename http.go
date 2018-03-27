package main

import (
	"encoding/json"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"gitlab.systemz.pl/systemz/tasktab/model"
	"io"
	"log"
	"net/http"
	"os"
)

func respondOk(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	respondOk(w)
	// In the future we could report back on the status of our DB, or our cache
	// (e.g. Redis) by performing a simple PING, and include them in the response.
	io.WriteString(w, `{"alive": true}`)
}

func taskListHandler(w http.ResponseWriter, r *http.Request) {
	res := model.ListTasksForGroup(1)

	//pagesJson, err := json.Marshal(res)
	pagesJson, err := json.MarshalIndent(res,"","\t")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": true}`)
	}
	respondOk(w)
	io.WriteString(w, string(pagesJson))
}

func tasksTodoListHandler(w http.ResponseWriter, r *http.Request) {
	res := model.ListTasksToDoForGroup(1)

	//pagesJson, err := json.Marshal(res)
	pagesJson, err := json.MarshalIndent(res,"","\t")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, `{"error": true}`)
	}
	respondOk(w)
	io.WriteString(w, string(pagesJson))
}

func syncHandler(w http.ResponseWriter, r *http.Request) {
	go getTasksForAllGroups()
	respondOk(w)
	io.WriteString(w, `{"syncing": true}`)
}

func httpStart() {
	//FIXME configurable port
	host := ":3000"
	log.Println("HTTP server started on " + host)

	r := mux.NewRouter()
	r.HandleFunc("/health", healthCheckHandler).Methods("GET")
	r.HandleFunc("/api/v1/sync", syncHandler).Methods("GET")
	r.HandleFunc("/api/v1/tasks/all/{uid}", taskListHandler).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/v1/tasks/todo/{uid}", tasksTodoListHandler).Methods("GET", "OPTIONS")

	//http.Handle("/static", http.FileServer(http.Dir("./frontend/dist/static")))
	//http.Handle("/", http.FileServer(rice.MustFindBox("frontend").HTTPBox()))

	//box := rice.MustFindBox("frontend")
	//cssFileServer := http.StripPrefix("/dist/", http.FileServer(box.HTTPBox()))
	//http.Handle("/static/", cssFileServer)

	http.Handle("/", r)

	loggedRouter := handlers.LoggingHandler(os.Stdout, r)

	//corsObj:=handlers.AllowedOrigins([]string{"*"})
	//log.Fatal(http.ListenAndServe(host, handlers.CORS(corsObj)(loggedRouter)))

	handler := cors.Default().Handler(loggedRouter)

	log.Fatal(http.ListenAndServe(host, handler))

	// for https visit: https://gist.github.com/denji/12b3a568f092ab951456
}
