package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/systemz/hometab/internal/model"
	"net/http"
	"strconv"
)

type PaginateFileResponse struct {
	Pagination struct {
		AllRecords int `json:"allRecords"`
	} `json:"pagination"`
	Files []model.File `json:"files"`
}

func FileSimilar(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	rawRes := model.SimilarFiles(vars["sha256"], int(userInfo.Id), 50)

	// prepare JSON result
	fileList, err := json.MarshalIndent(rawRes, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(fileList)
}

func OneFile(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	vars := mux.Vars(r)
	fileId, err := strconv.Atoi(vars["id"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	rawRes := model.GetFileByIdForUser(fileId, int(userInfo.Id))

	// prepare JSON result
	file, err := json.MarshalIndent(rawRes, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return file
	w.WriteHeader(http.StatusOK)
	w.Write(file)
}

func FilePaginate(w http.ResponseWriter, r *http.Request) {
	authUserOk, userInfo := CheckApiAuth(w, r)
	if !authUserOk {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// get limitStr on one page
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}
	// get nextIdStr on one page
	nextIdStr := r.URL.Query().Get("nextId")
	nextId, err := strconv.Atoi(nextIdStr)
	if err != nil || nextId < 1 {
		nextId = 0
	}
	// get prevIdStr on one page
	prevIdStr := r.URL.Query().Get("prevId")
	prevId, err := strconv.Atoi(prevIdStr)
	if err != nil || prevId < 1 {
		prevId = 0
	}

	// get search term
	decoder := json.NewDecoder(r.Body)
	var fileQuery PaginateQueryRequest
	decoder.Decode(&fileQuery)
	searchTerm := fileQuery.Query

	// gather data, convert from DB model to API model
	var rawRes PaginateFileResponse
	//var counters []CounterApi
	dbFiles, allRecords := model.FileListPaginate(int(userInfo.Id), limit, nextId, prevId, searchTerm)

	rawRes.Pagination.AllRecords = allRecords
	rawRes.Files = dbFiles
	// prevent null result in JSON, make empty array instead
	if allRecords < 1 {
		rawRes.Files = []model.File{}
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// prepare JSON result
	fileList, err := json.MarshalIndent(rawRes, "", "  ")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// all ok, return list
	w.WriteHeader(http.StatusOK)
	w.Write(fileList)
}

//func Scan /api/v1/file/scan
/*
	type FileScanRequestBody struct {
		FilePath string   `json:"filePath"`
		Tags     []string `json:"tags"`
		ParentId int      `json:"parentId"`
	}

	// parse JSON
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	var requestBody FileScanRequestBody
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Printf("Error when handling file scan: %v", err.Error())
		panic(err)
	}

	// make all hard work
	fileInDb := core.AddFile(db, requestBody.FilePath, core.AddFileOptions{
		CalcSimilarity: true,
		GenerateThumbs: true,
		Tags:           requestBody.Tags,
		ParentId:       requestBody.ParentId,
	})

	// send reponse to user
	jsonResponse, err := json.Marshal(fileInDb)
	w.Write(jsonResponse)
*/
