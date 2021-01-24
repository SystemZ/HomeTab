package server

import (
	"bytes"
	"github.com/gorilla/mux"
	"github.com/systemz/hometab/internal/model"
	"github.com/systemz/hometab/internal/service/gotagcore"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"strings"
)

// TODO check for auth
// TODO check if user have this file in collection

func FullImg(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var imgInDb model.File
	model.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)
	var mimeInDb model.Mime
	model.DB.First(&mimeInDb, imgInDb.MimeId)

	writeRawFileApi(w, r, imgInDb.FilePath, mimeInDb.Mime)
}

func Thumb(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	var imgInDb model.File
	model.DB.Where("sha256 = ?", vars["sha256"]).First(&imgInDb)
	var mimeInDb model.Mime
	model.DB.First(&mimeInDb, imgInDb.MimeId)

	// FIXME check if img exists in DB

	width, _ := strconv.ParseUint(vars["w"], 10, 32)
	height, _ := strconv.ParseUint(vars["h"], 10, 32)

	// create thumb on disk if needed
	done := make(chan bool)
	go gotagcore.CreateThumb(imgInDb.FilePath, imgInDb.Sha256, mimeInDb.Mime, uint(width), uint(height), done, true)
	<-done
	debug.FreeOSMemory()

	// push thumb to browser, thumb will be always .jpg
	writeRawFileApi(w, r, gotagcore.ThumbPath(imgInDb.Sha256, uint(width), uint(height)), "image/jpeg")
}

func writeRawFileApi(w http.ResponseWriter, r *http.Request, filePath string, mime string) {
	writeCacheApi("foo", r, w)

	imgFile, _ := os.Open(filePath)
	defer imgFile.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(imgFile)

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image via writeRawFile")
	}
}

func writeCacheApi(key string, r *http.Request, w http.ResponseWriter) {
	e := `"` + key + `"`
	w.Header().Set("Etag", e)
	//w.Header().Set("Cache-Control", "max-age=2592000") // 30 days
	w.Header().Set("Cache-Control", "max-age=60") // 1 minute for easier dev

	if match := r.Header.Get("If-None-Match"); match != "" {
		if strings.Contains(match, e) {
			w.WriteHeader(http.StatusNotModified)
			return
		}
	}
}
