package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/gotag/model2"
	"log"
	"net/http"
)

type TagRequest struct {
	Tag string `json:"tag"`
}

func TagList(w http.ResponseWriter, r *http.Request) {
	tags := model2.TagList()
	// prepare JSON result
	tagList, err := json.MarshalIndent(tags, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(tagList)
}

func TagAdd(w http.ResponseWriter, r *http.Request) {
	// get SHA256 from URL
	vars := mux.Vars(r)

	// get tag name from JSON body
	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var tagAddReq TagRequest
	decoder.Decode(&tagAddReq)

	if len(tagAddReq.Tag) < 1 {
		log.Printf("Tag '%v' too short!", tagAddReq.Tag)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var imgInDb model2.File
	model2.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)
	model2.AddTagToFile(model2.DB, tagAddReq.Tag, imgInDb.Id)
}

func TagDelete(w http.ResponseWriter, r *http.Request) {
	// get SHA256 from URL
	vars := mux.Vars(r)

	// get tag name from JSON body
	// FIXME validate
	decoder := json.NewDecoder(r.Body)
	var tagDelReq TagRequest
	decoder.Decode(&tagDelReq)

	if len(tagDelReq.Tag) < 1 {
		log.Printf("Tag '%v' too short!", tagDelReq.Tag)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var imgInDb model2.File
	model2.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)

	var tagInDb model2.Tag
	model2.DB.Where("tag = ?", tagDelReq.Tag).First(&tagInDb)

	// remove link between tag and file
	model2.DB.Where("file_id = ? AND tag_id = ?", imgInDb.Id, tagInDb.Id).Delete(model2.FileTag{})
}

//TODO method to totally remove tag from all files
