package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

type CounterApi struct {
	Id         uint     `json:"id"`
	Name       string   `json:"name"`
	Tags       []string `json:"tags"`
	Seconds    uint     `json:"seconds"`
	InProgress bool     `json:"inProgress"`
}

func ApiCounterList(w http.ResponseWriter, r *http.Request) {
	authDeviceOk, deviceInfo := DeviceApiCheckAuth(w, r)
	authUserOk, userInfo := CheckApiAuth(w, r)
	// deny access if neither auth method works
	if !authUserOk && !authDeviceOk {
		w.Write([]byte{})
		return
	}

	var userId uint
	if authDeviceOk {
		userId = deviceInfo.UserId
	}
	if authUserOk {
		userId = userInfo.Id
	}

	// gather data, convert from DB model to API model
	var counters []CounterApi
	dbCounters := model.CountersLongList(userId)
	for _, counter := range dbCounters {
		counters = append(counters, CounterApi{
			Id:         counter.Id,
			Name:       counter.Name,
			Tags:       []string{counter.Tags},
			Seconds:    counter.SecondsAll,
			InProgress: counter.Running == 1,
		})
	}

	// prepare JSON
	counterList, err := json.MarshalIndent(counters, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(counterList)

}
