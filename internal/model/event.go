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
	ValueStr  sql.NullString `gorm:"column:val_str" json:"valStr"`
	ValueInt  sql.NullInt64  `gorm:"column:val_int"`
	CreatedAt *time.Time     `gorm:"column:created_at" json:"createdAt"`
}

type EventCode uint

const (
	Noop EventCode = iota
	DeviceBatteryPercent
	DeviceScreenOn
	DeviceMusicTrack
	DeviceMusicArtist
	DeviceScreenOff
	TaskDone
	//DeviceUnlock
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

func TaskDoneEvent(userId uint, taskId int) {
	event := Event{
		UserId: userId,
		Code:   TaskDone,
		ValueInt: sql.NullInt64{
			Int64: int64(taskId),
			Valid: true,
		},
	}
	DB.Create(&event)
}

type TaskDoneT struct {
	Subject   string     `json:"subject"`
	Username  string     `json:"username"`
	CreatedAt *time.Time `json:"createdAt"`
}

func TaskDoneEvents7days() (result []TaskDoneT) {
	query := `
SELECT 
  tasks.subject,
  users.username,
  events.created_at
FROM events
JOIN users ON events.user_id = users.id
JOIN tasks ON events.val_int = tasks.id
WHERE events.code = ?
AND (events.created_at >= DATE_SUB(CURDATE(), INTERVAL 7 DAY))
ORDER BY events.created_at DESC
`

	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(TaskDone)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var list TaskDoneT
		err := rows.Scan(&list.Subject, &list.Username, &list.CreatedAt)
		if err != nil {
			return
		}
		result = append(result, list)
	}
	return result
}
