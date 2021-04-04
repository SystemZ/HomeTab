package server

import (
	"encoding/json"
	"net/http"

	"github.com/systemz/hometab/internal/model"
)

type ApiNewUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func ApiNewUser(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var newUser ApiNewUserRequest
	decoder.Decode(&newUser)

	// reject if username is too short
	if len(newUser.Username) < 3 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create user
	model.CreateUser(newUser.Username, newUser.Email, newUser.Password)

	// all ok, return user
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
