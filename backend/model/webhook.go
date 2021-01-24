package model

import (
	"time"
)

type WebhookReceiver struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	ProjectId uint       `gorm:"column:project_id"`
	ActionId  uint       `gorm:"column:action_id"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type WebhookAction struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	Extra     string     `gorm:"column:extra"`
	Type      uint       `gorm:"column:type"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}
