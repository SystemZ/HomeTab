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
