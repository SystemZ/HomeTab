package web

import (
	"gitlab.com/systemz/tasktab/model"
	"net/http"
)

func Tasks(w http.ResponseWriter, r *http.Request) {
	var tasks []model.Task
	model.DB.Order("updated_at desc").Limit(10).Find(&tasks)
	//log.Printf("%v", tasks)
	type TasksPage struct {
		Tasks []model.Task
	}

	var templateVars TasksPage
	templateVars.Tasks = tasks

	display.HTML(w, http.StatusOK, "tasks", templateVars)
	//tmpl := template.Must(template.ParseFiles("templates/tasks.html"))
	//tmpl.Execute(w, templateVars)
}
