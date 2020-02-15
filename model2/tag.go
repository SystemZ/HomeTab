package model2

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Tag struct {
	Id  int    `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	Tag string `gorm:"column:tag" json:"tag"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`

	// not in table
	Counter int `gorm:"-" json:"counter"`
}

func AddTagToFile(db *gorm.DB, tag string, fileId int) {
	// get tag if exists already
	// FIXME use transactions
	var tagInDb Tag
	db.Where("tag = ?", tag).First(&tagInDb)

	// tag is not present in DB, create new
	if tagInDb.Id < 1 {
		tagInDb = Tag{Tag: tag}
		db.Create(&tagInDb)
	}

	// tag is now present in DB
	//log.Printf("Tag in DB: %+v", tagInDb)

	// check if link file <-> tag is present too
	var fileTagInDb FileTag
	db.Where("tag_id = ? AND file_id = ?", tagInDb.Id, fileId).First(&fileTagInDb)

	// link is not in DB, create new
	//log.Printf("%v", fileTagInDb)
	if fileTagInDb.Id < 1 {
		fileTagInDb = FileTag{
			FileId: fileId,
			TagId:  tagInDb.Id,
		}
		db.Create(&fileTagInDb)
	}

}

func TagList() (result []Tag) {
	query := `
SELECT COUNT(tags.id) AS counter, tags.id, tags.tag, tags.created_at, tags.updated_at
FROM tags
LEFT JOIN file_tags ON tags.id = file_tags.tag_id
WHERE file_tags.deleted_at IS NULL
GROUP BY tags.id
ORDER BY COUNT(tags.id)
DESC 
`
	//LIMIT 100
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list Tag
		err := rows.Scan(&list.Counter, &list.Id, &list.Tag, &list.CreatedAt, &list.UpdatedAt)
		if err != nil {
			log.Printf("sql scan error: %v", err)
			return
		}
		result = append(result, list)
	}
	return result
}
