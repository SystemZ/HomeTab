package server

import (
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/systemz/hometab/internal/config"
	"github.com/systemz/hometab/internal/model"
	"net/http"
	"time"
)

func Refresh(w http.ResponseWriter, r *http.Request) {
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

	// Now, create a new session token for the current user
	newSessionToken := uuid.NewV4().String()
	//_, err = cache.Do("SETEX", newSessionToken, "120", fmt.Sprintf("%s",response))
	sessionDuration := time.Second * time.Duration(config.SESSION_VALID_S)
	status := model.Redis.Set(newSessionToken, fmt.Sprintf("%s", redisRes.String()), sessionDuration)
	//if err != nil {
	if status.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Delete the older session token
	statusDel := model.Redis.Del(sessionToken)
	//_, err = cache.Do("DEL", sessionToken)
	if statusDel.Err() != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Set the new token as the users `session_token` cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   newSessionToken,
		Expires: time.Now().Add(sessionDuration),
	})
}
