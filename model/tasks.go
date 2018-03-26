package model

import (
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/api/gmail/v1"
	"log"
)

var (
	title string
)

type Task struct {
	Id    int    `json:"id"`
	Title string `json:"title"`
}

func ListTasksForGroup(groupId int) []Task {
	stmt, err := DB.Prepare("SELECT id, title FROM tasks WHERE group_id = ?")
	checkErr(err)

	rows, err := stmt.Query(groupId)
	checkErr(err)

	defer rows.Close()
	var result []Task
	for rows.Next() {
		err := rows.Scan(&id, &title)
		checkErr(err)
		result = append(result, Task{Id: id, Title: title})
	}
	return result
}

// returns row ID
func ImportGitlabTask(issue *gitlab.Issue, instanceId int, groupId int) int64 {
	stmt, err := DB.Prepare("INSERT IGNORE tasks SET title = ?, instance_task_id = ?, created_at = ?, instance_id = ?, group_id = ?")
	checkErr(err)

	res, err := stmt.Exec(issue.Title, issue.ID, issue.CreatedAt.Unix(), instanceId, groupId)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

// returns row ID
func ImportGithubTask(issue *github.Issue, instanceId int, groupId int) int64 {
	stmt, err := DB.Prepare("INSERT IGNORE tasks SET title = ?, instance_task_id = ?, created_at = ?, instance_id = ?, group_id = ?")
	checkErr(err)

	res, err := stmt.Exec(issue.Title, issue.ID, issue.CreatedAt.Unix(), instanceId, groupId)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

// returns row ID
func ImportGmailTask(message *gmail.Message, instanceId int, groupId int) int64 {
	var subject string

	for _, header := range message.Payload.Headers {
		if header.Name == "Subject" {
			log.Printf("Mail ID: %v, Subject: %v", message.Id, header.Value)
			subject = header.Value
		}
	}

	stmt, err := DB.Prepare("INSERT IGNORE tasks SET title = ?, instance_task_id = ?, created_at = ?, instance_id = ?, group_id = ?")
	checkErr(err)

	res, err := stmt.Exec(subject, message.Id, message.InternalDate/1000, instanceId, groupId)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}
