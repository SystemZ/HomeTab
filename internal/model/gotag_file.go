package model

import (
	"database/sql"
	"log"
	"sort"
	"time"
)

type File struct {
	Id       int    `gorm:"AUTO_INCREMENT" json:"id"`
	Filename string `gorm:"column:file_name" json:"filename"`
	FilePath string `gorm:"column:file_path;type:varchar(4096)" json:"filePath"`
	SizeB    int    `gorm:"column:size_b" json:"sizeB"`
	MimeId   int    `gorm:"column:mime_id" json:"-"`
	PhashA   int    `gorm:"column:phash_a;type:bigint(16)" json:"-"`
	PhashB   int    `gorm:"column:phash_b;type:bigint(16)" json:"-"`
	PhashC   int    `gorm:"column:phash_c;type:bigint(16)" json:"-"`
	PhashD   int    `gorm:"column:phash_d;type:bigint(16)" json:"-"`
	Sha256   string `gorm:"column:sha256;type:char(64)" json:"sha256"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"-"`

	// helpers, not present in DB
	Mime     string `gorm:"-" json:"mime"`
	Tags     string `gorm:"-" json:"tags"`
	Distance int    `gorm:"-" json:"distance"`
}

func GetFileByIdForUser(fileId int, userId int) (res File) {
	query := `
     SELECT
      files.id,
      files.sha256,
      files.file_name,
      files.file_path,
      files.size_b,
      files.created_at,
      files.updated_at,
      mimes.mime
    FROM files
    INNER JOIN mimes ON files.mime_id = mimes.id
    INNER JOIN file_users ON file_users.file_id = files.id
    WHERE file_users.user_id = ?
    AND files.id = ?
    GROUP BY files.id
    LIMIT 1`

	err := DB.DB().
		QueryRow(query, userId, fileId).
		Scan(&res.Id, &res.Sha256, &res.Filename, &res.FilePath, &res.SizeB, &res.CreatedAt, &res.UpdatedAt, &res.Mime)
	if err != nil {
		log.Printf("sql scan error: %v", err)
		return
	}
	return
}

