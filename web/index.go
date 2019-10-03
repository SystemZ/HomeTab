package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type TasksPage struct {
	AuthOk bool
	Tasks  []model.Task
	User   model.User
}

func Index(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)

	var tasks []model.Task
	model.DB.Order("updated_at desc").Limit(10).Find(&tasks)

	var templateVars TasksPage
	templateVars.Tasks = tasks
	templateVars.User = user
	templateVars.AuthOk = authOk

	display.HTML(w, http.StatusOK, "index", templateVars)
}
