package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

func ApiProjectList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var projects []model.Project
	model.DB.Find(&projects)

	// prepare JSON
	noteList, err := json.MarshalIndent(projects, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(noteList)

}
