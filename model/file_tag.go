package model

import (
	"database/sql"
	"log"
	"time"
)

type FileTag struct {
	Id     int `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	UserId int `gorm:"column:user_id"`
	FileId int `gorm:"column:file_id"`
	TagId  int `gorm:"column:tag_id"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type FileTagList struct {
	FileId int    `json:"fileId"`
	Tags   string `json:"tags"`
}

// TODO optimize, don't run it in loop
func FileTagsList(fileList []FileTagList, userId int) []FileTagList {
	var finalList []FileTagList
	for _, file := range fileList {
		query := `
	   SELECT 
	     GROUP_CONCAT(tags.tag SEPARATOR ',') AS tagz
       FROM tags
       INNER JOIN file_tags ON tags.id = file_tags.tag_id
       WHERE file_tags.file_id = ?
       AND file_tags.deleted_at IS NULL
	   AND file_tags.user_id = ?
	   `
		stmt, err := DB.DB().Prepare(query)
		if err != nil {
			log.Printf("%v", err)
			return fileList
		}
		defer stmt.Close()
		rows, err := stmt.Query(file.FileId, userId)
		if err != nil {
			log.Printf("%v", err)
			return fileList
		}
		defer rows.Close()
		for rows.Next() {
			var tagsRaw sql.NullString
			err := rows.Scan(&tagsRaw)
			if err != nil {
				log.Printf("sql scan error: %v", err)
				return fileList
			}
			// append only if we have some tags
			// we don't need to send unnecessary data
			if !tagsRaw.Valid {
				continue
			}
			finalList = append(finalList, FileTagList{
				FileId: file.FileId,
				Tags:   tagsRaw.String,
			})
		}
	}
	return finalList
}
