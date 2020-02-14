package model2

import (
	"log"
	"sort"
	"time"
)

type File struct {
	Id       int    `gorm:"AUTO_INCREMENT" json:"id"`
	UserId   int    `gorm:"column:user_id"`
	Filename string `gorm:"column:file_name" json:"filename"`
	FilePath string `gorm:"column:file_path;type:varchar(4096)"`
	SizeB    int    `gorm:"column:size_b"`
	MimeId   int    `gorm:"column:mime_id"`
	PhashA   int    `gorm:"column:phash_a;type:bigint(16)" json:"pHashA"`
	PhashB   int    `gorm:"column:phash_b;type:bigint(16)"`
	PhashC   int    `gorm:"column:phash_c;type:bigint(16)"`
	PhashD   int    `gorm:"column:phash_d;type:bigint(16)"`
	Sha256   string `gorm:"column:sha256;type:char(64)" json:"sha256"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	Thumb int `gorm:"-" json:"thumb"`
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
  files.phash_a,
  files.created_at,
  files.updated_at,
  mimes.mime
FROM files
INNER JOIN mimes on files.mime_id = mimes.id
WHERE files.user_id = ? AND files.id ` + whereSign + ` ?
AND files.file_name LIKE ?
GROUP BY files.id
ORDER BY files.id ` + sortType + `
LIMIT ?
`
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
		var mime string
		err := rows.Scan(&list.Id, &list.Sha256, &list.Filename, &list.PhashA, &list.CreatedAt, &list.UpdatedAt, &mime)
		if err != nil {
			return
		}

		if mime == "image/jpeg" ||
			mime == "image/png" ||
			mime == "image/gif" ||
			mime == "video/webm" ||
			mime == "video/mp4" {
			list.Thumb = 1
		}
		result = append(result, list)
	}

	sort.Slice(result, func(p, q int) bool {
		return result[p].Id < result[q].Id
	})
	return result, allRecords
}
