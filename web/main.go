package web

import (
	"github.com/unrolled/render"
	"log"
	"net/http"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
	"a":     "b",
}
var (
	display *render.Render
)

func init() {
	display = render.New(render.Options{
		IndentJSON: true,
		Extensions: []string{".html"},
	})
}

func StartWebInterface() {
	log.Printf("Starting web interface...")
	http.HandleFunc("/", Index)
	http.HandleFunc("/login", Login)
	http.HandleFunc("/welcome", Welcome)
	http.HandleFunc("/refresh", Refresh)
	http.HandleFunc("/tasks", Tasks)
	// start the server on port 8000
	log.Fatal(http.ListenAndServe(":8000", nil))
}
