package model

import (
	"github.com/jinzhu/gorm"
	"log"
)

// Host represents the host for this application
// swagger:model user
type Project struct {
	gorm.Model
	// ID
	//
	// required: true
	Id uint `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`

	// Username
	//
	// required: true
	Name    string `gorm:"column:name" json:"username"`
	GroupId uint   `gorm:"column:group_id" json:"group_id"`
}

func CreateProject(name string, groupId uint) uint {
	var project Project

	project.Name = name
	project.GroupId = groupId
	//FIXME timezones

	err := DB.Save(&project).Error
	if err != nil {
		log.Printf("%v", err)
	}

	return project.Id
}
