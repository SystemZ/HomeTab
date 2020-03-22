package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

func ApiDeviceList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get info from DB
	deviceList := model.GetListOfDevices()

	// prepare JSON
	res, err := json.MarshalIndent(deviceList, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
