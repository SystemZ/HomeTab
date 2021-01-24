package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"io/ioutil"
	"log"
	"net/http"
)

// Android
type IncomingEvent struct {
	Status struct {
		Screen  string `json:"scr"`
		Battery uint   `json:"bat"`
	} `json:"stat"`
	Music struct {
		Track  string `json:"trk"`
		Artist string `json:"art"`
	} `json:"music"`
	Version uint   `json:"v"`
	Token   string `json:"tok"`
}

// Android
func ApiEvent(w http.ResponseWriter, r *http.Request) {
	// FIXME use DeviceApiCheckAuth()
	//enforce POST only
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//get raw JSON
	EventRaw, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("%s", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//parse from JSON
	var event IncomingEvent
	err = json.Unmarshal(EventRaw, &event)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//enforce currently used format for data submit
	if event.Version != 1 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//at this point we have proper JSON

	//get device from DB by token
	var device model.Device
	model.DB.Where("token = ?", event.Token).First(&device)

	if device.UserId < 1 {
		log.Printf("Unknown device tried to submit data")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// device authorized, save presented data

	// battery %
	model.DeviceEventAddInt(model.DeviceBatteryPercent, device.UserId, device.Id, int(event.Status.Battery))
	// screen state
	if event.Status.Screen == "on" {
		model.DeviceEventAdd(model.DeviceScreenOn, device.UserId, device.Id)
	}
	if event.Status.Screen == "off" {
		model.DeviceEventAdd(model.DeviceScreenOff, device.UserId, device.Id)
	}
	// music info
	if event.Music.Track != "%mt_track" {
		model.DeviceEventAddStr(model.DeviceMusicTrack, device.UserId, device.Id, event.Music.Track)
	}
	if event.Music.Artist != "%mt_artist" {
		model.DeviceEventAddStr(model.DeviceMusicArtist, device.UserId, device.Id, event.Music.Artist)
	}
}

func ApiEventList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// get info from DB
	events := model.TaskDoneEvents7days()

	// prepare JSON
	res, err := json.MarshalIndent(events, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(res)

}
