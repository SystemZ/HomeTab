package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
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
		Funcs: []template.FuncMap{
			{
				"formatDate": func(date time.Time) string {
					loc, err := time.LoadLocation("Europe/Warsaw")
					if err != nil {
						log.Printf("%v", err)
					}
					return date.In(loc).Format("15:04:05 02.01.2006")
				},
			},
		},
	})

}

func StartWebInterface() {
	// create multiple routes
	r := mux.NewRouter()
	// main course
	r.HandleFunc("/", Index)
	r.HandleFunc("/count", Count)
	r.HandleFunc("/count/log", CountLog)
	r.HandleFunc("/device", Device)
	// settings
	r.HandleFunc("/account", Account)
	// auth
	if config.REGISTER_ON {
		r.HandleFunc("/register", Register)
	}
	r.HandleFunc("/login", Login)
	r.HandleFunc("/logout", Logout)
	r.HandleFunc("/refresh", Refresh) // FIXME
	// API
	r.HandleFunc("/api/v1/event", ApiEvent)
	r.HandleFunc("/api/v1/counter", ApiCounterList)
	r.HandleFunc("/api/v1/counter/{id}", ApiCounter)
	// start internal http server with logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("HTTP server started on :3000")
	log.Fatal(http.ListenAndServe(":3000", loggedRouter))
}

func ApiCheckAuth(w http.ResponseWriter, r *http.Request) (ok bool, userId uint) {
	token := r.Header.Get("Authorization")
	//get device from DB by token
	var device model.Device
	model.DB.Where("token = ?", token).First(&device)
	// check auth
	if device.UserId < 1 {
		log.Printf("Unknown device tried access to API")
		w.WriteHeader(http.StatusBadRequest)
		return false, 0
	}
	return true, device.UserId
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

	_, err = res.Result()
	if res.Err() != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

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
