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
	model.DB.Order("updated_at desc").Where("project_id = ? AND (snooze_to <= ? OR snooze_to IS NULL) AND done_at IS NULL", projectId, time.Now()).Find(&tasks)
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

type NewTaskApiReq struct {
	Title string `json:"title"`
}

func ApiTaskCreate(w http.ResponseWriter, r *http.Request) {
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

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var newTask NewTaskApiReq
	decoder.Decode(&newTask)

	// add to DB
	task := model.Task{
		Subject:          newTask.Title,
		ProjectId:        uint(projectId),
		AssignedUserId:   0,
		Repeating:        0,
		NeverEnding:      0,
		RepeatUnit:       "",
		RepeatMin:        0,
		RepeatBest:       0,
		RepeatMax:        0,
		EstimateS:        0,
		MasterTaskId:     0,
		SeparateChildren: 0,
	}
	model.CreateTask(task)

}

type EditTaskApiReq struct {
	Id     int        `json:"id"`
	Title  string     `json:"title"`
	Snooze *time.Time `json:"snoozeTo"`
	Delete bool       `json:"delete"`
	Done   bool       `json:"done"`
}

func ApiTaskEdit(w http.ResponseWriter, r *http.Request) {
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

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var editedTasks []EditTaskApiReq
	decoder.Decode(&editedTasks)

	for _, task := range editedTasks {
		var taskInDb model.Task
		model.DB.Where("project_id = ? AND id = ?", projectId, task.Id).Find(&taskInDb)
		// check if task exists
		if taskInDb.Id < 1 {
			continue
		}
		// snooze if date is in the future
		if task.Snooze != nil && task.Snooze.After(time.Now()) {
			taskInDb.SnoozeTo = task.Snooze
			model.DB.Save(&taskInDb)
		}

		// alternative snooze by seconds
		/*
			if taskAction == "snooze" {
				now := time.Now()
				taskSnoozeSecondsRaw := r.FormValue("taskSnoozeSeconds")
				taskSnoozeSeconds, err := strconv.Atoi(taskSnoozeSecondsRaw)
				if err != nil {
					return
				}
				snoozeTime := now.Add(time.Second * time.Duration(taskSnoozeSeconds))
				taskIdInt, err := strconv.Atoi(taskId)
				if err != nil {
					// skip this task if something is wrong
					continue
				}
			}
		*/

		// soft delete task
		if task.Delete {
			model.DB.Delete(&taskInDb)
		}
		// edit title
		if len(task.Title) > 0 {
			taskInDb.Subject = task.Title
			model.DB.Save(&taskInDb)
		}
		// set as done
		if task.Done {
			timeNow := time.Now()
			taskInDb.DoneAt = &timeNow
			model.DB.Save(&taskInDb)
		}
	}

}
