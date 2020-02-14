package model2

import (
	"database/sql"
	"log"
	"sort"
	"time"
)

type File struct {
	Id       int    `gorm:"AUTO_INCREMENT" json:"id"`
	UserId   int    `gorm:"column:user_id" json:"userId"`
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
	Mime string `gorm:"-" json:"mime"`
	Tags string `gorm:"-" json:"tags"`
}

func FileListPaginate(userId int, limit int, nextId int, prevId int, qTerm string) (result []File, allRecords int) {
	// count... counters
	if len(qTerm) < 1 {
		qTerm = "%"
	}
	scoutQuery := `SELECT COUNT(*) FROM files WHERE file_name LIKE ?`
	stmt1, err := DB.DB().Prepare(scoutQuery)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt1.Close()
	rows1, err := stmt1.Query(qTerm)
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

	// get counters
	query := `
SELECT
  files.id,
  files.sha256,
  files.file_name,
  files.file_path,
  files.size_b,
  files.created_at,
  files.updated_at,
  mimes.mime,
  (SELECT GROUP_CONCAT(tags.tag SEPARATOR ',')
   FROM tags
   INNER JOIN file_tags on tags.id = file_tags.tag_id
   WHERE file_tags.file_id = files.id
   AND file_tags.deleted_at IS NULL
  ) AS tagz
FROM files
INNER JOIN mimes on files.mime_id = mimes.id
WHERE files.user_id = ? AND files.id ` + whereSign + ` ?
AND files.file_name LIKE ?
GROUP BY files.id
ORDER BY files.id ` + sortType + `
LIMIT ?
`
	/* this doesn't show not tagged files :(
	SELECT
	  files.id,
	  files.sha256,
	  files.file_name,
	  files.file_path,
	  files.size_b,
	  files.created_at,
	  files.updated_at,
	  mimes.mime,
	  GROUP_CONCAT(tags.tag SEPARATOR ',') as tags
	FROM files
	INNER JOIN mimes on files.mime_id = mimes.id
	INNER JOIN file_tags on file_tags.file_id = files.id
	INNER JOIN tags on file_tags.tag_id = tags.id
	WHERE files.user_id = ? AND files.id ` + whereSign + ` ?
	AND files.file_name LIKE ?
	GROUP BY files.id
	ORDER BY files.id ` + sortType + `
	LIMIT ?
	*/

	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, nextId, qTerm, limit)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list File
		var tagz sql.NullString
		err := rows.Scan(&list.Id, &list.Sha256, &list.Filename, &list.FilePath, &list.SizeB, &list.CreatedAt, &list.UpdatedAt, &list.Mime, &tagz)
		if tagz.Valid {
			list.Tags = tagz.String
		}
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
