package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
	"strconv"
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

type CounterApiPagination struct {
	Pagination struct {
		AllRecords int `json:"allRecords"`
	} `json:"pagination"`
	Counters []CounterApi `json:"counters"`
}

func ApiCounterListPagination(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userId := userInfo.Id

	// get limitStr on one page
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	// get nextIdStr on one page
	nextIdStr := r.URL.Query().Get("nextId")
	nextId, err := strconv.Atoi(nextIdStr)
	if err != nil || nextId < 1 {
		nextId = 0
	}
	// get prevIdStr on one page
	prevIdStr := r.URL.Query().Get("prevId")
	prevId, err := strconv.Atoi(prevIdStr)
	if err != nil || prevId < 1 {
		prevId = 0
	}

	// gather data, convert from DB model to API model
	var rawRes CounterApiPagination
	var counters []CounterApi
	dbCounters, allRecords := model.CountersLongListPaginate(userId, limit, nextId, prevId)
	for _, counter := range dbCounters {
		counters = append(counters, CounterApi{
			Id:         counter.Id,
			Name:       counter.Name,
			Tags:       []string{counter.Tags},
			Seconds:    counter.SecondsAll,
			InProgress: counter.Running == 1,
		})
	}
	rawRes.Counters = counters
	rawRes.Pagination.AllRecords = allRecords

	// prepare JSON
	counterList, err := json.MarshalIndent(rawRes, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(counterList)

}
