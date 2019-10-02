package web

import (
	"log"
	"net/http"
)

func StartWebInterface() {
	log.Printf("Starting web interface...")
	http.HandleFunc("/login", Login)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
	"a":     "b",
}
