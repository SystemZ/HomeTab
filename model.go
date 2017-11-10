package main

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

func dbInit() (db *sql.DB) {
	db, err := sql.Open("sqlite3", "./gotag.sqlite3")
	checkErr(err)

	//change to https://github.com/mattes/migrate
	migrator, _ := gomigrate.NewMigrator(db, gomigrate.Sqlite3{}, "./migrations")
	err = migrator.Migrate()

	return db
}

type File struct {
	Fid    int
	Name   string
	Size   int
	Sha256 string
}

type Files []File

func dbList(db *sql.DB, page int64) (found bool, sha256s map[int]File) {
	found = false

	limit := int64(15)
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	// query
	rows, err := db.Query("SELECT id, last_path, size, sha256 FROM files ORDER BY id LIMIT ?, ?", offset, limit)
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
		sha256s[id] = File{id, filename, size, sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func dbListRandom(db *sql.DB, page int64) (found bool, sha256s map[int]File) {
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
		sha256s[id] = File{id, filename, size, sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func dbListSha256(db *sql.DB, search string) (found bool, sha256s map[int]File) {
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
		sha256s[id] = File{id, filename, size, sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func dbFind(db *sql.DB, sha256 string) (found bool, file File) {
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
		file = File{id, last_path, size, sha256}
		break
	}
	return found, file
}

func dbFindById(db *sql.DB, fileId int) (found bool, file File) {
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
		file = File{fileId, last_path, size, sha256}
		break
	}
	return found, file
}

func dbFindSha256(db *sql.DB, sha string) (found bool, res int, lastPath string, mime string) {
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

func dbInsert(db *sql.DB, lastPath string, size int64, mime string, md5 string, sha1 string, sha256 string) (id int64) {
	stmt, err := db.Prepare("INSERT INTO files(last_path, size, mime, md5, sha1, sha256) VALUES(?,?,?,?,?,?)")
	checkErr(err)

	res, err := stmt.Exec(lastPath, size, mime, md5, sha1, sha256)
	checkErr(err)

	id, err = res.LastInsertId()
	checkErr(err)

	return id
}

func dbFindSert(db *sql.DB, lastPath string, size int64, mime string, md5 string, sha1 string, sha256 string) {
	found, _ := dbFind(db, sha256)
	if !found {
		dbInsert(db, lastPath, size, mime, md5, sha1, sha256)
	}
}

func dbUpdatePath(db *sql.DB, sha256sum string, newPath string) {
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

func dbFindPHash(db *sql.DB, sha256 string) (found bool) {
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

func dbUpdatePHash(db *sql.DB, sha256sum string, pHash string) {
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

func dbFilesWithPHash(db *sql.DB, page int64, sha256sum string) (distances map[int]Distance) {
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
