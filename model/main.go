package model

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
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
	//migrator, _ := gomigrate.NewMigrator(db, gomigrate.Sqlite3{}, "./migrations")
	//err = migrator.Migrate()

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
