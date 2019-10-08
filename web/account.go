package web

import (
	"gitlab.com/systemz/tasktab/model"
	"log"
	"net/http"
	"strconv"
)

type AccountPage struct {
	AuthOk   bool
	User     model.User
	Projects []model.Project
	Devices  []model.Device
}

func Account(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)

	//FIXME validation
	//FIXME possible race condition

	// project was changed via form
	if r.Method == http.MethodPost && len(r.FormValue("defaultProject")) > 0 {

		projectIdInt, err := strconv.Atoi(r.FormValue("defaultProject"))
		if err != nil {
			log.Printf("%v", err.Error())
			return
		}
		user.DefaultProjectId = uint(projectIdInt)
		model.DB.Save(&user)
		//model.DB.Model(&model.User{}).Where("id = ?", user.Id).Update("defaultProject", r.FormValue("defaultProject"))
		http.Redirect(w, r, "/account", 302)
		return
	}
	// project was created via form
	if r.Method == http.MethodPost && len(r.FormValue("newProjectName")) > 0 {
		model.CreateProject(r.FormValue("newProjectName"), 0) //FIXME
		//newProject := model.Project{Name: r.FormValue("newProjectName")}
		//model.DB.Create(newProject)
		http.Redirect(w, r, "/account", 302)
		return
	}
	// device was created via form
	if r.Method == http.MethodPost && len(r.FormValue("newDeviceName")) > 0 {
		model.CreateDevice(r.FormValue("newDeviceName"), user.Id)
		//newProject := model.Project{Name: r.FormValue("newProjectName")}
		//model.DB.Create(newProject)
		http.Redirect(w, r, "/account", 302)
		return
	}

	// var for all data from DB
	var templateVars AccountPage
	templateVars.User = user
	templateVars.AuthOk = authOk

	// get data from DB
	var projects []model.Project
	model.DB.Order("created_at desc").Find(&projects)
	templateVars.Projects = projects
	model.DB.Order("created_at desc").Find(&templateVars.Devices)

	// render HTML
	display.HTML(w, http.StatusOK, "account", templateVars)
}
