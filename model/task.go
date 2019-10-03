package model

import (
	"github.com/jinzhu/gorm"
	"log"
	"time"
)

type Task struct {
	gorm.Model
	Id             uint       `gorm:"primary_key;type:uint(10)" json:"id"`
	Subject        string     `gorm:"column:subject" json:"subject"`
	ProjectId      uint       `gorm:"column:project_id"`
	AssignedUserId uint       `gorm:"column:assigned_user_id"`
	Repeating      uint       `gorm:"column:repeating"`
	NeverEnding    uint       `gorm:"column:never_ending"`
	RepeatUnit     string     `gorm:"column:repeat_unit"`
	RepeatMin      uint       `gorm:"column:repeat_min"`
	RepeatBest     uint       `gorm:"column:repeat_best"`
	RepeatMax      uint       `gorm:"column:repeat_max"`
	RepeatFrom     *time.Time `gorm:"column:repeat_from"`
	// RepeatFrom       mysql.NullTime `gorm:"column:repeat_from"`
	EstimateS        uint `gorm:"column:estimate_s"`
	MasterTaskId     uint `gorm:"column:master_task_id"`
	SeparateChildren uint `gorm:"column:separate_children"`
}

func CreateTask(task Task) {
	now := time.Now()
	if task.RepeatFrom == nil {
		task.RepeatFrom = &now
	}

	err := DB.Save(&task).Error
	if err != nil {
		log.Printf("%v", err)
	}
}
