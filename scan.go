package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"runtime/debug"
	"github.com/nfnt/resize"
)

func scan(dir string, scanFunction filepath.WalkFunc) {
	filepath.Walk(dir, scanFunction)
	//err := filepath.Walk(root, visit)
	//fmt.Printf("filepath.Walk() returned %v\n", err)
}

// http://www.mrwaggel.be/post/generate-md5-hash-of-a-file-in-golang/
func hashFileMd5(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnMD5String string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnMD5String, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := md5.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnMD5String, err
	}

	//Get the 16 bytes hash
	hashInBytes := hash.Sum(nil)[:16]

	//Convert the bytes to a string
	returnMD5String = hex.EncodeToString(hashInBytes)

	return returnMD5String, nil

}

func hashFileSha1(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnString string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnString, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha1.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnString, err
	}

	//Get the 20 bytes hash
	hashInBytes := hash.Sum(nil)[:20]

	//Convert the bytes to a string
	returnString = hex.EncodeToString(hashInBytes)

	return returnString, nil

}

func hashFileSha256(filePath string) (string, error) {
	//Initialize variable returnMD5String now in case an error has to be returned
	var returnString string

	//Open the passed argument and check for any error
	file, err := os.Open(filePath)
	if err != nil {
		return returnString, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//Open a new hash interface to write to
	hash := sha256.New()

	//Copy the file in the hash interface and check for any error
	if _, err := io.Copy(hash, file); err != nil {
		return returnString, err
	}

	//Get the 32 bytes hash
	hashInBytes := hash.Sum(nil)[:32]

	//Convert the bytes to a string
	returnString = hex.EncodeToString(hashInBytes)

	return returnString, nil

}

func getType(filePath string) (string, error) {
	var returnString string

	file, err := os.Open(filePath)
	if err != nil {
		return returnString, err
	}

	//Tell the program to call the following function when the current function returns
	defer file.Close()

	//var header [512]byte
	header := make([]byte, 512)
	_, err = io.ReadFull(file, header[:])
	if err != nil {
		return returnString, err
	}

	//buff := make([]byte, 512) // docs tell that it take only first 512 bytes into consideration
	//if _, err = io.Copy(buff,file); err != nil {
	//	fmt.Println(err) // do something with that error
	//	return
	//}

	//fmt.Println() // do something based on your detection.
	returnString = http.DetectContentType(header)
	return returnString, nil
}

func thumbDirPath(sha256 string) (path string) {
	parent := "./cache"
	lvl1 := string(sha256[0]) + string(sha256[1])
	lvl2 := string(sha256[2]) + string(sha256[3])
	lvl3 := string(sha256[4]) + string(sha256[5])
	path = parent + "/" + lvl1 + "/" + lvl2 + "/" + lvl3 + "/"
	return path
}

func thumbPath(sha256 string, width uint, height uint) (path string) {
	dir := thumbDirPath(sha256)
	path = dir + sha256 + "_" + string(width) + "_" + string(height)
	return path
}

func createThumb(pathToImg string, sha256 string, mime string, maxWidth uint, maxHeight uint, done chan bool) {
	// don't overwrite current file
	// TODO check if image is ok then rewrite if needed
	imgThumbPath := thumbPath(sha256, maxWidth, maxHeight)
	if _, err := os.Stat(imgThumbPath); !os.IsNotExist(err) {
		done <- true
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

	done <- true
}

func makeThumbs(path string, sha256sum string, mime string) {
	done1 := make(chan bool)
	done2 := make(chan bool)
	go createThumb(path, sha256sum, mime, 300, 250, done1)
	go createThumb(path, sha256sum, mime, 560, 500, done2)
	<-done1
	<-done2
	debug.FreeOSMemory()
}

func visit(db *sql.DB, generateThumbs bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		fmt.Printf("Visiting: %s\n", path)

		// don't continue if folder
		info, nil := os.Stat(path)
		if info.IsDir() {
			return nil
		}

		// just create thumbs if specified
		sha256sum, _ := hashFileSha256(path)
		fmt.Printf("SHA256: %s\n", sha256sum)
		alreadyInDb, _, _, mime := dbFindSha256(db, sha256sum)

		// thumbs for files already in DB
		if generateThumbs && alreadyInDb {
			makeThumbs(path, sha256sum, mime)
		}

		// end here for files present in DB
		if alreadyInDb {
			return nil
		}

		// TODO update lastPath even in found in db, use size for fast check

		name := info.Name()
		fmt.Printf("Name: %s\n", name)
		size := info.Size()
		fmt.Printf("Size: %d\n", size)
		mime, _ = getType(path)
		fmt.Printf("MIME: %s\n", mime)
		md5sum, _ := hashFileMd5(path)
		fmt.Printf("MD5: %s\n", md5sum)
		sha1sum, _ := hashFileSha1(path)
		fmt.Printf("SHA1: %s\n", sha1sum)

		fmt.Printf("%s", "Writing to DB...")
		dbFindSert(db, path, size, mime, md5sum, sha1sum, sha256sum)
		fmt.Printf("%s", " done\n")

		// thumbs for new files
		makeThumbs(path, sha256sum, mime)

		return nil
	}
}
