package web

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"log"
	"net/http"
	"time"
)

type ApiLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApiLoginResponse struct {
	Token string `json:"token"`
}

func ApiLogin(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginCreds ApiLoginRequest
	decoder.Decode(&loginCreds)

	log.Printf("used credentials: %+v", loginCreds)

	//TODO check length of credentials first
	passOk, user := model.IsPasswordOk(loginCreds.Username, loginCreds.Password)

	if !passOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Create a new random session token
	sessionToken := uuid.NewV4().String()
	// Set the token in the cache, along with the user whom it represents
	sessionDuration := time.Second * time.Duration(config.SESSION_VALID_S)
	sessionCreation := model.Redis.Set(sessionToken, user.Id, sessionDuration)
	if sessionCreation.Err() != nil {
		// something wrong with saving token to DB
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// prepare JSON
	rawResponse := ApiLoginResponse{Token: sessionToken}
	jsonResponse, err := json.Marshal(&rawResponse)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(jsonResponse)
}
