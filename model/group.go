package model

import (
	"log"
	"time"
)

type Group struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	Name      string     `gorm:"column:name" json:"username"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
}

func CreateGroup(name string) uint {
	var group Group

	group.Name = name
	//FIXME timezones
	now := time.Now()
	group.CreatedAt = &now
	group.UpdatedAt = &now

	save := DB.Save(&group)

	if save.Error != nil {
		log.Printf("%v", save.Error.Error())
	}

	return group.Id
}
