package model2

import "time"

type File struct {
	Id       int    `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	UserId   int    `gorm:"column:user_id"`
	Filename string `gorm:"column:file_name"`
	FilePath string `gorm:"column:file_path;type:varchar(4096)"`
	SizeB    int    `gorm:"column:size_b"`
	MimeId   int    `gorm:"column:mime_id"`
	PhashA   int    `gorm:"column:phash_a;type:bigint(16)"`
	PhashB   int    `gorm:"column:phash_b;type:bigint(16)"`
	PhashC   int    `gorm:"column:phash_c;type:bigint(16)"`
	PhashD   int    `gorm:"column:phash_d;type:bigint(16)"`
	Sha256   string `gorm:"column:sha256;type:char(64)"`

	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
