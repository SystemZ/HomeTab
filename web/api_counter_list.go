package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"log"
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
	token := r.Header.Get("Authorization")
	//get device from DB by token
	var device model.Device
	model.DB.Where("token = ?", token).First(&device)
	// check auth
	if device.UserId < 1 {
		log.Printf("Unknown device tried access counter list")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// gather data, convert from DB model to API model
	var counters []CounterApi
	dbCounters := model.CountersLongList(device.UserId)
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
