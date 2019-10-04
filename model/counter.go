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

type CounterSession struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	CounterId uint       `gorm:"column:counter_id" json:"counter_id"`
	UserId    uint       `gorm:"column:user_id" json:"user_id"`
	StartedAt *time.Time `gorm:"column:started_at" json:"started_at"`
	EndedAt   *time.Time `gorm:"column:ended_at" json:"ended_at"`
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

func StartCounterSession(counterId uint, userId uint) {
	var session CounterSession
	//FIXME timezones
	now := time.Now()
	session.CounterId = counterId
	session.UserId = userId
	session.StartedAt = &now
	session.CreatedAt = &now
	session.UpdatedAt = &now
	DB.Save(&session)
}

func StopCounterSession(counterId uint, userId uint) {
	var session CounterSession
	DB.Where(&CounterSession{UserId: userId, CounterId: counterId}).First(&session)
	//FIXME timezones
	now := time.Now()
	session.EndedAt = &now
	session.UpdatedAt = &now
	DB.Save(&session)
}

// SELECT TIMESTAMPDIFF(SECOND,started_at, ended_at) AS time_diff FROM `counter_sessions`
// SELECT SUM(TIMESTAMPDIFF(SECOND,started_at, ended_at)) AS time_diff FROM `counter_sessions`
