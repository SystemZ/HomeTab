package web

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"gitlab.com/systemz/gotag/core"
	"gitlab.com/systemz/gotag/model2"
	"net/http"
	"runtime/debug"
	"strconv"
)

type PaginateQueryRequest struct {
	Query string `json:"q"`
}

type PaginateFileResponse struct {
	Pagination struct {
		AllRecords int `json:"allRecords"`
	} `json:"pagination"`
	Files []model2.File `json:"files"`
}

func FullImg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var imgInDb model2.File
	model2.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)
	var mimeInDb model2.Mime
	model2.DB.First(&mimeInDb, imgInDb.MimeId)

	writeRawFileApi(w, r, imgInDb.FilePath, mimeInDb.Mime)
}

func Thumb(w http.ResponseWriter, r *http.Request) {
	debug.FreeOSMemory()
	vars := mux.Vars(r)

	var imgInDb model2.File
	model2.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)
	var mimeInDb model2.Mime
	model2.DB.First(&mimeInDb, imgInDb.MimeId)

	// FIXME check if img exists in DB

	width, _ := strconv.ParseUint(vars["w"], 10, 32)
	height, _ := strconv.ParseUint(vars["h"], 10, 32)

	// create thumb on disk if needed
	done := make(chan bool)
	go core.CreateThumb(imgInDb.FilePath, imgInDb.Sha256, mimeInDb.Mime, uint(width), uint(height), done, true)
	<-done
	debug.FreeOSMemory()

	// push thumb to browser, thumb will be always .jpg
	writeRawFileApi(w, r, core.ThumbPath(imgInDb.Sha256, uint(width), uint(height)), "image/jpeg")
}

func FileSimilar(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	rawRes := model2.SimilarFiles(vars["sha256"])

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

func FilePaginate(w http.ResponseWriter, r *http.Request) {
	//authUserOk, userInfo := CheckApiAuth(w, r)
	//if !authUserOk {
	//	w.WriteHeader(http.StatusUnauthorized)
	//	return
	//}
	//userId := userInfo.Id
	userId := 0

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
	dbFiles, allRecords := model2.FileListPaginate(userId, limit, nextId, prevId, searchTerm)

	rawRes.Pagination.AllRecords = allRecords
	rawRes.Files = dbFiles
	// prevent null result in JSON, make empty array instead
	if allRecords < 1 {
		rawRes.Files = []model2.File{}
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
