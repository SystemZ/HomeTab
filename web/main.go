package web

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/unrolled/render"
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/service/feed_backends/helios"
	"gitlab.com/systemz/tasktab/service/feed_backends/kernelcare"
	"gitlab.com/systemz/tasktab/service/feed_backends/tm_gdynia"
	"gitlab.com/systemz/tasktab/service/feed_backends/tm_poznan"
	"gitlab.com/systemz/tasktab/service/feed_backends/tw_szczecin"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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

const (
	STATIC_DIR = "/new/"
)

func StartWebInterface() {
	// create multiple routes
	r := mux.NewRouter()
	// main course
	r.PathPrefix(STATIC_DIR).Handler(http.StripPrefix(STATIC_DIR, http.FileServer(http.Dir("."+STATIC_DIR))))
	r.HandleFunc("/", Index)
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
	r.HandleFunc("/api/v1/login", ApiLogin)

	// for frontend
	r.HandleFunc("/api/v1/user", ApiUserList).Methods("GET")
	r.HandleFunc("/api/v1/project", ApiProjectList).Methods("GET")
	r.HandleFunc("/api/v1/project/{id}/task", ApiTaskList).Methods("GET")
	r.HandleFunc("/api/v1/project/{id}/task", ApiTaskCreate).Methods("POST")
	r.HandleFunc("/api/v1/project/{id}/task", ApiTaskEdit).Methods("PUT")
	r.HandleFunc("/api/v1/note", ApiNoteList).Methods("GET")
	r.HandleFunc("/api/v1/note", ApiNoteNew).Methods("POST")
	r.HandleFunc("/api/v1/note/{id}", ApiNote).Methods("GET")
	r.HandleFunc("/api/v1/note/{id}", ApiNoteEdit).Methods("PUT")
	r.HandleFunc("/api/v1/counter", ApiCounterAdd).Methods("POST")
	r.HandleFunc("/api/v1/counter-page", ApiCounterListPagination).Methods("POST")
	r.HandleFunc("/api/v1/counter/{id}/info", ApiCounterFrontend).Methods("GET")
	r.HandleFunc("/api/v1/counter/{id}/start", ApiCounterStartFrontend).Methods("PUT")
	r.HandleFunc("/api/v1/counter/{id}/stop", ApiCounterStopFrontend).Methods("PUT")
	r.HandleFunc("/api/v1/event", ApiEventList).Methods("GET")
	r.HandleFunc("/api/v1/device", ApiDeviceList).Methods("GET")

	// for Android
	r.HandleFunc("/api/v1/mq/access", ApiMqCredential)
	r.HandleFunc("/api/v1/event", ApiEvent).Methods("POST")
	r.HandleFunc("/api/v1/counter", ApiCounterList).Methods("GET")
	r.HandleFunc("/api/v1/counter/{id}", ApiCounter)
	r.HandleFunc("/api/v1/counter/{id}/start", ApiCounterStart).Methods("POST")
	r.HandleFunc("/api/v1/counter/{id}/stop", ApiCounterStop).Methods("POST")

	// Feeds
	r.HandleFunc("/feed/kernelcare", kernelcare.Serve)
	r.HandleFunc("/feed/helios", helios.Serve)
	r.HandleFunc("/feed/tw-szczecin", tw_szczecin.Serve)
	r.HandleFunc("/feed/tm-gdynia", tm_gdynia.Serve)
	r.HandleFunc("/feed/tm-poznan", tm_poznan.Serve)

	// Webhooks
	r.HandleFunc("/wh/gitlab", WebhookGitlab)
	// start internal http server with logging
	loggedRouter := handlers.LoggingHandler(os.Stdout, r)
	log.Println("HTTP server started on :3000")

	if config.DEV_MODE {
		log.Println("DEV env detected, CORS wildcard")
		loggedRouter = cors.AllowAll().Handler(loggedRouter)
	}

	log.Fatal(http.ListenAndServe(":3000", loggedRouter))
}

func DeviceApiCheckAuth(w http.ResponseWriter, r *http.Request) (ok bool, device model.Device) {
	token := r.Header.Get("Authorization")
	//get device from DB by token
	model.DB.Where("token = ?", token).First(&device)
	// check auth
	if device.UserId < 1 {
		log.Printf("Unknown device tried access to API")
		//w.WriteHeader(http.StatusBadRequest)
		return false, device
	}
	return true, device
}

func CheckApiAuth(w http.ResponseWriter, r *http.Request) (ok bool, user model.User) {
	tokenInHeader := r.Header.Get("Authorization")
	if len(tokenInHeader) != 43 {
		return false, user
	}
	tokenSplit := strings.Split(tokenInHeader, " ")
	if tokenSplit[0] != "Bearer" {
		return false, user
	}
	if len(tokenSplit[1]) != 36 {
		return false, user
	}
	//Bearer 0b97c6a3-2415-4b5e-b144-268fdf6af6da

	res := model.Redis.Get(tokenSplit[1])
	_, err := res.Result()
	if res.Err() != nil {
		return false, user
	}
	if len(res.String()) < 1 {
		// If the session token is not present in cache, return an unauthorized error
		return false, user
	}

	userId, err := res.Uint64()
	if err != nil {
		return false, user
	}

	_, user = model.GetUserById(uint(userId))
	return true, user
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
