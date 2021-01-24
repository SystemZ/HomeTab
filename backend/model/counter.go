package model

import (
	"github.com/go-sql-driver/mysql"
	"log"
	"sort"
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
	Precise   uint       `gorm:"column:precise" json:"precise"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

type CounterTag struct {
	Id        uint       `json:"id" gorm:"primary_key;type:uint(10)" json:"id"`
	CounterId uint       `gorm:"column:counter_id" json:"counter_id"`
	Name      string     `gorm:"column:name" json:"name"`
	CreatedAt *time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt *time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func CreateCounter(name string) uint {
	var counter Counter

	counter.Name = name
	//FIXME timezones
	now := time.Now()
	counter.ProjectId = 1
	counter.CreatedAt = &now
	counter.UpdatedAt = &now

	err := DB.Save(&counter).Error
	if err != nil {
		log.Printf("%v", err)
	}

	return counter.Id
}

func StartCounterSession(counterId uint, userId uint) uint {
	var session CounterSession
	//FIXME timezones
	now := time.Now()
	session.CounterId = counterId
	session.UserId = userId
	session.Precise = 1
	session.StartedAt = &now
	session.CreatedAt = &now
	session.UpdatedAt = &now
	DB.Save(&session)
	return session.Id
}

func StopCounterSession(counterId uint, userId uint) (sessionId uint, sessionTaken string) {
	var session CounterSession
	res := DB.Order("ended_at asc").Where("user_id = ? AND counter_id = ? AND ended_at IS NULL", userId, counterId).First(&session)
	if res.RowsAffected < 1 {
		return 0, ""
	}
	//FIXME timezones
	now := time.Now()
	session.EndedAt = &now
	session.UpdatedAt = &now
	DB.Save(&session)

	sessionId = session.Id

	started := time.Unix(session.StartedAt.Unix(), 0)
	ended := time.Unix(session.EndedAt.Unix(), 0)
	sessionTaken = PrettyTime(uint(ended.Sub(started).Seconds()))

	return
}

type CounterList struct {
	Counter
	Tags                string `json:"tags"`
	Seconds7d           uint   `json:"secondsD7"`
	Seconds30d          uint   `json:"secondsD30"`
	SecondsAll          uint   `json:"secondsAll"`
	Seconds7dFormatted  string `json:"secondsD7F"`
	Seconds30dFormatted string `json:"secondsD30F"`
	SecondsAllFormatted string `json:"secondsAllF"`
	Running             uint   `json:"running"`
}

func CountersLongList(userId uint) (result []CounterList) {
	query := `
SELECT
  counters.id,
  counters.name,
  (SELECT GROUP_CONCAT(counter_tags.name SEPARATOR ',') FROM counter_tags WHERE counter_tags.counter_id = counters.id) AS tags,
  counters.created_at,
  counters.updated_at,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.started_at > NOW() - INTERVAL 7 DAY
  ), 0) AS seconds_7d,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.started_at > NOW() - INTERVAL 30 DAY
  ), 0) AS seconds_30d,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
  ), 0) AS seconds_all,
  (
    SELECT COUNT(*)
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.ended_at IS NULL
  ) AS running
FROM counters
GROUP BY counters.id
ORDER BY counters.id DESC
`
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, userId, userId, userId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list CounterList
		err := rows.Scan(&list.Id, &list.Name, &list.Tags, &list.CreatedAt, &list.UpdatedAt, &list.Seconds7d, &list.Seconds30d, &list.SecondsAll, &list.Running)
		if err != nil {
			return
		}

		list.Seconds7dFormatted = PrettyTime(list.Seconds7d)
		list.Seconds30dFormatted = PrettyTime(list.Seconds30d)
		list.SecondsAllFormatted = PrettyTime(list.SecondsAll)
		result = append(result, list)
	}
	return result
}

func CounterStats(counterId, userId uint) (result []CounterList) {
	query := `
SELECT
  counters.id,
  counters.name,
  (SELECT GROUP_CONCAT(counter_tags.name SEPARATOR ',') FROM counter_tags WHERE counter_tags.counter_id = counters.id) AS tags,
  counters.created_at,
  counters.updated_at,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.started_at > NOW() - INTERVAL 7 DAY
  ), 0) AS seconds_7d,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.started_at > NOW() - INTERVAL 30 DAY
  ), 0) AS seconds_30d,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
  ), 0) AS seconds_all,
  (
    SELECT COUNT(*)
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.ended_at IS NULL
  ) AS running
FROM counters
WHERE counters.id = ?
GROUP BY counters.id
ORDER BY counters.id DESC
`
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, userId, userId, userId, counterId)
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list CounterList
		err := rows.Scan(&list.Id, &list.Name, &list.Tags, &list.CreatedAt, &list.UpdatedAt, &list.Seconds7d, &list.Seconds30d, &list.SecondsAll, &list.Running)
		if err != nil {
			return
		}

		list.Seconds7dFormatted = PrettyTime(list.Seconds7d)
		list.Seconds30dFormatted = PrettyTime(list.Seconds30d)
		list.SecondsAllFormatted = PrettyTime(list.SecondsAll)
		result = append(result, list)
	}
	return result
}

func CountersLongListPaginate(userId uint, limit int, nextId int, prevId int, qTerm string) (result []CounterList, allRecords int) {
	// count... counters
	if len(qTerm) < 1 {
		qTerm = "%"
	}
	scoutQuery := `SELECT COUNT(*) FROM counters WHERE counters.name LIKE ?`
	stmt1, err := DB.DB().Prepare(scoutQuery)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt1.Close()
	rows1, err := stmt1.Query(qTerm)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows1.Close()
	for rows1.Next() {
		err := rows1.Scan(&allRecords)
		if err != nil {
			return
		}
	}

	// no records visible? don't ask DB again
	if allRecords < 1 {
		return result, allRecords
	}

	// we have some results so let's get details
	whereSign := ">"
	sortType := "ASC"
	if nextId < prevId {
		nextId = prevId
		whereSign = "<"
		sortType = "DESC"
	}

	// get counters
	query := `
