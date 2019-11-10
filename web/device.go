package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type DevicePage struct {
	AuthOk  bool
	User    model.User
	Devices []model.DeviceList
}

func Device(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)
	if !authOk {
		return
	}

	var page DevicePage
	page.User = user
	page.AuthOk = authOk
	page.Devices = model.GetListOfDevices()
	// render HTML
	display.HTML(w, http.StatusOK, "device", page)
}
