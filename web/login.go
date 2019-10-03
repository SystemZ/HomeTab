package web

import (
	"github.com/satori/go.uuid"
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
	"time"
)

// https://www.sohamkamani.com/blog/2018/03/25/golang-session-authentication/
// Create a struct that models the structure of a user, both in the request body, and in the DB
type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type LoginPage struct {
	AuthFailed bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	page := LoginPage{}
	if r.Method != http.MethodPost {
		display.HTML(w, http.StatusOK, "login", nil)
		return
	}

	var creds Credentials
	//// Get the JSON body and decode into credentials
	//err := json.NewDecoder(r.Body).Decode(&creds)
	//if err != nil {
	//	// If the structure of the body is wrong, return an HTTP error
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	creds.Username = r.FormValue("username")
	creds.Password = r.FormValue("password")

	// Get the expected password from our in memory map
	expectedPassword, ok := users[creds.Username]

	// If a password exists for the given user
	// AND, if it is the same as the password we received, the we can move ahead
	// if NOT, then we return an "Unauthorized" status
	if !ok || expectedPassword != creds.Password {
		//w.WriteHeader(http.StatusUnauthorized)
		page.AuthFailed = true
		display.HTML(w, http.StatusUnauthorized, "login", page)
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	sessionDuration := time.Second * time.Duration(config.SESSION_VALID_S)
	sessionCreation := model.Redis.Set(sessionToken, creds.Username, sessionDuration)
	//_, err = cache.Do("SETEX", sessionToken, "120", creds.Username)
	if sessionCreation.Err() != nil {
		// If there is an error in setting the cache, return an internal server error
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	// we also set an expiry time of 120 seconds, the same as the cache
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(sessionDuration),
	})

	http.Redirect(w, r, "/tasks", http.StatusTemporaryRedirect)
}
