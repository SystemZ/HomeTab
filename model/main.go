package model

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"

	"github.com/DavidHuie/gomigrate"
	"github.com/carlogit/phash"
	_ "github.com/mattn/go-sqlite3"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func DbInit() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./gotag.sqlite3")
	checkErr(err)

	//change to https://github.com/mattes/migrate
	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Sqlite3{}, "./migrations")
	err = migrator.Migrate()

	return db
}

type File struct {
	Fid      int    `json:"id"`
	Name     string `json:"lastPath"`
	Size     int    `json:"size"`
	Sha256   string `json:"sha256"`
	Mime     string `json:"mime"`
	ParentId int    `json:"parentId"`
	Path     string `json:"path"`
	Phash    string `json:"phash"`
}

type Files []File

func List(db *sql.DB, page int64) (found bool, sha256s map[int]File) {
	found = false

	limit := int64(15)
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	// query
	rows, err := db.Query("SELECT id, last_path, size, sha256, phash FROM files ORDER BY id LIMIT ?, ?", offset, limit)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var last_path string
		var size int
		var sha256 string
		var phash sql.NullString
		err = rows.Scan(&id, &last_path, &size, &sha256, &phash)

		filename := filepath.Base(last_path)

		checkErr(err)

		if sha256s == nil {
			sha256s = make(map[int]File)
		}
		sha256s[id] = File{Fid: id, Name: filename, Size: size, Sha256: sha256, Path: last_path, Phash: phash.String}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func ListRandom(db *sql.DB, page int64) (found bool, sha256s map[int]File) {
	found = false

	limit := int64(15)
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	// query
	rows, err := db.Query("SELECT id, last_path, size, sha256 FROM files ORDER BY random() LIMIT ?, ?", offset, limit)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var last_path string
		var size int
		var sha256 string
		err = rows.Scan(&id, &last_path, &size, &sha256)

		filename := filepath.Base(last_path)

		checkErr(err)

		if sha256s == nil {
			sha256s = make(map[int]File)
		}
		sha256s[id] = File{Fid: id, Name: filename, Size: size, Sha256: sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func ListSha256(db *sql.DB, search string) (found bool, sha256s map[int]File) {
	found = false
	rows, err := db.Query("SELECT id, last_path, size, sha256 FROM files WHERE sha256 = ? ORDER BY id LIMIT 100", search)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var last_path string
		var size int
		var sha256 string
		err = rows.Scan(&id, &last_path, &size, &sha256)

		filename := filepath.Base(last_path)

		checkErr(err)

		if sha256s == nil {
			sha256s = make(map[int]File)
		}
		sha256s[id] = File{Fid: id, Name: filename, Size: size, Sha256: sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func Find(db *sql.DB, sha256 string) (found bool, file File) {
	rows, err := db.Query("SELECT id, last_path, size FROM files WHERE sha256 = ?", sha256)
	checkErr(err)
	defer rows.Close()
	var id int
	var last_path string
	var size int
	found = false

	for rows.Next() {
		err = rows.Scan(&id, &last_path, &size)
		checkErr(err)
		found = true
		file = File{Fid: id, Name: last_path, Size: size, Sha256: sha256}
		break
	}
	return found, file
}

func FindById(db *sql.DB, fileId int) (found bool, file File) {
	rows, err := db.Query("SELECT last_path, size, sha256 FROM files WHERE id = ?", fileId)
	checkErr(err)
	defer rows.Close()
	var last_path string
	var size int
	var sha256 string
	found = false

	for rows.Next() {
		err = rows.Scan(&last_path, &size, &sha256)
		checkErr(err)
		found = true
		file = File{Fid: fileId, Name: last_path, Size: size, Sha256: sha256}
		break
	}
	return found, file
}

func FindSha256(db *sql.DB, sha string) (found bool, res int, lastPath string, mime string) {
	rows, err := db.Query("SELECT id, last_path, mime FROM files WHERE sha256 = ?", sha)
	checkErr(err)
	defer rows.Close()
	var id int
	var sha256 string

	var result int
	found = false

	for rows.Next() {
		err = rows.Scan(&id, &sha256, &mime)
		checkErr(err)
		result = id
		lastPath = sha256
		found = true
		break
	}
	return found, result, lastPath, mime
}

func FindByFile(db *sql.DB, filePath string) (found bool, result File) {
	rows, err := db.Query("SELECT id, size FROM files WHERE last_path = ?", filePath)
	checkErr(err)
	defer rows.Close()
	found = false
	var f File

	for rows.Next() {
		err = rows.Scan(&f.Fid, &f.Size)
		checkErr(err)
		return true, f
		break
	}
	return found, File{}
}

func Insert(db *sql.DB, lastPath string, size int64, mime string, sha256 string) (id int64) {
	stmt, err := db.Prepare("INSERT INTO files(last_path, size, mime, sha256) VALUES(?,?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(lastPath, size, mime, sha256)
	checkErr(err)

	id, err = res.LastInsertId()
	checkErr(err)

	return id
}

func FindSert(db *sql.DB, lastPath string, size int64, mime string, sha256 string) {
	found, _ := Find(db, sha256)
	if !found {
		Insert(db, lastPath, size, mime, sha256)
	}
}

func UpdatePath(db *sql.DB, sha256sum string, newPath string) {
	trashSQL, err := db.Prepare("UPDATE files SET last_path=? WHERE sha256=?")
	if err != nil {
		fmt.Println(err)
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	_, err = tx.Stmt(trashSQL).Exec(newPath, sha256sum)
	if err != nil {
		fmt.Println("Doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func UpdateParentId(db *sql.DB, sha256sum string, parentId int) {
	trashSQL, err := db.Prepare("UPDATE files SET parent_id=? WHERE sha256=?")
	if err != nil {
		fmt.Println(err)
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	_, err = tx.Stmt(trashSQL).Exec(parentId, sha256sum)
	if err != nil {
		fmt.Println("Doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

func FindPHash(db *sql.DB, sha256 string) (found bool) {
	rows, err := db.Query("SELECT phash FROM files WHERE phash IS NOT NULL AND sha256 = ?", sha256)
	defer rows.Close()

	checkErr(err)
	found = false

	for rows.Next() {
		//err = rows.Scan(&phash)
		checkErr(err)
		found = true
		break
	}
	return found //, phash
}

func UpdatePHash(db *sql.DB, sha256sum string, pHash string) {
	trashSQL, err := db.Prepare("UPDATE files SET phash=? WHERE sha256=?")
	if err != nil {
		fmt.Println(err)
	}
	tx, err := db.Begin()
	if err != nil {
		fmt.Println(err)
	}
	_, err = tx.Stmt(trashSQL).Exec(pHash, sha256sum)
	if err != nil {
		log.Println("Doing rollback")
		tx.Rollback()
	} else {
		tx.Commit()
	}
}

type Distance struct {
	IdA  int
	IdB  int
	Dist int
}

func FilesWithPHash(db *sql.DB, page int64, sha256sum string) (distances map[int]Distance) {
	distances = make(map[int]Distance)

	rowsOneFile, err := db.Query("SELECT id, phash FROM files WHERE sha256 = ? AND phash IS NOT NULL", sha256sum)
	checkErr(err)

	var id int
	var pHash string
	for rowsOneFile.Next() {
		err = rowsOneFile.Scan(&id, &pHash)
		checkErr(err)
		break
	}
	rowsOneFile.Close()

	limit := int64(1000000) //FIXME
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	rows, err := db.Query("SELECT id, phash FROM files WHERE sha256 != ? AND phash IS NOT NULL AND mime IN ('image/jpeg', 'image/png') ORDER BY id LIMIT ?, ?", sha256sum, offset, limit)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var rId int
		var rPHash string
		err = rows.Scan(&rId, &rPHash)

		checkErr(err)

		if id != rId && id != 0 {
			distance := phash.GetDistance(pHash, rPHash)
			/*
				writes when reading sqlite causes locks and crashes
				aggregate to map then write outside this function
			*/
			distances[rId] = Distance{id, rId, distance}
		}
	}
	return distances
}

func CountAllFiles(db *sql.DB) (filesCounted int) {
	rows, err := db.Query("SELECT COUNT(id) FROM files")
	defer rows.Close()
	checkErr(err)

	for rows.Next() {
		err = rows.Scan(&filesCounted)
		checkErr(err)
		break
	}
	return filesCounted
}
