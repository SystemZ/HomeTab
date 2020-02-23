package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/tasktab/model"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TaskApiResponse struct {
	Id        int        `json:"id"`
	Title     string     `json:"title"`
	CreatedAt *time.Time `json:"createdAt"`
}

func ApiTaskList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	projectId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong project ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var rawResponse []TaskApiResponse
	// get data from DB, prepare other format
	var tasks []model.Task
	model.DB.Order("updated_at desc").Where("project_id = ? AND (snooze_to <= ? OR snooze_to IS NULL)", projectId, time.Now()).Find(&tasks)
	//model.DB.Order("updated_at desc").Where(&model.Task{ProjectId: user.DefaultProjectId}).Find(&tasks)
	var project model.Project
	model.DB.Where(&model.Project{Id: uint(projectId)}).First(&project)

	for _, task := range tasks {
		rawResponse = append(rawResponse, TaskApiResponse{
			Id:        int(task.Id),
			Title:     task.Subject,
			CreatedAt: task.CreatedAt,
		})
	}

	// prepare JSON
	noteList, err := json.MarshalIndent(rawResponse, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(noteList)

}
