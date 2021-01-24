package server

import (
	"encoding/json"
	"github.com/systemz/hometab/internal/model"
	"net/http"
)

func ApiProjectList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get info from DB
	var projects []model.Project
	model.DB.Find(&projects)

	// add special case
	// project with ID 0
	// all one user tasks from all projects
	var finalProjectList []model.Project
	finalProjectList = append(finalProjectList, model.Project{
		Id:   0,
		Name: "My tasks",
	})
	for _, project := range projects {
		finalProjectList = append(finalProjectList, project)
	}

	// prepare JSON
	res, err := json.MarshalIndent(finalProjectList, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
