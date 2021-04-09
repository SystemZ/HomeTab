package server

import (
	"encoding/json"
	"net/http"

	"github.com/systemz/hometab/internal/model"
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

type ApiNewProjectRequest struct {
	Name    string `json:"name"`
	GroupId uint   `json:"groupId"`
}

func ApiNewProject(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var newProject ApiNewProjectRequest
	decoder.Decode(&newProject)

	// reject if project name is empty
	if len(newProject.Name) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create project, default group ID is 0
	model.CreateProject(newProject.Name, 0)

	// all ok, return project
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