func FileListPaginate(userId int, limit int, nextId int, prevId int, qTerm string) (result []File, allRecords int) {
	var tagSearch bool
	if len(qTerm) < 1 {
		qTerm = "%"
	} else {
		tagSearch = true
	}

	filename := "%"
	// count how many results we will have for pagination
	// FIXME allow search by tag and filename at the same time
	scoutQuery := `
		SELECT COUNT(*)
		FROM files
		JOIN file_users ON files.id = file_users.file_id
		WHERE file_users.user_id = ?
		AND file_name LIKE ?`
	if tagSearch && qTerm == "none" {
		scoutQuery = `
			SELECT COUNT(*)
			FROM files
    		WHERE (SELECT COUNT(id) FROM file_tags WHERE file_tags.user_id = ? AND file_tags.file_id = files.id AND file_tags.deleted_at IS NULL) < 1
			AND files.file_name LIKE ?`
	} else if tagSearch { // file with tag X
		// FIXME ignore untagged
		scoutQuery = `
			SELECT COUNT(tags.id) AS counter
			FROM tags
	   			LEFT JOIN file_tags ON tags.id = file_tags.tag_id
	   			WHERE file_tags.user_id = ?
       			AND tags.tag = ?
	    	 	GROUP BY tags.id
			ORDER BY COUNT(tags.id) DESC`
	}
	stmt1, err := DB.DB().Prepare(scoutQuery)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt1.Close()
	var rows1 *sql.Rows

	if qTerm != "none" && qTerm != "%" {
		// tag search
		rows1, err = stmt1.Query(userId, qTerm)
	} else {
		// untagged and all files
		rows1, err = stmt1.Query(userId, filename)
	}
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		err := rows1.Scan(&allRecords)
		if err != nil {
			return
		}
	}

	// no records visible? don't ask DB again
	if allRecords < 1 {
		return result, allRecords
	}

	// we have some results so let's get details
	whereSign := ">"
	sortType := "ASC"
	if nextId < prevId {
		nextId = prevId
		whereSign = "<"
		sortType = "DESC"
	}

	// all files - standard search
	query := `
     SELECT
      files.id,
      files.sha256,
      files.file_name,
      files.file_path,
      files.size_b,
      files.created_at,
      files.updated_at,
      mimes.mime
    FROM files
    INNER JOIN mimes ON files.mime_id = mimes.id
    INNER JOIN file_users ON file_users.file_id = files.id
    WHERE file_users.user_id = ?
    AND files.id ` + whereSign + ` ?
    AND files.file_name LIKE ?
    GROUP BY files.id
    ORDER BY files.id ` + sortType + `
    LIMIT ?`

	if tagSearch && qTerm == "none" {
		query = `
    SELECT
      files.id,
      files.sha256,
      files.file_name,
      files.file_path,
      files.size_b,
      files.created_at,
      files.updated_at,
      mimes.mime
    FROM files
    INNER JOIN mimes on files.mime_id = mimes.id
    INNER JOIN file_users ON file_users.file_id = files.id
    WHERE (SELECT COUNT(id) FROM file_tags WHERE file_tags.file_id = files.id AND file_tags.deleted_at IS NULL) < 1
    AND file_users.user_id = ?
    AND files.id ` + whereSign + ` ?
	AND files.file_name LIKE ?
    ORDER BY files.id ` + sortType + `
	LIMIT ?
`
	} else if tagSearch {
		// show only files with tag X, scoped search
		// this doesn't show untagged files :(
		// subquery is nasty but without it we see only one tag
		query = `
    SELECT
	  files.id,
	  files.sha256,
	  files.file_name,
	  files.file_path,
	  files.size_b,
	  files.created_at,
	  files.updated_at,
	  mimes.mime
	FROM files
	INNER JOIN mimes on files.mime_id = mimes.id
	INNER JOIN file_tags on file_tags.file_id = files.id
	INNER JOIN tags on file_tags.tag_id = tags.id
    INNER JOIN file_users ON file_users.file_id = files.id
	WHERE file_users.user_id = ?
	AND files.id ` + whereSign + ` ?
	AND files.file_name LIKE ?
	AND tags.tag = ?
	AND file_tags.deleted_at IS NULL
	GROUP BY files.id
	ORDER BY files.id ` + sortType + `
	LIMIT ?`
	}

	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt.Close()
	var rows *sql.Rows
	// TODO make searchable by qTerm and tag simultaneously
	if tagSearch && qTerm == "none" {
		rows, err = stmt.Query(userId, nextId, "%", limit)
	} else if tagSearch {
		// in this case qTerm is tag, fix variable names
		rows, err = stmt.Query(userId, nextId, "%", qTerm, limit)
	} else {
		rows, err = stmt.Query(userId, nextId, qTerm, limit)
	}
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list File
		// TODO is this necessary?
		err := rows.Scan(&list.Id, &list.Sha256, &list.Filename, &list.FilePath, &list.SizeB, &list.CreatedAt, &list.UpdatedAt, &list.Mime)
		if err != nil {
			log.Printf("sql scan error: %v", err)
			return
		}
		result = append(result, list)
	}

	sort.Slice(result, func(p, q int) bool {
		return result[p].Id < result[q].Id
	})
	return result, allRecords
}

// TODO userId as a option, cross user suggestions
func SimilarFiles(sha256 string, userId int, limit int) (result []File) {
	var imgInDb File
	DB.Where("sha256 = ?", sha256).First(&imgInDb)

	// file with hash not found
	if imgInDb.Id < 1 {
		log.Printf("file with hash %v not found", sha256)
		return
	}

	query1 := `
SELECT HAMMINGDISTANCE(?,?,?,?,files.phash_a,files.phash_b,files.phash_c,files.phash_d) AS dist,
       files.id,
       sha256,
       file_name,
       file_path,
       size_b,
       mimes.mime
FROM files
INNER JOIN mimes ON files.mime_id = mimes.id
INNER JOIN file_users ON files.id = file_users.file_id
WHERE sha256 != ?
AND file_users.user_id = ?
AND files.phash_a != 0
AND files.phash_b != 0
AND files.phash_c != 0
AND files.phash_d != 0
ORDER BY dist ASC
LIMIT ?
`

	stmt1, err := DB.DB().Prepare(query1)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt1.Close()
	rows1, err := stmt1.Query(imgInDb.PhashA, imgInDb.PhashB, imgInDb.PhashC, imgInDb.PhashD, imgInDb.Sha256, userId, limit)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		var list File
		err := rows1.Scan(&list.Distance, &list.Id, &list.Sha256, &list.Filename, &list.FilePath, &list.SizeB, &list.Mime)
		if err != nil {
			log.Printf("sql scan error: %v", err)
			return
		}
		result = append(result, list)
	}

	return result
}
