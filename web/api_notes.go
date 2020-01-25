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

type NoteApiResponse struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Short     string     `json:"short"`
	Tags      string     `json:"tags"`
	CreatedAt *time.Time `json:"createdAt"`
}

func ApiNoteList(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	var rawResponse []NoteApiResponse
	// get data from DB, prepare other format
	notes := model.NoteLongList()
	for _, note := range notes {
		rawResponse = append(rawResponse, NoteApiResponse{
			Id:        note.Id,
			Title:     note.Title,
			Short:     note.BodyShort,
			Tags:      note.Tags,
			CreatedAt: note.CreatedAt,
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

// FIXME for note edit
type OneNoteApi struct {
	Id        uint       `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Tags      string     `json:"tags"`
	CreatedAt *time.Time `json:"createdAt"`
}

func ApiNote(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
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

	// get note from DB
	notesInDb := model.OneNote(uint(noteId))
	// no such note
	if len(notesInDb) < 1 {
		log.Printf("No note with ID %v found", noteId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var rawResponse OneNoteApi
	rawResponse.Id = notesInDb[0].Id
	rawResponse.Title = notesInDb[0].Title
	rawResponse.Body = notesInDb[0].Body
	rawResponse.Tags = notesInDb[0].Tags
	rawResponse.CreatedAt = notesInDb[0].CreatedAt

	// prepare JSON
	response, err := json.MarshalIndent(rawResponse, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func ApiNoteEdit(w http.ResponseWriter, r *http.Request) {
	// check auth
	ok, _ := CheckApiAuth(w, r)
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
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

	// get note from DB
	notesInDb := model.OneNote(uint(noteId))
	// no such note
	if len(notesInDb) < 1 {
		log.Printf("No note with ID %v found", noteId)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var noteEdited OneNoteApi
	decoder.Decode(&noteEdited)

	oneNote := notesInDb[0]
	oneNote.Body = noteEdited.Body
	model.DB.Save(&oneNote)

	// all ok, return list
	w.WriteHeader(http.StatusOK)
}
