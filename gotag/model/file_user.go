package model

import "time"

type FileUser struct {
	Id     int `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	FileId int `gorm:"column:file_id"`
	UserId int `gorm:"column:user_id"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
