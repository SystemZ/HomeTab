package model

import (
	"log"
	"time"
)

type Counter struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	Name      string     `gorm:"column:name" json:"username"`
	ProjectId uint       `gorm:"column:project_id" json:"group_id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func CreateCounter(name string) uint {
	var counter Counter

	counter.Name = name
	//FIXME timezones
	now := time.Now()
	counter.CreatedAt = &now
	counter.UpdatedAt = &now

	err := DB.Save(&counter).Error
	if err != nil {
		log.Printf("%v", err)
	}

	return counter.Id
}
