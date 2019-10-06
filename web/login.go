package web

import (
	"github.com/satori/go.uuid"
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
	"time"
)

// https://www.sohamkamani.com/blog/2018/03/25/golang-session-authentication/
type LoginPage struct {
	AuthFailed bool
	RegisterOn bool
	AuthOk     bool
}

func Login(w http.ResponseWriter, r *http.Request) {
	page := LoginPage{}
	page.AuthOk = false
	if config.REGISTER_ON {
		page.RegisterOn = true
	}
	if r.Method != http.MethodPost {
		display.HTML(w, http.StatusOK, "login", nil)
		return
	}

	//TODO check length of credentials first
	passOk, user := model.IsPasswordOk(r.FormValue("username"), r.FormValue("password"))
	if !passOk {
		page.AuthFailed = true
		display.HTML(w, http.StatusUnauthorized, "login", page)
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	// The token has an expiry time of 120 seconds
	sessionDuration := time.Second * time.Duration(config.SESSION_VALID_S)
	sessionCreation := model.Redis.Set(sessionToken, user.Id, sessionDuration)
	if sessionCreation.Err() != nil {
		// If there is an error in setting the cache, return an internal server error
		page.AuthFailed = true
		display.HTML(w, http.StatusInternalServerError, "login", page)
		return
	}

	// Finally, we set the client cookie for "session_token" as the session token we just generated
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(sessionDuration),
	})
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