SELECT
  counters.id,
  counters.name,
  (SELECT GROUP_CONCAT(counter_tags.name SEPARATOR ',') FROM counter_tags WHERE counter_tags.counter_id = counters.id) AS tags,
  counters.created_at,
  counters.updated_at,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.started_at > NOW() - INTERVAL 7 DAY
  ), 0) AS seconds_7d,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.started_at > NOW() - INTERVAL 30 DAY
  ), 0) AS seconds_30d,
  IFNULL((
    SELECT SUM(TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())))
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
  ), 0) AS seconds_all,
  (
    SELECT COUNT(*)
    FROM counter_sessions
    WHERE
      counters.id = counter_sessions.counter_id
    AND
      counter_sessions.user_id = ?
    AND
      counter_sessions.ended_at IS NULL
  ) AS running
FROM counters
WHERE id ` + whereSign + ` ?
AND counters.name LIKE ?
GROUP BY counters.id
ORDER BY counters.id ` + sortType + `
LIMIT ?
`
	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, userId, userId, userId, nextId, qTerm, limit)
	if err != nil {
		log.Printf("%v", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list CounterList
		err := rows.Scan(&list.Id, &list.Name, &list.Tags, &list.CreatedAt, &list.UpdatedAt, &list.Seconds7d, &list.Seconds30d, &list.SecondsAll, &list.Running)
		if err != nil {
			return
		}

		list.Seconds7dFormatted = PrettyTime(list.Seconds7d)
		list.Seconds30dFormatted = PrettyTime(list.Seconds30d)
		list.SecondsAllFormatted = PrettyTime(list.SecondsAll)
		result = append(result, list)
	}

	sort.Slice(result, func(p, q int) bool {
		return result[p].Id < result[q].Id
	})
	return result, allRecords
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

type CounterSessionList struct {
	CounterId         uint
	Id                uint
	UserId            uint
	Name              string
	Tags              string
	StartedAt         time.Time
	EndedAt           mysql.NullTime
	Duration          uint
	DurationFormatted string
	Running           bool
}

func CounterLogList(userId uint, limit int) (result []CounterSessionList) {
	query := `
SELECT 
  counter_sessions.counter_id,
  counter_sessions.id,
  counter_sessions.user_id,
  counters.name, 
  (SELECT GROUP_CONCAT(counter_tags.name SEPARATOR ',') FROM counter_tags WHERE counter_tags.counter_id = counters.id) AS tags,
  counter_sessions.started_at, 
  counter_sessions.ended_at,
  TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())) AS duration
FROM counter_sessions
JOIN counters ON counters.id = counter_sessions.counter_id
WHERE counter_sessions.deleted_at IS NULL
  AND user_id = ?
ORDER BY counter_sessions.started_at DESC
LIMIT ?
`

	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, limit)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list CounterSessionList
		err := rows.Scan(&list.CounterId, &list.Id, &list.UserId, &list.Name, &list.Tags, &list.StartedAt, &list.EndedAt, &list.Duration)
		if err != nil {
			log.Printf("%v", err.Error())
			return
		}
		list.DurationFormatted = PrettyTime(list.Duration)
		list.Running = !list.EndedAt.Valid
		result = append(result, list)
	}
	return result
}

func CounterLog(counterId int, userId uint, limit int) (result []CounterSessionList) {
	query := `
SELECT 
  counter_sessions.counter_id,
  counter_sessions.id,
  counter_sessions.user_id,
  counters.name, 
  (SELECT GROUP_CONCAT(counter_tags.name SEPARATOR ',') FROM counter_tags WHERE counter_tags.counter_id = counters.id) AS tags,
  counter_sessions.started_at, 
  counter_sessions.ended_at,
  TIMESTAMPDIFF(SECOND, counter_sessions.started_at,IFNULL(counter_sessions.ended_at, NOW())) AS duration
FROM counter_sessions
JOIN counters ON counters.id = counter_sessions.counter_id
WHERE counter_sessions.deleted_at IS NULL
  AND user_id = ?
  AND counters.id = ?
ORDER BY counter_sessions.started_at DESC
LIMIT ?
`

	stmt, err := DB.DB().Prepare(query)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(userId, counterId, limit)
	if err != nil {
		log.Printf("%v", err.Error())
		return
	}
	defer rows.Close()
	for rows.Next() {
		var list CounterSessionList
		err := rows.Scan(&list.CounterId, &list.Id, &list.UserId, &list.Name, &list.Tags, &list.StartedAt, &list.EndedAt, &list.Duration)
		if err != nil {
			log.Printf("%v", err.Error())
			return
		}
		list.DurationFormatted = PrettyTime(list.Duration)
		list.Running = !list.EndedAt.Valid
		result = append(result, list)
	}
	return result
}
