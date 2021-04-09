package server

import (
	"encoding/json"
	"net/http"

	"github.com/systemz/hometab/internal/model"
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

type PushRegisterRequest struct {
	PushToken string `json:"pushToken"`
}

func ApiPushRegister(w http.ResponseWriter, r *http.Request) {
	authOk, device := DeviceApiCheckAuth(w, r)
	if !authOk {
		w.Write([]byte{})
		return
	}

	decoder := json.NewDecoder(r.Body)
	var newRegistration PushRegisterRequest
	decoder.Decode(&newRegistration)

	// reject if title empty
	if len(newRegistration.PushToken) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// update push token if necessary
	var deviceInDb model.Device
	if device.TokenPush != newRegistration.PushToken {
		model.DB.Model(&deviceInDb).Where("id = ?", device.Id).UpdateColumn("token_push", newRegistration.PushToken)
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
}

type ApiNewDeviceRequest struct {
	Name string `json:"name"`
}

func ApiNewDevice(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, userInfo := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var newDevice ApiNewDeviceRequest
	decoder.Decode(&newDevice)

	// reject if device name is too short
	if len(newDevice.Name) < 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// create device
	model.CreateDevice(newDevice.Name, userInfo.Id)

	// all ok, return device
	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
