// +build windows,amd64 linux,arm

package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

func createThumb(pathToImg string, sha256 string, mime string, maxWidth uint, maxHeight uint, done chan bool) {
	// don't overwrite current file
	// TODO check if image is ok then rewrite if needed
	imgThumbPath := thumbPath(sha256, maxWidth, maxHeight)
	if _, err := os.Stat(imgThumbPath); !os.IsNotExist(err) {
		log.Printf("Creating %d x %d thumb is not needed, skipping...\n", maxWidth, maxHeight)
		done <- true
		return
	} else {
		log.Printf("Creating %d x %d thumb...\n", maxWidth, maxHeight)
	}

	// check if we still have original file on disk, prevent crash
	if _, err := os.Stat(pathToImg); os.IsNotExist(err) {
		log.Printf("File %s missing, skipping...\n", pathToImg)
		done <- true
		return
	}

	// open original file to read
	fileRead, err := os.Open(pathToImg)
	if err != nil {
		log.Fatalln(err)
	}
	defer fileRead.Close()

	// decide what type of file we are reading (decoding) to prevent creating useless files
	var img image.Image
	switch mime {
	case "image/jpeg":
		img, _ = jpeg.Decode(fileRead)
	case "image/png":
		img, _ = png.Decode(fileRead)
	//case "video/webm":
	//	img, _ = webm.Decode(fileRead)
	//case "video/mp4":
	//	img, _ = mp4.Decode(fileRead)
	default:
		log.Printf("%s", "This file type is not supported yet, sorry :(")
		done <- true
		// exit this function if not jpg,png,webm or mp4
		return
	}

	// create path for thumb
	dirPath := thumbDirPath(sha256)
	os.MkdirAll(dirPath, 0700)

	// open thumb file to write
	fileWrite, err := os.Create(thumbPath(sha256, maxWidth, maxHeight))
	if err != nil {
		log.Fatalln(err)
	}
	defer fileWrite.Close()

	// resize original file to something smaller
	thumb := resize.Thumbnail(maxWidth, maxHeight, img, resize.Bilinear)

	// create thumbnail in jpg
	err = jpeg.Encode(fileWrite, thumb, nil)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Creating %d x %d thumb done\n", maxWidth, maxHeight)
	done <- true
}
