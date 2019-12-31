package web

import (
	"encoding/json"
	"gitlab.com/systemz/tasktab/model"
	"net/http"
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
