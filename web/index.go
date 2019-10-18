package web

import (
	"gitlab.com/systemz/tasktab/model"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type TasksPage struct {
	AuthOk     bool
	Tasks      []model.Task
	TasksCount uint
	User       model.User
	Project    model.Project
	Inspire    string
}

func Index(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)

	if !authOk {
		return
	}

	// new task was added via form
	if r.Method == http.MethodPost && len(r.FormValue("newTask")) > 0 {
		task := model.Task{
			Subject:          r.FormValue("newTask"),
			ProjectId:        user.DefaultProjectId,
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
	// task action via form
	if r.Method == http.MethodPost && len(r.FormValue("taskAction")) > 0 {
		//r.ParseForm() // Required if you don't call r.FormValue()
		taskAction := r.FormValue("taskAction")
		//FIXME validation
		for _, taskId := range r.Form["taskId"] {
			// TODO add actions to log for creating task timeline
			log.Printf("action: %v, taskId: %v", taskAction, taskId)
			if taskAction == "delete" {
				model.DB.Where("id = ?", taskId).Delete(&model.Task{})
			} else if taskAction == "snooze" {
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
				model.DB.Model(&model.Task{Id: uint(taskIdInt)}).Update("SnoozeTo", &snoozeTime)
			}
		}
	}

	// always redirect from POST to prevent F5 problems
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/", 302)
		return
	}

	// get stuff from DB
	var tasks []model.Task
	model.DB.Order("updated_at desc").Where("project_id = ? AND (snooze_to <= ? OR snooze_to IS NULL)", user.DefaultProjectId, time.Now()).Find(&tasks)
	//model.DB.Order("updated_at desc").Where(&model.Task{ProjectId: user.DefaultProjectId}).Find(&tasks)
	var project model.Project
	model.DB.Where(&model.Project{Id: user.DefaultProjectId}).First(&project)

	var templateVars TasksPage
	templateVars.Tasks = tasks
	templateVars.TasksCount = uint(len(tasks))
	templateVars.User = user
	templateVars.Project = project
	templateVars.AuthOk = authOk

	inspirePool := []string{
		"Take a walk",
		"Take a shower",
		"Buy a yacht",
		"Brush your teeth",
		"Clean your desk",
		"Pet a cat",
	}

	rand.Seed(time.Now().UnixNano())
	luckyNum := rand.Intn(len(inspirePool) - 1)
	templateVars.Inspire = inspirePool[luckyNum]

	display.HTML(w, http.StatusOK, "index", templateVars)
}
