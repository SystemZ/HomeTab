package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/gotag/model"
	"log"
	"net/http"
)

type TagRequest struct {
	Tag string `json:"tag"`
}

func TagList(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	tags := model.TagList(int(userInfo.Id))
	// prevent returning null in JSON
	if len(tags) < 1 {
		tags = []model.Tag{}
	}
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
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	var imgInDb model.File
	model.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)
	model.AddTagToFile(model.DB, tagAddReq.Tag, imgInDb.Id, int(userInfo.Id))
}

func TagDelete(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

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

	var imgInDb model.File
	model.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)

	var tagInDb model.Tag
	model.DB.Where("tag = ?", tagDelReq.Tag).First(&tagInDb)

	// remove link between tag and file
	model.DB.Where("file_id = ? AND tag_id = ? AND user_id = ?", imgInDb.Id, tagInDb.Id, int(userInfo.Id)).Delete(model.FileTag{})
}

//TODO method to totally remove tag from all files
