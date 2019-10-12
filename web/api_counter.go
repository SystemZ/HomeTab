package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/tasktab/model"
	"log"
	"net/http"
	"strconv"
)

func ApiCounter(w http.ResponseWriter, r *http.Request) {
	authOk, userId := ApiCheckAuth(w, r)
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
	}

	// gather data, convert from DB model to API model
	var counter CounterApi
	dbCounters := model.CounterLogList(userId)
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
