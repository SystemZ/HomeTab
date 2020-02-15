package model

import (
	"database/sql"
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
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
