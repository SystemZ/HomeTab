package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type DevicePage struct {
	AuthOk  bool
	User    model.User
	Devices []model.Device
}

func Device(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)

	//FIXME validation
	//FIXME possible race condition

	// var for all data from DB
	var templateVars AccountPage
	templateVars.User = user
	templateVars.AuthOk = authOk

	// get data from DB
	model.DB.Order("created_at desc").Find(&templateVars.Devices)

	// render HTML
	display.HTML(w, http.StatusOK, "device", templateVars)
}
