package server

import (
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/systemz/hometab/internal/config"
	"github.com/systemz/hometab/internal/model"
	"github.com/systemz/hometab/internal/service/feed/generator/helios"
	"github.com/systemz/hometab/internal/service/feed/generator/kernelcare"
	"github.com/systemz/hometab/internal/service/feed/generator/tm_gdynia"
	"github.com/systemz/hometab/internal/service/feed/generator/tm_poznan"
	"github.com/systemz/hometab/internal/service/feed/generator/tw_szczecin"
	"github.com/unrolled/render"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

var (
	display *render.Render
)

func init() {
	display = render.New(render.Options{
		Directory:  config.TEMPLATE_PATH,
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
	r.HandleFunc("/api/v1/user", ApiNewUser).Methods("POST")
	r.HandleFunc("/api/v1/project", ApiProjectList).Methods("GET")
	r.HandleFunc("/api/v1/project/{id}/task", ApiTaskList).Methods("GET")
	r.HandleFunc("/api/v1/project/{id}/task", ApiTaskCreate).Methods("POST")
	r.HandleFunc("/api/v1/project/{id}/task", ApiTaskEdit).Methods("PUT")
	r.HandleFunc("/api/v1/note", ApiNoteList).Methods("GET")
	r.HandleFunc("/api/v1/note", ApiNoteNew).Methods("POST")
	r.HandleFunc("/api/v1/note/{id}", ApiNote).Methods("GET")
	r.HandleFunc("/api/v1/note/{id}", ApiNoteEdit).Methods("PUT")
	r.HandleFunc("/api/v1/note/{id}", ApiNoteDelete).Methods("DELETE")
	r.HandleFunc("/api/v1/counter", ApiCounterAdd).Methods("POST")
	r.HandleFunc("/api/v1/counter-page", ApiCounterListPagination).Methods("POST")
	r.HandleFunc("/api/v1/counter/{id}/info", ApiCounterFrontend).Methods("GET")
	r.HandleFunc("/api/v1/counter/{id}/start", ApiCounterStartFrontend).Methods("PUT")
	r.HandleFunc("/api/v1/counter/{id}/stop", ApiCounterStopFrontend).Methods("PUT")
	r.HandleFunc("/api/v1/event", ApiEventList).Methods("GET")
	r.HandleFunc("/api/v1/device", ApiDeviceList).Methods("GET")

	// for Android
	r.HandleFunc("/api/v1/push/register", ApiPushRegister)
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

	//
	//
	// gotag part
	//
	//

	// JSON API for JS frontend
	r.HandleFunc("/api/v1/scan", Scan).Methods("POST")
	r.HandleFunc("/api/v1/files", FilePaginate).Methods("POST")
	r.HandleFunc("/api/v1/file/{id}", OneFile).Methods("GET")
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
	log.Println("HTTP server started on :" + config.HTTP_PORT)

	if config.DEV_MODE {
		log.Println("DEV env detected, CORS wildcard")
		loggedRouter = cors.AllowAll().Handler(loggedRouter)
	}

	log.Fatal(http.ListenAndServe(":"+config.HTTP_PORT, loggedRouter))
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
