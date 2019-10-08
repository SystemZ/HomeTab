package model

import (
	"database/sql"
	"time"
)

type Event struct {
	Id        uint           `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	UserId    uint           `gorm:"column:user_id"`
	DeviceId  uint           `gorm:"column:device_id"`
	Code      EventCode      `gorm:"column:code" json:"code"`
	ValueStr  sql.NullString `gorm:"column:val_str" json:"val_str"`
	ValueInt  sql.NullInt64  `gorm:"column:val_int"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"created_at"`
}

type EventCode uint

const (
	Noop EventCode = iota
	DeviceBatteryPercent
	DeviceScreenOn
	DeviceMusicTrack
	DeviceMusicArtist
	DeviceScreenOff
	//DeviceBoot
	//DeviceShutdown
	//DeviceGoneSleep
	//DeviceWokeUp
	//DeviceHSConnected
	//DeviceHSDisconnected
	//DeviceBatteryCharged
	//DeviceBatteryFull
	//DeviceRinging
	//DeviceOffhook
	// user health
	//UserGoneSleep
	//UserRemPhase
	//UserWokeUp
	// user location
	//UserEnteredHome
	//UserLeftHome
	// user music
	//UserStartedMusic
	//UserStoppedMusic
	// user browser
	//UserBrowserStarted
	//UserBrowserClosed
	//UserBrowserUrl
	// by ping
	//printer online
	//printer offline
)

func DeviceEventAdd(code EventCode, userId uint, deviceId uint) {
	event := Event{
		UserId:   userId,
		DeviceId: deviceId,
		Code:     code,
	}
	DB.Create(&event)
}

func DeviceEventAddInt(code EventCode, userId uint, deviceId uint, val int) {
	event := Event{
		UserId:   userId,
		DeviceId: deviceId,
		Code:     code,
		ValueInt: struct {
			Int64 int64
			Valid bool
		}{Int64: int64(val), Valid: true},
	}
	DB.Create(&event)
}

func DeviceEventAddStr(code EventCode, userId uint, deviceId uint, str string) {
	event := Event{
		UserId:   userId,
		DeviceId: deviceId,
		Code:     code,
		ValueStr: struct {
			String string
			Valid  bool
		}{String: str, Valid: true},
	}
	DB.Create(&event)
}
