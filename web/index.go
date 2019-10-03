package web

import (
	"gitlab.com/systemz/tasktab/model"
	"math/rand"
	"net/http"
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

	// if new task was added via form
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

	// always redirect from POST to prevent F5 problems
	if r.Method == http.MethodPost {
		http.Redirect(w, r, "/", 301)
	}

	// get stuff from DB
	var tasks []model.Task
	model.DB.Order("updated_at desc").Where(&model.Task{ProjectId: user.DefaultProjectId}).Find(&tasks)
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
