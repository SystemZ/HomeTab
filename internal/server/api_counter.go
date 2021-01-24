package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/systemz/hometab/internal/model"
	"github.com/systemz/hometab/internal/service"
	"log"
	"net/http"
	"strconv"
)

func ApiCounter(w http.ResponseWriter, r *http.Request) {
	authOk, device := DeviceApiCheckAuth(w, r)
	if !authOk {
		w.Write([]byte{})
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	counterId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong counter ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// gather data, convert from DB model to API model
	var counter CounterApi
	dbCounters := model.CounterLogList(device.UserId, 100)
	for _, v := range dbCounters {
		if v.CounterId == uint(counterId) {
			counter = CounterApi{
				Id:         v.CounterId,
				Name:       v.Name,
				Tags:       []string{v.Tags},
				Seconds:    v.Duration,
				InProgress: v.Running,
			}
			break
		}
	}

	// prepare JSON
	counterList, err := json.MarshalIndent(counter, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(counterList)
}

func ApiCounterFrontend(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	counterId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong counter ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// gather data, convert from DB model to API model
	dbCounters := model.CounterLog(counterId, userInfo.Id, 100)

	// set info about counter
	var counter CounterApi
	if len(dbCounters) > 0 {
		counter.Id = dbCounters[0].Id
		counter.Name = dbCounters[0].Name
		counter.Tags = []string{dbCounters[0].Tags}
		counter.Seconds = dbCounters[0].Duration
		counter.InProgress = dbCounters[0].Running
	}

	// set counter sessions
	for _, v := range dbCounters {
		counter.Sessions = append(counter.Sessions, CounterLogApi{
			Duration:          v.Duration,
			DurationFormatted: v.DurationFormatted,
			Start:             v.StartedAt,
			End:               v.EndedAt.Time,
		})
	}

	// set counter stats
	counter.Stats = model.CounterStats(uint(counterId), userInfo.Id)[0]

	// prepare JSON
	counterList, err := json.MarshalIndent(counter, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(counterList)
}

func ApiCounterStart(w http.ResponseWriter, r *http.Request) {
	authOk, device := DeviceApiCheckAuth(w, r)
	if !authOk {
		w.Write([]byte{})
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	counterId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong counter ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//FIXME validation for user permissions
	sessionId := model.StartCounterSession(uint(counterId), device.UserId)

	// notify mobile app
	var user model.User
	model.DB.Where(model.User{Id: device.UserId}).First(&user)
	go service.SendCounterNotification(true, user, uint(counterId), sessionId, "")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func ApiCounterStop(w http.ResponseWriter, r *http.Request) {
	authOk, device := DeviceApiCheckAuth(w, r)
	if !authOk {
		w.Write([]byte{})
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	counterId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong counter ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//FIXME validation for user permissions
	sessionId, sessionTaken := model.StopCounterSession(uint(counterId), device.UserId)

	// notify mobile app
	var user model.User
	model.DB.Where(model.User{Id: device.UserId}).First(&user)
	go service.SendCounterNotification(false, user, uint(counterId), sessionId, sessionTaken)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func ApiCounterStartFrontend(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	counterId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong counter ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//FIXME validation for user permissions
	sessionId := model.StartCounterSession(uint(counterId), userInfo.Id)

	// notify mobile app
	var user model.User
	model.DB.Where(model.User{Id: userInfo.Id}).First(&user)
	go service.SendCounterNotification(true, user, uint(counterId), sessionId, "")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

func ApiCounterStopFrontend(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	counterId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong counter ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	//FIXME validation for user permissions
	sessionId, sessionTaken := model.StopCounterSession(uint(counterId), userInfo.Id)
	// notify mobile app
	var user model.User
	model.DB.Where(model.User{Id: userInfo.Id}).First(&user)
	go service.SendCounterNotification(false, user, uint(counterId), sessionId, sessionTaken)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}

type NewCounterRequest struct {
	Tag  string `json:"tag"`
	Name string `json:"name"`
}

func ApiCounterAdd(w http.ResponseWriter, r *http.Request) {
	authUserOk, _ := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var newCounter NewCounterRequest
	decoder.Decode(&newCounter)

	if len(newCounter.Tag) < 1 || len(newCounter.Name) < 1 {
		log.Printf("Too short counter")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	counterId := model.CreateCounter(newCounter.Name)
	newTag := model.CounterTag{
		CounterId: counterId,
		Name:      newCounter.Tag,
	}
	model.DB.Save(&newTag)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte{})
}
