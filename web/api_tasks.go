package web

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/service"
	"log"
	"net/http"
	"strconv"
	"time"
)

type TaskApiResponse struct {
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	AssignedTo  int        `json:"assignedTo"`
	RepeatUnit  string     `json:"repeatUnit"`
	RepeatEvery int        `json:"repeatEvery"`
	CreatedAt   *time.Time `json:"createdAt"`
}

func ApiTaskList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, userInfo := CheckApiAuth(w, r)
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
	// prevent null response
	rawResponse = []TaskApiResponse{}
	// get data from DB, prepare other format
	var tasks []model.Task

	var project model.Project
	// special case, all tasks assigned to user
	// from all projects
	if projectId == 0 {
		model.DB.Order("updated_at desc").Where("(snooze_to <= ? OR snooze_to IS NULL) AND done_at IS NULL AND assigned_user_id = ?", time.Now(), userInfo.Id).Find(&tasks)
	} else {
		model.DB.Order("updated_at desc").Where("project_id = ? AND (snooze_to <= ? OR snooze_to IS NULL) AND done_at IS NULL", projectId, time.Now()).Find(&tasks)
		model.DB.Where(&model.Project{Id: uint(projectId)}).First(&project)
		if project.Id < 1 {
			// no such project, cancel rest of work
			return
		}
	}

	for _, task := range tasks {
		rawResponse = append(rawResponse, TaskApiResponse{
			Id:         int(task.Id),
			Title:      task.Subject,
			AssignedTo: int(task.AssignedUserId),
			RepeatUnit: task.RepeatUnit,
			// FIXME for random intervals
			RepeatEvery: int(task.RepeatBest),
			CreatedAt:   task.CreatedAt,
		})
	}

	// prepare JSON
	res, err := json.MarshalIndent(rawResponse, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(res)

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
	Id          int        `json:"id"`
	Title       string     `json:"title"`
	Snooze      *time.Time `json:"snoozeTo"`
	Delete      bool       `json:"delete"`
	Done        bool       `json:"done"`
	AssignTo    int        `json:"assignTo"`
	RepeatUnit  string     `json:"repeatUnit"`
	RepeatEvery int        `json:"repeatEvery"`
}

func ApiTaskEdit(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, userInDb := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	_, err := strconv.Atoi(vars["id"])
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
		// FIXME validate project ID
		model.DB.Where("id = ?", task.Id).Find(&taskInDb)
		// check if task exists
		if taskInDb.Id < 1 {
			continue
		}
		// snooze if date is in the future
		if task.Snooze != nil && task.Snooze.After(time.Now()) {
			taskInDb.SnoozeTo = task.Snooze
			model.DB.Save(&taskInDb)

			go service.SendGenericNotificationToAllDevices("Task snoozed", taskInDb.Subject, []int{int(userInDb.Id)})
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

		// set as done and finish
		// this is separate from edit form
		if task.Done {
			// add event in event stream
			model.TaskDoneEvent(userInDb.Id, int(taskInDb.Id))
			// set as done in main table
			timeNow := time.Now()
			done := sql.NullTime{
				Time:  timeNow,
				Valid: true,
			}
			taskInDb.DoneAt = &done
			taskInDb.DoneAt.Valid = true
			taskInDb.DoneAt.Time = timeNow
			model.DB.Save(&taskInDb)

			go service.SendGenericNotificationToAllDevices("Task done", taskInDb.Subject, []int{int(userInDb.Id)})
			continue
		}

		// soft delete task
		if task.Delete {
			model.DB.Delete(&taskInDb)
		}
		// edit title
		if len(task.Title) > 0 {
			taskInDb.Subject = task.Title
			model.DB.Save(&taskInDb)
		}
		// assign to users
		if task.AssignTo >= 0 {
			taskInDb.AssignedUserId = uint(task.AssignTo)
			model.DB.Save(&taskInDb)
		}
		// set repeat
		if len(task.RepeatUnit) == 1 && task.RepeatEvery > 0 {
			taskInDb.Repeating = 1
			// FIXME validate char
			taskInDb.RepeatUnit = task.RepeatUnit
			// FIXME use better uint/int
			taskInDb.RepeatMin = uint(task.RepeatEvery)
			taskInDb.RepeatBest = uint(task.RepeatEvery)
			taskInDb.RepeatMax = uint(task.RepeatEvery)
			model.DB.Save(&taskInDb)
		} else {
			taskInDb.Repeating = 0
			taskInDb.RepeatUnit = ""
			taskInDb.RepeatMin = 0
			taskInDb.RepeatBest = 0
			taskInDb.RepeatMax = 0
			model.DB.Save(&taskInDb)
		}

	}

}
