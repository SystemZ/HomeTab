package model

import (
	"github.com/google/go-github/github"
	"github.com/xanzy/go-gitlab"
	"google.golang.org/api/gmail/v1"
	"log"
	"strconv"
	"strings"
)

var (
	title          string
	checkedAt      int
	updatedAt      int
	createdAt      int
	instanceTaskId string
	projectTaskId  string
)

type Task struct {
	Id             int    `json:"id"`
	InstanceId     int    `json:"instanceId"`
	InstanceTaskId string `json:"instanceTaskId"`
	ProjectId      int    `json:"projectId"`
	ProjectTaskId  string `json:"projectTaskId"`
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

func GetTaskByStringId(taskIdString string) Task {
	taskId, _ := strconv.ParseInt(taskIdString, 10, 64)
	return getTask("taskId", int(taskId))[0]
}

func GetTasksByInstanceTaskId(instanceId int, instanceTaskId string) []Task {
	return getTaskInstance("instanceTaskId", instanceId, instanceTaskId)
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
	query := "SELECT id, instance_task_id, project_id, project_task_id, instance_id, done, title, (SELECT type_id FROM instances WHERE instances.id = tasks.instance_id) AS type_id, checked_at, updated_at, created_at FROM tasks "
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
	defer stmt.Close()

	rows, err := stmt.Query(id)
	checkErr(err)

	var taskType, done, projectId int
	var doneBool bool

	defer rows.Close()
	var result []Task
	for rows.Next() {
		err := rows.Scan(&id, &instanceTaskId, &projectId, &projectTaskId, &instanceId, &done, &title, &taskType, &checkedAt, &updatedAt, &createdAt)
		checkErr(err)
		if done >= 1 {
			doneBool = true
		} else {
			doneBool = false
		}
		result = append(result, Task{
			Id:             id,
			InstanceTaskId: instanceTaskId,
			ProjectId:      projectId,
			ProjectTaskId:  projectTaskId,
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

func getTaskInstance(typeId string, id int, idString string) []Task {
	query := "SELECT id, instance_task_id, project_id, project_task_id, instance_id, done, title, (SELECT type_id FROM instances WHERE instances.id = tasks.instance_id) AS type_id, checked_at, updated_at, created_at FROM tasks "
	switch typeId {
	case "instanceTaskId":
		query += "WHERE instance_id = ? AND instance_task_id = ?"
	}
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(id, idString)
	checkErr(err)

	var taskType, done, projectId int
	var doneBool bool

	defer rows.Close()
	var result []Task
	for rows.Next() {
		err := rows.Scan(&id, &instanceTaskId, &projectId, &projectTaskId, &instanceId, &done, &title, &taskType, &checkedAt, &updatedAt, &createdAt)
		checkErr(err)
		if done >= 1 {
			doneBool = true
		} else {
			doneBool = false
		}
		result = append(result, Task{
			Id:             id,
			InstanceTaskId: instanceTaskId,
			ProjectId:      projectId,
			ProjectTaskId:  projectTaskId,
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
func ImportGitlabTask(issue *gitlab.Issue, instanceId int, groupId int, projectId int) (alreadyExists bool, task Task) {
	getTask := GetTasksByInstanceTaskId(instanceId, strconv.Itoa(issue.ID))
	if len(getTask) > 0 {
		return true, getTask[0]
	}

	done := 0
	if issue.State == "closed" {
		done = 1
	}

	stmt, err := DB.Prepare("INSERT IGNORE tasks SET done = ?, title = ?, instance_task_id = ?, project_task_id = ?, project_id = ?, created_at = ?, instance_id = ?, group_id = ?")
	defer stmt.Close()
	checkErr(err)

	_, err = stmt.Exec(done, issue.Title, issue.ID, strconv.Itoa(issue.IID), projectId, issue.CreatedAt.Unix(), instanceId, groupId)
	checkErr(err)

	getTask = GetTasksByInstanceTaskId(instanceId, strconv.Itoa(issue.ID))
	return false, getTask[0]
}

// returns row ID
func ImportGithubTask(issue *github.Issue, instanceId int, groupId int, projectId int) (alreadyExists bool, task Task) {
	getTask := GetTasksByInstanceTaskId(instanceId, strconv.Itoa(int(*issue.ID)))
	if len(getTask) > 0 {
		return true, getTask[0]
	}

	done := 0
	if *issue.State == "closed" {
		done = 1
	}

	stmt, err := DB.Prepare("INSERT IGNORE tasks SET done = ?, title = ?, instance_task_id = ?, project_task_id = ?, project_id = ?, created_at = ?, instance_id = ?, group_id = ?")
	defer stmt.Close()
	checkErr(err)

	_, err = stmt.Exec(done, issue.Title, issue.ID, issue.Number, projectId, issue.CreatedAt.Unix(), instanceId, groupId)
	checkErr(err)

	getTask = GetTasksByInstanceTaskId(instanceId, strconv.Itoa(int(*issue.ID)))
	return false, getTask[0]
}

// returns row ID
func ImportGmailTask(message *gmail.Message, instanceId int, groupId int) int64 {
	task := GetTasksByInstanceTaskId(instanceId, message.Id)
	if len(task) > 0 {
		return 0
	}

	var subject string

	for _, header := range message.Payload.Headers {
		if header.Name == "Subject" {
			log.Printf("Mail ID: %v, Subject: %v", message.Id, header.Value)
			subject = header.Value
		}
	}

	// e-mails in INBOX are always not finished, so default done = 0 is OK
	stmt, err := DB.Prepare("INSERT IGNORE tasks SET title = ?, instance_task_id = ?, created_at = ?, instance_id = ?, group_id = ?")
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(subject, message.Id, message.InternalDate/1000, instanceId, groupId)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

func RefreshDone(task Task, state string) {
	//github and gitlab
	if task.Done && state == "closed" {
		return
	}
	//github
	if !task.Done && state == "open" {
		return
	}
	//gitlab
	if !task.Done && state == "opened" {
		return
	}

	stmt, err := DB.Prepare("UPDATE tasks SET done=? WHERE id = ?")
	defer stmt.Close()
	checkErr(err)

	if task.Done && state == "open" {
		_, err = stmt.Exec(0, task.Id)
	}
	if task.Done && state == "opened" {
		_, err = stmt.Exec(0, task.Id)
	}
	if !task.Done && state == "closed" {
		_, err = stmt.Exec(1, task.Id)
	}

	checkErr(err)
}

func SetAsDone(task Task) {
	stmt, err := DB.Prepare("UPDATE tasks SET done=? WHERE id = ?")
	defer stmt.Close()
	checkErr(err)

	_, err = stmt.Exec(1, task.Id)
	checkErr(err)
}

func SetAsNotDone(task Task) {
	stmt, err := DB.Prepare("UPDATE tasks SET done=? WHERE id = ?")
	defer stmt.Close()
	checkErr(err)

	_, err = stmt.Exec(0, task.Id)
	checkErr(err)
}

func MarkAllGmailTasksDoneExcept(instanceId int, whitelist []string) {
	realWhitelist := "'" + strings.Join(whitelist, "','") + "'"
	// mark tasks as not done if moved from archived to inbox
	stmt, err := DB.Prepare("UPDATE tasks SET done = 0 WHERE instance_id = ? AND tasks.instance_task_id IN ( ? )")
	defer stmt.Close()
	checkErr(err)
	toDoRes, err := stmt.Exec(instanceId, realWhitelist)
	todoAffected, _ := toDoRes.RowsAffected()
	log.Printf("Emails set as todo: %v", todoAffected)
	checkErr(err)

	// NOT IN ( ? ) doesn't work, have to make hack like this :/
	blacklist := " AND tasks.instance_task_id != '" + strings.Join(whitelist, "' AND tasks.instance_task_id != '") + "'"
	// mark all tasks as done if not in inbox
	stmt2, err := DB.Prepare("UPDATE tasks SET done = 1 WHERE instance_id = ?" + blacklist)
	defer stmt2.Close()
	checkErr(err)
	doneRes, err := stmt2.Exec(instanceId)
	doneAffected, _ := doneRes.RowsAffected()
	log.Printf("Emails set as done: %v", doneAffected)
	checkErr(err)
}
