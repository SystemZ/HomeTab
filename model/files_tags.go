package model

import (
	"database/sql"
	"path/filepath"
)

func FileTagSearchByName(db *sql.DB, page int64, tag string) (found bool, sha256s map[int]File) {
	found = false

	limit := int64(15)
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	rows, err := db.Query("SELECT files.id, files.last_path, files.size, files.sha256 FROM tags JOIN files ON files.id = tags.fid WHERE tags.name = ? ORDER BY files.id LIMIT ?, ?", tag, offset, limit)
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
		sha256s[id] = File{Fid: id, LastPath: filename, Size: size, Sha256: sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func FilesWithoutTags(db *sql.DB, page int64) (found bool, sha256s map[int]File) {
	found = false

	limit := int64(15)
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	rows, err := db.Query("SELECT f.id, f.last_path, f.size, f.sha256, (SELECT COUNT(id) FROM tags t WHERE t.fid = f.id) AS tags FROM files f WHERE tags < 1 ORDER BY id LIMIT ?, ?", offset, limit)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var last_path string
		var size int
		var sha256 string
		var tags int
		err = rows.Scan(&id, &last_path, &size, &sha256, &tags)

		filename := filepath.Base(last_path)

		checkErr(err)

		if sha256s == nil {
			sha256s = make(map[int]File)
		}
		sha256s[id] = File{Fid: id, LastPath: filename, Size: size, Sha256: sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}

func FilesWithoutTagsRandom(db *sql.DB, page int64) (found bool, sha256s map[int]File) {
	found = false

	limit := int64(15)
	offset := (page - 1) * limit
	if offset <= 0 {
		offset = 0
	}

	rows, err := db.Query("SELECT f.id, f.last_path, f.size, f.sha256, (SELECT COUNT(id) FROM tags t WHERE t.fid = f.id) AS tags FROM files f WHERE tags < 1 ORDER BY random() LIMIT ?, ?", offset, limit)
	checkErr(err)
	defer rows.Close()

	for rows.Next() {
		var id int
		var last_path string
		var size int
		var sha256 string
		var tags int
		err = rows.Scan(&id, &last_path, &size, &sha256, &tags)

		filename := filepath.Base(last_path)

		checkErr(err)

		if sha256s == nil {
			sha256s = make(map[int]File)
		}
		sha256s[id] = File{Fid: id, LastPath: filename, Size: size, Sha256: sha256}

		if !found {
			found = true
		}
	}
	return found, sha256s
}
