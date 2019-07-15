package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"

	"gitlab.com/systemz/gotag/model"
)

func scan(dir string, scanFunction filepath.WalkFunc) {
	filepath.Walk(dir, scanFunction)
	//err := filepath.Walk(root, visit)
	//fmt.Printf("filepath.Walk() returned %v\n", err)
}

func scanNg(db *sql.DB, dir string) {
	//log.Printf("Will work with %v threads", runtime.GOMAXPROCS(0))
	log.Printf("Scanning...")
	var fileList []string
	_ = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	log.Printf("Scanning complete")

	//maxGoroutines := 2
	maxGoroutines := runtime.GOMAXPROCS(0)
	guard := make(chan struct{}, maxGoroutines)
	// other implementation of concurrency:
	// https://stackoverflow.com/questions/43789362/is-it-possible-to-limit-how-many-goroutines-run-per-second/43792222#43792222

	// limit concurrency
	start := time.Now()
	fmt.Println("getting ready to do some work...")

	for _, file := range fileList {
		log.Printf("%v", file)
		guard <- struct{}{} // would block if guard channel is already filled

		go func() {
			log.Printf("Starting work...")
			addFile(db, file, AddFileOptions{
				generateThumbs: true,
				calcSimilarity: true,
				onlyAddNew:     true,
			})
			fmt.Println("Finished work:", time.Now())
			<-guard
		}()
	}

	dur := time.Since(start)
	fmt.Println("scanned in", dur)

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
	path = dir + sha256 + "_" + strconv.Itoa(int(width)) + "_" + strconv.Itoa(int(height))
	return path
}

func makeThumbs(path string, sha256sum string, mime string) {
	done1 := make(chan bool)
	done2 := make(chan bool)
	go createThumb(path, sha256sum, mime, 300, 300, done1)
	go createThumb(path, sha256sum, mime, 700, 700, done2)
	<-done1
	<-done2
	debug.FreeOSMemory()
}

type AddFileOptions struct {
	onlyAddNew     bool
	generateThumbs bool
	calcSimilarity bool
	tags           []string
	parentId       int
}

func addFile(db *sql.DB, path string, options AddFileOptions) (dbFile model.File) {
	log.Printf("Checking: %s\n", path)

	// don't continue if folder
	info, nil := os.Stat(path)
	if info.IsDir() {
		log.Printf("%s\n", "It's a dir, skipping...")
		return
	}

	fileInfo, err := os.Stat(path)
	size := fileInfo.Size()
	log.Printf("Size: %d\n", size)

	if options.onlyAddNew {
		isInDb, fileInDb := model.FindByFile(db, path)
		if isInDb && fileInDb.Size == int(size) {
			// this file already exists in DB, skip it
			log.Printf("Skipping...")
			return
		}
	}

	//TODO use size for fast check
	log.Printf("%s\n", "Calculating SHA256...")
	sha256sum, _ := hashFileSha256(path)
	log.Printf("SHA256: %s\n", sha256sum)
	isInDb, _, _, mime := model.FindSha256(db, sha256sum)

	// check if it's already in DB
	_, fileInDb := model.Find(db, sha256sum)

	// add new file to DB if needed
	if !isInDb {
		log.Printf("%s\n", "File not in DB, check info and add to DB!")

		if err != nil {
			log.Printf("Error when getting info for %v: %v", path, err)
			return
		}

		name := fileInfo.Name()
		log.Printf("Name: %s\n", name)
		mime, _ = getType(path)
		log.Printf("MIME: %s\n", mime)

		log.Printf("%s", "Writing to DB...")
		model.FindSert(db, path, size, mime, sha256sum)
		log.Printf("%s", "...done\n")
	}

	// apply tags if provided
	for _, tag := range options.tags {
		model.TagFindSert(db, tag, fileInDb.Fid)
	}

	// update path to file if necessary
	if isInDb && fileInDb.Name != path {
		log.Printf("Updating path from %s to %s ...", fileInDb.Name, path)
		model.UpdatePath(db, sha256sum, path)
		log.Printf("%s\n", "Updating path done")
	} else if isInDb && fileInDb.Name == path {
		log.Printf("%s\n", "File path is up to date")
	}

	//// update parent id
	//if isInDb && fileInDb.ParentId != options.parentId {
	//	model.UpdateParentId(db, fileInDb.Sha256, options.parentId)
	//}

	// check similarity to other images
	if options.calcSimilarity {
		log.Println("Calculating similarity to other images, this may take a while")
		// calc and add perceptual hash to DB for images
		pHashFound := model.FindPHash(db, sha256sum)
		if (mime == "image/jpeg" || mime == "image/png") && !pHashFound {
			log.Printf("%s\n", "pHash not found, calculating...")
			pHash := getPHash(path)
			log.Printf("%s %s\n", "pHash:", pHash)
			model.UpdatePHash(db, sha256sum, pHash)
		}

		// calc distance between this and rest of images
		// FIXME support more than 1m rows with count rows before starting
		//start := timeStart()
		distances := model.FilesWithPHash(db, 1, sha256sum)
		tx, stmt := model.DistanceInsertPrepare(db)
		//i := 0
		for _, v := range distances {
			model.DistanceInsert(stmt, v.IdA, v.IdB, v.Dist)
			//i++
		}
		model.DistanceInsertEnd(tx)
		log.Println("Calculating similarity done")
	}

	// thumbs creation
	if options.generateThumbs {
		log.Printf("%s\n", "Creating thumbs...")
		makeThumbs(path, sha256sum, mime)
		log.Printf("%s\n", "Thumbs work done")
	}

	return fileInDb
}

