package model

import (
	"log"
	"strconv"
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
	res := DB.Order("ended_at asc").Where("user_id = ? AND counter_id = ? AND ended_at IS NULL", userId, counterId).First(&session)
	if res.RowsAffected < 1 {
		return
	}
	//FIXME timezones
	now := time.Now()
	session.EndedAt = &now
	session.UpdatedAt = &now
	DB.Save(&session)
}

type CounterList struct {
	Counter
	Seconds       uint
	TimeFormatted string
}

func CountersLongList() (result []CounterList) {
	query := `
	   SELECT
	   counters.id,
	   counters.name,
	   IFNULL(
	     IFNULL(
	         SUM(TIMESTAMPDIFF(SECOND,counter_sessions.started_at, counter_sessions.ended_at)),
	         SUM(TIMESTAMPDIFF(SECOND,counter_sessions.started_at, NOW()))
	      ),
	      0
	   ) AS seconds
	   FROM counters
	   LEFT JOIN counter_sessions
	   ON counters.id = counter_sessions.counter_id
	   GROUP BY counters.id`
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query() //
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list CounterList
		err := rows.Scan(&list.Id, &list.Name, &list.Seconds)
		if err != nil {
			return
		}

		list.TimeFormatted = PrettyTime(list.Seconds)
		result = append(result, list)
	}
	return result
}

func PrettyTime(s uint) string {
	var h int
	var m int
	for s >= 3600 {
		s -= 3600
		h++
	}
	for s >= 60 {
		s -= 60
		m++
	}
	return strconv.Itoa(h) + "h " + strconv.Itoa(m) + "m " + strconv.Itoa(int(s)) + "s"
}
