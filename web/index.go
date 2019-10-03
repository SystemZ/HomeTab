package web

import (
	"gitlab.com/systemz/tasktab/model"
	"math/rand"
	"net/http"
	"time"
)

type TasksPage struct {
	AuthOk  bool
	Tasks   []model.Task
	User    model.User
	Inspire string
}

func Index(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)

	// if new task was added via form
	if r.Method == http.MethodPost {
		task := model.Task{
			Subject:          r.FormValue("newTask"),
			ProjectId:        0,
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
		http.Redirect(w, r, "/", 302)
		return
	}

	var tasks []model.Task
	model.DB.Order("updated_at desc").Limit(10).Find(&tasks)

	var templateVars TasksPage
	templateVars.Tasks = tasks
	templateVars.User = user
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
