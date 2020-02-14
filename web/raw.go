package web

import (
	"bytes"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

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