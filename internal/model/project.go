package model

import (
	"log"
	"time"
)

// Host represents the host for this application
// swagger:model user
type Project struct {
	// ID
	//
	// required: true
	Id uint `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`

	// Username
	//
	// required: true
	Name      string     `gorm:"column:name" json:"name"`
	GroupId   uint       `gorm:"column:group_id" json:"groupId"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updatedAt"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deletedAt"`
}

func CreateProject(name string, groupId uint) uint {
	var project Project

	project.Name = name
	project.GroupId = groupId
	//FIXME timezones
	now := time.Now()
	project.CreatedAt = &now
	project.UpdatedAt = &now

	err := DB.Save(&project).Error
	if err != nil {
		log.Printf("%v", err)
	}

	return project.Id
}
