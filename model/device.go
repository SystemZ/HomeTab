package model

import (
	uuid "github.com/satori/go.uuid"
	"log"
	"time"
)

type Device struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	UserId    uint       `gorm:"column:user_id"`
	Name      string     `gorm:"column:name"`
	Token     string     `gorm:"token"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func CreateDevice(name string, userId uint) uint {
	var device Device

	device.Name = name
	device.UserId = userId
	device.Token = uuid.NewV4().String()

	err := DB.Save(&device).Error
	if err != nil {
		log.Printf("%v", err)
	}

	return device.Id
}

func GetListOfDevices() {
	query := `
SELECT
id,
user_id,
name,
created_at,
(SELECT events.val_int FROM events WHERE events.device_id = devices.id ORDER BY events.created_at ASC LIMIT 1) AS battery_left
FROM devices
`
	log.Printf("%v", query)
}
