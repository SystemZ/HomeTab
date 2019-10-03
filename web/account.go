package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type AccountPage struct {
	AuthOk bool
	User   model.User
}

func Account(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)

	var templateVars TasksPage
	templateVars.User = user
	templateVars.AuthOk = authOk

	display.HTML(w, http.StatusOK, "account", templateVars)
}
