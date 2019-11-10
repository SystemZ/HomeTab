package web

import (
	"github.com/gorilla/mux"
	"gitlab.com/systemz/tasktab/model"
	"log"
	"net/http"
	"strconv"
)

type NewNotePage struct {
	AuthOk   bool
	User     model.User
	Projects []model.Project
	Note     model.Note
	Edit     bool
}

func NewNote(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)
	if !authOk {
		return
	}

	// note was created via form
	if r.Method == http.MethodPost && len(r.FormValue("title")) > 0 {
		projectId, err := strconv.Atoi(r.FormValue("project"))
		// something wrong with project ID
		if err != nil {
			http.Redirect(w, r, "/notes", 302)
			return
		}
		model.CreateNote(
			r.FormValue("title"),
			r.FormValue("body"),
			r.FormValue("tags"),
			uint(projectId),
		)
		// all ok, redirect to note list
		http.Redirect(w, r, "/notes", 302)
		return
	}

	var page NewNotePage
	page.User = user
	page.AuthOk = authOk
	// get data from DB
	model.DB.Order("created_at desc").Find(&page.Projects)
	// render HTML
	display.HTML(w, http.StatusOK, "note_edit", page)
}

type NotesListPage struct {
	AuthOk   bool
	User     model.User
	Projects []model.Project
	Notes    []model.Note
}

func Notes(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)
	if !authOk {
		return
	}
	var page NotesListPage
	page.User = user
	page.AuthOk = authOk
	page.Notes = model.NoteLongList()
	// render HTML
	display.HTML(w, http.StatusOK, "notes", page)
}

type NotePage struct {
	AuthOk bool
	User   model.User
	Note   model.Note
}

func Note(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)
	if !authOk {
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	noteId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong note ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var page NotePage
	page.User = user
	page.AuthOk = authOk
	page.Note = model.OneNote(uint(noteId))[0]
	// render HTML
	display.HTML(w, http.StatusOK, "note", page)
}

type NoteEditPage struct {
	AuthOk   bool
	User     model.User
	Projects []model.Project
	Note     model.Note
	Edit     bool
}

func NoteEdit(w http.ResponseWriter, r *http.Request) {
	authOk, user := CheckAuth(w, r)
	if !authOk {
		return
	}

	// check ID in URL
	vars := mux.Vars(r)
	noteId, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Printf("Wrong note ID requested")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var page NoteEditPage
	page.User = user
	page.AuthOk = authOk
	page.Note = model.OneNote(uint(noteId))[0]
	page.Edit = true
	// get data from DB
	model.DB.Order("created_at desc").Find(&page.Projects)

	// note was edited
	if r.Method == http.MethodPost && len(r.FormValue("title")) > 0 {
		projectId, err := strconv.Atoi(r.FormValue("project"))
		// something wrong with project ID
		if err != nil {
			http.Redirect(w, r, "/note/"+strconv.Itoa(noteId), 302)
			return
		}
		page.Note.Id = uint(noteId)
		page.Note.Title = r.FormValue("title")
		page.Note.Body = r.FormValue("body")
		page.Note.Tags = r.FormValue("tags")
		page.Note.ProjectId = uint(projectId)
		//note.Tags = r.FormValue("tags")
		model.DB.Save(&page.Note)
		// all ok, redirect to note list
		http.Redirect(w, r, "/notes", 302)
		return
	}

	// render HTML
	display.HTML(w, http.StatusOK, "note_edit", page)
}
