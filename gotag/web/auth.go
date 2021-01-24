package web

import (
	"encoding/json"
	uuid "github.com/satori/go.uuid"
	"gitlab.com/systemz/gotag/config"
	"gitlab.com/systemz/gotag/model"
	"net/http"
	"strings"
	"time"
)

type ApiLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type ApiLoginResponse struct {
	Token string `json:"token"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var loginCreds ApiLoginRequest
	decoder.Decode(&loginCreds)

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
