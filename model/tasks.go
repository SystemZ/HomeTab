package model

import (
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/api/gmail/v1"
	"log"
)

var (
	title          string
	checkedAt      int
	updatedAt      int
	createdAt      int
	instanceTaskId string
)

type Task struct {
	Id             int    `json:"id"`
	InstanceId     int    `json:"instanceId"`
	InstanceTaskId string `json:"instanceTaskId"`
	Done           bool   `json:"done"`
	Type           string `json:"type"`
	Title          string `json:"title"`
	CheckedAt      int    `json:"checkedAt"`
	UpdatedAt      int    `json:"updatedAt"`
	CreatedAt      int    `json:"createdAt"`
}

func GetTaskById(taskId int) Task {
	return getTask("taskId", taskId)[0]
}

func ListTasksForGroup(groupId int) []Task {
	return getTask("groupId", groupId)
}

func ListTasksToDoForGroup(groupId int) []Task {
	return getTask("groupIdToDo", groupId)
}

func ListTasksForInstance(instanceIdInt int) []Task {
	return getTask("instanceId", instanceIdInt)
}

func getTask(typeId string, id int) []Task {
	query := "SELECT id, instance_task_id, instance_id, done, title, (SELECT type_id FROM instances WHERE instances.id = tasks.instance_id) AS type_id, checked_at, updated_at, created_at FROM tasks "
	switch typeId {
	case "taskId":
		query += "WHERE id = ? LIMIT 1"
	case "groupId":
		query += "WHERE group_id = ?"
	case "groupIdToDo":
		query += "WHERE done = 0 AND group_id = ?"
	case "instanceId":
		query += "WHERE instance_id = ?"
	}
	stmt, err := DB.Prepare(query)
	checkErr(err)
	rows, err := stmt.Query(id)
	checkErr(err)

	var taskType, done int
	var doneBool bool

	defer rows.Close()
	var result []Task
	for rows.Next() {
		err := rows.Scan(&id, &instanceTaskId, &instanceId, &done, &title, &taskType, &checkedAt, &updatedAt, &createdAt)
		checkErr(err)
		if done >= 1 {
			doneBool = true
		} else {
			doneBool = false
		}
		result = append(result, Task{
			Id:             id,
			InstanceTaskId: instanceTaskId,
			InstanceId:     instanceId,
			Done:           doneBool,
			Title:          title,
			Type:           taskTypePretty(taskType),
			CheckedAt:      checkedAt,
			UpdatedAt:      updatedAt,
			CreatedAt:      createdAt,
		})
	}
	return result
}

func taskTypePretty(typeId int) string {
	if typeId == 1 {
		return "gitlab"
	} else if typeId == 2 {
		return "github"
	} else if typeId == 3 {
		return "gmail"
	}
	return "unknown"
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

func SetAsDone(task Task) {
	stmt, err := DB.Prepare("UPDATE tasks SET done=? WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(1, task.Id)
	checkErr(err)
}

func SetAsNotDone(task Task) {
	stmt, err := DB.Prepare("UPDATE tasks SET done=? WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(0, task.Id)
	checkErr(err)
}
