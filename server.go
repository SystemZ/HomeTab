package main

import (
	"bytes"
	"database/sql"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"os"
	"strconv"

	"image/png"

	"strings"

	"runtime/debug"

	"github.com/nfnt/resize"
	"github.com/pilu/traffic"
	"github.com/unrolled/render"
)

type FileView struct {
	File File
	Tags map[int]Tag
}

type FilesView struct {
	Files    map[int]File
	PrevPage int64
	NextPage int64
}

func writeImage(w http.ResponseWriter, img *image.Image) {
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func writeCache(key string, r *traffic.Request, w http.ResponseWriter) {
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

func writeImageWithCache(w http.ResponseWriter, r *traffic.Request, img *image.Image) {
	writeCache("foobar", r, w)
	buffer := new(bytes.Buffer)
	if err := jpeg.Encode(buffer, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
		*img = nil
	} else {
		*img = nil
	}
}

func writeImageDirect(w http.ResponseWriter, img *image.Image) {
	w.Header().Set("Content-Type", "image/jpeg")
	if err := jpeg.Encode(w, *img, nil); err != nil {
		log.Println("unable to encode image.")
	}
}

func writeRawFile(w http.ResponseWriter, r *traffic.Request, filePath string, mime string) {
	writeCache("foo", r, w)

	imgFile, _ := os.Open(filePath)
	defer imgFile.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(imgFile)

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}
}

func server(db *sql.DB) {
	r := render.New(render.Options{
		IndentJSON: true,
		Extensions: []string{".tmpl", ".html"},
	})
	router := traffic.New()
	traffic.SetHost("0.0.0.0")

	router.Get("/", func(w traffic.ResponseWriter, req *traffic.Request) {
		r.HTML(w, http.StatusOK, "home", nil)
	})

	router.Get("/tags", func(w traffic.ResponseWriter, req *traffic.Request) {
		_, tags := dbTagRank(db)
		r.HTML(w, http.StatusOK, "tag_rank", tags)
	})

	router.Get("/files/random", func(w traffic.ResponseWriter, req *traffic.Request) {
		_, files := dbListRandom(db, 1)
		r.HTML(w, http.StatusOK, "random_files", files)
	})

	router.Get("/files/random/notag", func(w traffic.ResponseWriter, req *traffic.Request) {
		_, files := dbFilesWithoutTagsRandom(db, 1)
		r.HTML(w, http.StatusOK, "random_files", files)
	})

	router.Get("/files/tag/name/:tag/page/:page", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		page, _ := strconv.ParseInt(params.Get("page"), 10, 64)
		_, files := dbFileTagSearchByName(db, page, params.Get("tag"))
		prevPage := page - 1
		if page == 1 {
			prevPage = 1
		}
		r.HTML(w, http.StatusOK, "files", FilesView{files, prevPage, page + 1})
	})

	router.Get("/files/page/:page", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		page, _ := strconv.ParseInt(params.Get("page"), 10, 64)
		_, files := dbList(db, page)
		prevPage := page - 1
		if page == 1 {
			prevPage = 1
		}
		r.HTML(w, http.StatusOK, "files", FilesView{files, prevPage, page + 1})
	})

	router.Get("/files/notag/page/:page", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		page, _ := strconv.ParseInt(params.Get("page"), 10, 64)
		_, files := dbFilesWithoutTags(db, page)
		prevPage := page - 1
		if page == 1 {
			prevPage = 1
		}
		r.HTML(w, http.StatusOK, "files", FilesView{files, prevPage, page + 1})
	})

	router.Get("/file/:sha256", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		_, file := dbFind(db, params.Get("sha256"))
		_, tags := dbTagList(db, file.Fid)
		r.HTML(w, http.StatusOK, "file", FileView{file, tags})
	})

	router.Post("/file/:sha256/tag/add", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		found, files := dbListSha256(db, params.Get("sha256"))
		if !found {
			http.Error(w, "File doesn't exist", 404)
		}
		tag := req.FormValue("name")

		for k := range files {
			dbTagFindSert(db, tag, k)
		}

		w.Header().Set("Location", "/file/"+params.Get("sha256"))
		r.Text(w, 302, "")
	})

	router.Post("/file/:sha256/tag/delete", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		tag := req.FormValue("id")
		tagId, _ := strconv.ParseInt(tag, 10, 16)
		dbTagDelete(db, int(tagId))

		w.Header().Set("Location", "/file/"+params.Get("sha256"))
		r.Text(w, 302, "")
	})

	router.Get("/img/full/:sha256", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		_, _, lastPath, mime := dbFindSha256(db, params.Get("sha256"))
		writeRawFile(w, req, lastPath, mime)
	})

	router.Get("/img/thumb/:w/:h/:sha256", func(w traffic.ResponseWriter, req *traffic.Request) {
		params := req.URL.Query()
		sha256sum := params.Get("sha256")
		width64, _ := strconv.ParseUint(params.Get("w"), 10, 32)
		height64, _ := strconv.ParseUint(params.Get("h"), 10, 32)
		width := uint(width64)
		height := uint(height64)
		_, _, imgPath, mime := dbFindSha256(db, sha256sum)

		// create thumb on disk if needed
		done := make(chan bool)
		go createThumb(imgPath, sha256sum, mime, width, height, done)
		<-done
		debug.FreeOSMemory()

		// push thumb to browser
		writeRawFile(w, req, thumbPath(sha256sum, width, height), mime)
	})

	router.Get("/img/thumbs/:w/:h/:sha256", func(w traffic.ResponseWriter, req *traffic.Request) {
		//debug.FreeOSMemory()
		params := req.URL.Query()
		var img image.Image
		_, _, lastPath, mime := dbFindSha256(db, params.Get("sha256"))

		width, _ := strconv.ParseUint(params.Get("w"), 10, 32)
		height, _ := strconv.ParseUint(params.Get("h"), 10, 32)

		imgFile, _ := os.Open(lastPath)
		defer imgFile.Close()

		switch mime {
		case "image/jpeg":
			img, _ = jpeg.Decode(imgFile)
		case "image/png":
			img, _ = png.Decode(imgFile)
		case "image/gif":
			writeRawFile(w, req, lastPath, mime)
			//TODO config for gif thumbs
			//img, _ = gif.Decode(imgFile)
			//gif1, _ := gif.DecodeAll(imgFile)
		//case "video/webm":
		//	img, _ = webm.Decode(imgFile)
		//case "video/mp4":
		//	img, _ = mp4.Decode(imgFile)
		default:
			http.Error(w, "This file type is not supported yet, sorry :(", 500)
			return
		}
		//
		log.Printf("%v", img.ColorModel())
		//
		thumb := resize.Thumbnail(uint(width), uint(height), img, resize.Bilinear)
		//debug.FreeOSMemory()
		writeImageWithCache(w, req, &thumb)
	})

	router.Get("/", func(w traffic.ResponseWriter, req *traffic.Request) {
		r.JSON(w, http.StatusOK, map[string]string{"welcome": "This is rendered JSON!"})
	})

	router.Run()
}
