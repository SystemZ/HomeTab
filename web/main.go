package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gitlab.com/systemz/tasktab/model"
	"log"
	"net/http"
	"os"
)

var (
	display *render.Render
)

func init() {
	display = render.New(render.Options{
		IndentJSON: true,
		Extensions: []string{".html"},
		//DisableHTTPErrorRendering: false,
		Layout: "layout",
	})

}

func StartWebInterface() {
	// create multiple routes
	r := mux.NewRouter()
	r.HandleFunc("/", Index)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/account", Account)
	r.HandleFunc("/refresh", Refresh)

	// start internal http server with logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("HTTP server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", loggedRouter))
}

func CheckAuth(w http.ResponseWriter, r *http.Request) (ok bool, user model.User) {
	// We can obtain the session token from the requests cookies, which come with every request
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie is not set, return an unauthorized status
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			//w.WriteHeader(http.StatusUnauthorized)
			return false, user
		}
		// For any other type of error, return a bad request status
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		//w.WriteHeader(http.StatusBadRequest)
		return false, user
	}
	sessionToken := c.Value
	// We then get the name of the user from our cache, where we set the session token
	res := model.Redis.Get(sessionToken)

	//response, err := cache.Do("GET", sessionToken)
	//if err != nil {
	//	// If there is an error fetching from cache, return an internal server error status
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//if response == nil {
	if len(res.String()) < 1 {
		// If the session token is not present in cache, return an unauthorized error
		//w.WriteHeader(http.StatusUnauthorized)
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return false, user
	}

	userId, err := res.Uint64()
	_, user = model.GetUserById(uint(userId))
	return true, user
}
