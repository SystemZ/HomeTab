package model

import "time"

type FileTag struct {
	Id     int `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	FileId int `gorm:"column:file_id"`
	TagId  int `gorm:"column:tag_id"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
