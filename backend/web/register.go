package web

import (
	"gitlab.com/systemz/tasktab/config"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

// https://www.sohamkamani.com/blog/2018/03/25/golang-session-authentication/
type RegisterPage struct {
	AuthOk              bool
	RegisterFailed      bool
	RegisterTokenNeeded bool
}

func Register(w http.ResponseWriter, r *http.Request) {
	page := RegisterPage{}
	if config.REGISTER_WHITELIST {
		page.RegisterTokenNeeded = true
	}

	if r.Method == http.MethodGet {
		display.HTML(w, http.StatusOK, "register", page)
		return
	}
	if r.Method == http.MethodPost {
		//FIXME validation
		if config.REGISTER_WHITELIST && r.FormValue("token") != config.REGISTER_TOKEN {
			page.RegisterFailed = true
			display.HTML(w, http.StatusOK, "register", page)
			return
		}
		model.CreateUser(r.FormValue("username"), r.FormValue("email"), r.FormValue("password"))
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	// if something goes wrong
	http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
}
