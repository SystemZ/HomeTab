package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type CountPage struct {
	AuthOk   bool
	User     model.User
	Counters []model.Counter
}

func Count(w http.ResponseWriter, r *http.Request) {
	var page CountPage
	authOk, user := CheckAuth(w, r)

	//FIXME validation
	//FIXME possible race condition
	// project was created via form
	if r.Method == http.MethodPost && len(r.FormValue("newCounter")) > 0 {
		model.CreateCounter(r.FormValue("newCounter"))
		http.Redirect(w, r, "/count", 302)
		return
	}
	// get data from DB
	model.DB.Order("created_at desc").Find(&page.Counters)
	//display HTML
	page.User = user
	page.AuthOk = authOk
	display.HTML(w, http.StatusOK, "count", page)
}
