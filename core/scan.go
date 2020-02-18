package core

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/jinzhu/gorm"
	"gitlab.com/systemz/gotag/model"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"runtime/debug"
	"strconv"
	"time"
)

func ScanMulti(db *gorm.DB, dir string) {
	// search for folders and files recursively
	log.Printf("Searching for files to scan...")
	var fileList []string
	_ = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	log.Printf("Files and folders found: %v", len(fileList))

	// set how many threads to use
	maxGoroutines := runtime.GOMAXPROCS(0)
	log.Printf("Scanning with %v threads ...", maxGoroutines)
	guard := make(chan struct{}, maxGoroutines)

	// other implementation of concurrency:
	// https://stackoverflow.com/questions/43789362/is-it-possible-to-limit-how-many-goroutines-run-per-second/43792222#43792222

	// limit concurrency
	start := time.Now()
	fmt.Println("Starting file scan...")

	for _, file := range fileList {
		fileStart := time.Now()
		log.Printf("Scanning: %v", file)
		guard <- struct{}{} // would block if guard channel is already filled

		go func() {
			AddFile(db, file, AddFileOptions{
				GenerateThumbs: true,
				CalcSimilarity: true,
				OnlyAddNew:     true,
			})
			log.Printf("Done in %v: %v ", time.Since(fileStart), file)
			<-guard
		}()
	}

	// the end, show summary
	dur := time.Since(start)
	log.Printf("Scanned in %v", dur)
}

func ScanMono(db *gorm.DB, dir string, userId int) {
	// search for folders and files recursively
	log.Printf("Searching for files to scan...")
	var fileList []string
	_ = filepath.Walk(dir, func(path string, f os.FileInfo, err error) error {
		fileList = append(fileList, path)
		return nil
	})
	log.Printf("Files and folders found: %v", len(fileList))

	// limit concurrency
	start := time.Now()
	fmt.Println("Starting file scan...")

	for _, file := range fileList {
		fileStart := time.Now()
		log.Printf("Scanning: %v", file)
		AddFile(db, file, AddFileOptions{
			GenerateThumbs: true,
			CalcSimilarity: true,
			OnlyAddNew:     true,
			UserId:         userId,
		})
		log.Printf("Done in %v: %v ", time.Since(fileStart), file)
		log.Println("")
	}

	// the end, show summary
	dur := time.Since(start)
	log.Printf("Scanned in %v", dur)
}

type AddFileOptions struct {
	OnlyAddNew     bool
	GenerateThumbs bool
	CalcSimilarity bool
	Tags           []string
	ParentId       int
	UserId         int
}

func AddFile(db *gorm.DB, path string, options AddFileOptions) {
	log.Printf("Checking: %s\n", path)

	// don't continue if folder
	info, _ := os.Stat(path)
	if info.IsDir() {
		log.Printf("Dir, skip: %v", path)
		return
	}

	fileInfo, err := os.Stat(path)
	if err != nil {
		log.Printf("Problem with access, skip: %v", path)
		return
	}
	// TODO use size and filepath for quick scan
	size := fileInfo.Size()
	log.Printf("Size: %d %v", size, path)

	//if options.OnlyAddNew {
	//	isInDb, fileInDb := model.FindByFile(db, path)
	//	if isInDb && fileInDb.Size == int(size) {
	//		// this file already exists in DB, skip it
	//		log.Printf("Skipping...")
	//		return
	//	}
	//}

	////TODO use size for fast check
	log.Printf("%s\n", "Check SHA256...")
	sha256sum, _ := hashFileSha256(path)
	log.Printf("SHA256: %s\n", sha256sum)

	// check if file is already in DB
	var fileInDb model.File
	model.DB.Where("sha256 = ?", sha256sum).First(&fileInDb)

	// add new file to DB if needed
	if fileInDb.Id < 1 {
		log.Printf("%s\n", "File not in DB, adding...")

		// save MIME to DB
		mime, err := getType(path)
		if err != nil {
			log.Printf("Problem with getting MIME, skip: %v", path)
			return
		}
		mimeId := model.AddMime(db, mime)

		// check similarity to other images
		// calc and add perceptual hash to DB for images
		var pHash string
		var pHashA, pHashB, pHashC, pHashD int
		if mime == "image/jpeg" || mime == "image/png" {
			pHash = GetPHash(path)
			log.Printf("%s %s\n", "pHash:", pHash)
			// divide pHash for optimal storage in DB
			pHashA, _ = strconv.Atoi(pHash[0:16])
			pHashB, _ = strconv.Atoi(pHash[16:32])
			pHashC, _ = strconv.Atoi(pHash[32:48])
			pHashD, _ = strconv.Atoi(pHash[48:64])
		}

		// get original filename
		fileName := filepath.Base(path)

		// add file to DB, finally...
		fileInDb = model.File{
			Filename: fileName,
			FilePath: path,
			SizeB:    int(size),
			MimeId:   mimeId,
			PhashA:   pHashA,
			PhashB:   pHashB,
			PhashC:   pHashC,
			PhashD:   pHashD,
			Sha256:   sha256sum,
			Mime:     mime,
		}
		db.Save(&fileInDb)
		log.Printf("%+v", fileInDb)
	} else {
		log.Printf("Already in DB, skip add")
	}

	// add permissions to file if necessary
	var fileUser model.FileUser
	db.Where("file_id = ? AND user_id = ?", fileInDb.Id, options.UserId).First(&fileUser)

	// permissions not found, add to DB
	if fileUser.Id < 1 {
		log.Printf("%s\n", "Permissions not in DB, adding...")
		fileUser = model.FileUser{
			FileId:    fileInDb.Id,
			UserId:    options.UserId,
			CreatedAt: fileInDb.CreatedAt,
			UpdatedAt: fileInDb.UpdatedAt,
			DeletedAt: nil,
		}
		model.DB.Create(&fileUser)
	}

	// add tags to DB
	//found, tags := model.TagList(sqlite, img.Fid)
	//if !found {
	//	// finish work if no tags for this file
	//	continue
	//}
	//for _, tag := range tags {
	//	model.AddTagToFile(mysql, tag.Name, file.Id)
	//}
}

func makeThumbs(path string, sha256sum string, mime string) {
	done1 := make(chan bool)
	done2 := make(chan bool)
	go CreateThumb(path, sha256sum, mime, 300, 300, done1, false)
	go CreateThumb(path, sha256sum, mime, 700, 700, done2, false)
	<-done1
	<-done2
	debug.FreeOSMemory()
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
