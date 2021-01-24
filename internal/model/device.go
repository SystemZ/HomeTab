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
	TokenPush string     `gorm:"token_push"`
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

type DeviceList struct {
	Id                 uint      `json:"id"`
	UserId             uint      `json:"userId"`
	Name               string    `json:"name"`
	CreatedAt          time.Time `json:"createdAt"`
	DisplayState       string    `json:"displayState"`
	DisplayOnLastTime  time.Time `json:"displayLastOn"`
	DisplayOffLastTime time.Time `json:"displayLastOff"`
	BatteryLeft        uint      `json:"battery"`
	Username           string    `json:"username"`
	MusicTrack         string    `json:"musicTrack"`
	MusicArtist        string    `json:"musicArtist"`
	MusicLastPlayed    time.Time `json:"musicLastPlayed"`
}

func GetListOfDevices() (result []DeviceList) {
	query := `
SELECT
  id,
  user_id,
  name,
  created_at,
  (SELECT
       IF(events.code=?,'ON','OFF')
   FROM events
   WHERE events.device_id = devices.id
     AND events.code IN (?,?)
     AND events.val_int IS NULL
   ORDER BY events.created_at DESC
   LIMIT 1) AS display_current,
  (SELECT
    events.created_at
   FROM events
   WHERE events.device_id = devices.id
   AND events.code = ?
   AND events.val_int IS NULL
   ORDER BY events.created_at DESC
   LIMIT 1) AS last_display_on,
  (SELECT
    events.created_at
   FROM events
   WHERE events.device_id = devices.id
   AND events.code = ?
   AND events.val_int IS NULL
   ORDER BY events.created_at DESC
   LIMIT 1) AS last_display_off,
  (SELECT
    events.val_int
   FROM events
   WHERE events.device_id = devices.id
     AND events.code = ?
     AND events.val_int IS NOT NULL
   ORDER BY events.created_at DESC
   LIMIT 1) AS battery_left,
  (SELECT
    users.username
    FROM users
    WHERE users.id = devices.user_id
  ) AS username,
  (SELECT
    events.val_str
    FROM events
    WHERE events.device_id = devices.id
      AND events.code = ?
      AND events.val_int IS NULL
    ORDER BY events.created_at DESC
    LIMIT 1) AS music_title,
  (SELECT
    events.val_str
    FROM events
    WHERE events.device_id = devices.id
      AND events.code = ?
      AND events.val_int IS NULL
    ORDER BY events.created_at DESC
    LIMIT 1) AS music_artist,
  (SELECT
    events.created_at
    FROM events
    WHERE events.device_id = devices.id
      AND (events.code = ? OR events.code = ?)
      AND events.val_int IS NULL
    ORDER BY events.created_at DESC
    LIMIT 1) AS music_last_played
  FROM devices
`

	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(
		//1q
		DeviceScreenOn,
		DeviceScreenOn,
		DeviceScreenOff,
		//2q
		DeviceScreenOn,
		//3q
		DeviceScreenOff,
		//4q
		DeviceBatteryPercent,
		//5q
		DeviceMusicTrack,
		//6q
		DeviceMusicArtist,
		//7q
		DeviceMusicTrack,
		DeviceMusicArtist,
	)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list DeviceList
		err := rows.Scan(
			&list.Id,
			&list.UserId,
			&list.Name,
			&list.CreatedAt,
			&list.DisplayState,
			&list.DisplayOnLastTime,
			&list.DisplayOffLastTime,
			&list.BatteryLeft,
			&list.Username,
			&list.MusicTrack,
			&list.MusicArtist,
			&list.MusicLastPlayed,
		)
		if err != nil {
			return
		}

		result = append(result, list)
	}
	return result
}
