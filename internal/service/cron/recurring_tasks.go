package cron

import (
	"database/sql"
	"github.com/jinzhu/gorm"
	"github.com/systemz/hometab/internal/model"
	"log"
	"time"
)

func ScanRecurring(db *gorm.DB) {
	// we can skip events table
	// but then we can't access who done
	// that task previous time
	query := `
SELECT
id,
(SELECT created_at FROM events WHERE val_int = tasks.id AND code = ? ORDER by created_at DESC LIMIT 1) AS done_recently
FROM tasks
WHERE repeating = 1
AND done_at IS NOT NULL
`
	stmt, err := db.DB().Prepare(query)
	if err != nil {
		return
	}
	defer stmt.Close()
	rows, err := stmt.Query(model.TaskDone)
	if err != nil {
		return
	}

	defer rows.Close()
	for rows.Next() {
		var task model.Task
		var taskId int
		var repeatDoneAt sql.NullTime
		// get scoped, custom info from DB
		err := rows.Scan(&taskId, &repeatDoneAt)
		if err != nil {
			return
		}

		// get full info about task from DB
		db.First(&task, taskId)

		// this task is not done yet ever, skip
		if !repeatDoneAt.Valid {
			continue
		}

		var timeToUnDone time.Time
		if task.RepeatUnit == "y" {
			timeToUnDone = repeatDoneAt.Time.AddDate(int(task.RepeatBest), 0, 0)
		} else if task.RepeatUnit == "m" {
			timeToUnDone = repeatDoneAt.Time.AddDate(0, int(task.RepeatBest), 0)
		} else if task.RepeatUnit == "d" {
			timeToUnDone = repeatDoneAt.Time.AddDate(0, 0, int(task.RepeatBest))
		} else if task.RepeatUnit == "h" {
			timeToUnDone = repeatDoneAt.Time.Add(time.Hour * time.Duration(int(task.RepeatBest)))
		} else if task.RepeatUnit == "i" {
			timeToUnDone = repeatDoneAt.Time.Add(time.Minute * time.Duration(int(task.RepeatBest)))
		}

		if time.Now().After(timeToUnDone) {
			log.Printf("Task ID %v is due, setting as undone...", task.Id)
			task.DoneAt.Valid = false
			db.Save(&task)
		}

	}
}
