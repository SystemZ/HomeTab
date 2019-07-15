package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"fmt"
	"io/ioutil"

	"github.com/discordapp/lilliput"
	"github.com/nfnt/resize"
)

var EncodeOptions = map[string]map[int]int{
	".jpeg": map[int]int{lilliput.JpegQuality: 85},
	".png":  map[int]int{lilliput.PngCompression: 7},
	".webp": map[int]int{lilliput.WebpQuality: 85},
}

// https://github.com/discordapp/lilliput/blob/c86952445060cf68db423d77ce52ee4d4dbc819c/examples/main.go
func createVideoThumb(inputFilename string, sha256 string, mime string, maxWidth uint, maxHeight uint) {

	outputFilename := thumbPath(sha256, maxWidth, maxHeight)

	outputWidth := int(maxWidth)
	outputHeight := int(maxHeight)

	// decoder wants []byte, so read the whole file into a buffer
	inputBuf, err := ioutil.ReadFile(inputFilename)
	decoder, err := lilliput.NewDecoder(inputBuf)
	// this error reflects very basic checks,
	// mostly just for the magic bytes of the file to match known image formats
	if err != nil {
		fmt.Printf("error decoding image, %s\n", err)
		// skip this file
		return
	}
	defer decoder.Close()

	header, err := decoder.Header()
	// this error is much more comprehensive and reflects
	// format errors
	if err != nil {
		fmt.Printf("error reading image header, %s\n", err)
		// skip this file
		return
	}

	// print some basic info about the image
	fmt.Printf("image type: %s\n", decoder.Description())
	fmt.Printf("%dpx x %dpx\n", header.Width(), header.Height())

	// get ready to resize image,
	// using 8192x8192 maximum resize buffer size
	ops := lilliput.NewImageOps(8192)
	defer ops.Close()

	// create a buffer to store the output image, 50MB in this case
	outputImg := make([]byte, 50*1024*1024)
	outputType := ".jpg"
	//outputType := "." + strings.ToLower(decoder.Description())

	resizeMethod := lilliput.ImageOpsFit
	//if stretch {
	//resizeMethod := lilliput.ImageOpsResize
	//}

	opts := &lilliput.ImageOptions{
		FileType:             outputType,
		Width:                outputWidth,
		Height:               outputHeight,
		ResizeMethod:         resizeMethod,
		NormalizeOrientation: true,
		EncodeOptions:        EncodeOptions[outputType],
	}

	// resize and transcode image
	outputImg, err = ops.Transform(decoder, opts, outputImg)
	if err != nil {
		fmt.Printf("error transforming image, %s\n", err)
		// skip this file
		return
	}

	// image has been resized, now write file out
	err = ioutil.WriteFile(outputFilename, outputImg, 0664)
	if err != nil {
		fmt.Printf("error writing out resized image, %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("image written to %s\n", outputFilename)

}

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
	case "image/gif":
		log.Println("Creating .gif thumb...")
	case "video/webm":
		log.Println("Creating .webm thumb...")
	case "video/mp4":
		log.Println("Creating .mp4 thumb...")
	default:
		log.Printf("This MIME is not supported yet: %v", mime)
		done <- true
		// exit this function if not jpg,png,webm or mp4
		return
	}

	// create path for thumb
	dirPath := thumbDirPath(sha256)
	os.MkdirAll(dirPath, 0700)

	if mime == "image/jpeg" || mime == "image/png" {
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
	} else if mime == "image/gif" || mime == "video/webm" || mime == "video/mp4" {
		createVideoThumb(pathToImg, sha256, mime, maxWidth, maxHeight)
	}

	log.Printf("Creating %d x %d thumb done\n", maxWidth, maxHeight)
	done <- true
}
