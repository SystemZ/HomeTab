package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Log struct {
	Id        int        `json:"id" gorm:"AUTO_INCREMENT" json:"id"`
	Body      string     `gorm:"column:body"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
}

func AddLog(db *gorm.DB, body string) {
	log.Printf("%v", body)
	log := Log{Body: body}
	db.Create(&log)
}
