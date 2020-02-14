package model2

import (
	"github.com/jinzhu/gorm"
	"time"
)

type Mime struct {
	Id   int    `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	Mime string `gorm:"column:mime"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func AddMime(db *gorm.DB, mime string) int {
	// get tag if exists already
	// FIXME use transactions
	var mimeInDb Mime
	db.Where("mime = ?", mime).First(&mimeInDb)

	// tag is not present in DB, create new
	if mimeInDb.Id < 1 {
		mimeInDb = Mime{Mime: mime}
		db.Create(&mimeInDb)
	}

	// mime is now present in DB
	//log.Printf("Mime in DB: %+v", mimeInDb)
	return mimeInDb.Id
}
