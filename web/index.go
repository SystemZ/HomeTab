package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type TasksPage struct {
	AuthOk     bool
	Tasks      []model.Task
	TasksCount uint
	User       model.User
	Project    model.Project
	Inspire    string
}

func Index(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)
	if !authOk {
		return
	}

	var templateVars TasksPage
	templateVars.User = user
	templateVars.AuthOk = authOk

	display.HTML(w, http.StatusOK, "index", templateVars)
}
