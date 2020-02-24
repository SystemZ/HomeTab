package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type UserListRes struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
}

func ApiUserList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var usersInDb []model.User
	model.DB.Find(&usersInDb)

	var rawResponse []UserListRes
	// prevent null response
	rawResponse = []UserListRes{}
	for _, user := range usersInDb {
		rawResponse = append(rawResponse, UserListRes{
			Id:       int(user.Id),
			Username: user.Username,
		})
	}

	// prepare JSON
	response, err := json.MarshalIndent(rawResponse, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(response)

}
