package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
	"strconv"
)

type CountPage struct {
	AuthOk         bool
	User           model.User
	Counters       []model.CounterList
	CounterRunning bool
}

func Count(w http.ResponseWriter, r *http.Request) {
	var page CountPage
	authOk, user := CheckAuth(w, r)

	//FIXME validation
	//FIXME possible race condition

	// counter was created via form
	if r.Method == http.MethodPost && len(r.FormValue("newCounter")) > 0 {
		model.CreateCounter(r.FormValue("newCounter"))
		http.Redirect(w, r, "/count", 302)
		return
	}
	// counter start via form
	if r.Method == http.MethodPost && len(r.FormValue("startCounter")) > 0 {
		counterId, err := strconv.Atoi(r.FormValue("startCounter"))
		// something wrong with counter ID
		if err != nil {
			http.Redirect(w, r, "/count", 302)
			return
		}
		model.StartCounterSession(uint(counterId), user.Id)
		http.Redirect(w, r, "/count", 302)
		return
	}
	// counter stop via form
	if r.Method == http.MethodPost && len(r.FormValue("stopCounter")) > 0 {
		counterId, err := strconv.Atoi(r.FormValue("stopCounter"))
		// something wrong with counter ID
		if err != nil {
			http.Redirect(w, r, "/count", 302)
			return
		}
		model.StopCounterSession(uint(counterId), user.Id)
		http.Redirect(w, r, "/count", 302)
		return
	}
	// get data from DB
	page.Counters = model.CountersLongList(user.Id)
	for _, counter := range page.Counters {
		if counter.Running == 1 {
			page.CounterRunning = true
			break
		}
	}
	//display HTML
	page.User = user
	page.AuthOk = authOk
	display.HTML(w, http.StatusOK, "count", page)
}