// code below is almost duplicate of addFile()

func visit(db *sql.DB, generateThumbs bool) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		log.Printf("Visiting: %s\n", path)

		// don't continue if folder
		info, nil := os.Stat(path)
		if info.IsDir() {
			log.Printf("%s\n", "It's a dir, skipping...")
			return nil
		}

		// just create thumbs if specified
		log.Printf("%s\n", "Calculating SHA256...")
		//TODO use size for fast check
		sha256sum, _ := hashFileSha256(path)
		log.Printf("SHA256: %s\n", sha256sum)
		alreadyInDb, _, _, mime := model.FindSha256(db, sha256sum)

		// thumbs for files already in DB
		if generateThumbs && alreadyInDb {
			log.Printf("%s\n", "Creating thumbs for file already in DB...")
			makeThumbs(path, sha256sum, mime)
			log.Printf("%s\n", "Thumbs work done")
		} else {
			log.Printf("%s\n", "Thumb generation not enabled, skipping...")
		}

		_, fileInDb := model.Find(db, sha256sum)
		if alreadyInDb && fileInDb.Name != path {
			log.Printf("Updating path from %s to %s ...", fileInDb.Name, path)
			model.UpdatePath(db, sha256sum, path)
			log.Printf("%s\n", "Updating path done")
		} else if alreadyInDb && fileInDb.Name == path {
			log.Printf("%s\n", "File path is up to date")
		}

		log.Println("Calculating similarity to other images, this may take a while")
		// calc and add perceptual hash to DB for images
		pHashFound := model.FindPHash(db, sha256sum)
		if (mime == "image/jpeg" || mime == "image/png") && !pHashFound {
			log.Printf("%s\n", "pHash not found, calculating...")
			pHash := getPHash(path)
			log.Printf("%s %s\n", "pHash:", pHash)
			model.UpdatePHash(db, sha256sum, pHash)
		}

		// calc distance between this and rest of images
		// FIXME support more than 1m rows with count rows before starting
		//start := timeStart()
		distances := model.FilesWithPHash(db, 1, sha256sum)
		tx, stmt := model.DistanceInsertPrepare(db)
		//i := 0
		for _, v := range distances {
			model.DistanceInsert(stmt, v.IdA, v.IdB, v.Dist)
			//i++
		}
		model.DistanceInsertEnd(tx)
		log.Println("Calculating similarity done")
		//timeStop(start)
		//fmt.Printf("%d queries\n", i)

		// end here for files present in DB
		if alreadyInDb {
			log.Printf("%s\n", "File already in DB, ending...")
			return nil
		}

		log.Printf("%s\n", "File not in DB, check info and add to DB!")
		name := info.Name()
		log.Printf("Name: %s\n", name)
		size := info.Size()
		mime, _ = getType(path)
		log.Printf("Size: %d, MIME: %s\n ", size, mime)

		log.Printf("%s", "Writing to DB...")
		model.FindSert(db, path, size, mime, sha256sum)
		log.Printf("%s", "...done\n")

		// thumbs for new files
		makeThumbs(path, sha256sum, mime)

		return nil
	}
}
