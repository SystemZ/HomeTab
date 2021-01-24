package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
	"time"
)

func Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// (BEGIN) The code uptil this point is the same as the first part of the `Welcome` route
	c, err := r.Cookie("session_token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	sessionToken := c.Value

	redisRes := model.Redis.Get(sessionToken)
	//response, err := cache.Do("GET", sessionToken)
	if redisRes.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if len(redisRes.String()) < 1 {
		//if response == nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// (END) The code uptil this point is the same as the first part of the `Welcome` route

	// Delete the session token
	statusDel := model.Redis.Del(sessionToken)
	if statusDel.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   "",
		Expires: time.Date(1, 0, 0, 0, 0, 0, 0, time.UTC),
	})

	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
