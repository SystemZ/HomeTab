package model

import (
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type Project struct {
	Id                int    `json:"id"`
	InstanceId        int    `json:"instanceId"`
	InstanceProjectId int    `json:"instanceProjectId"`
	Title             string `json:"title"`
	CreatedAt         int    `json:"createdAt"`
	UpdatedAt         int    `json:"updatedAt"`
}

func GetProjectByInstanceIdAndProjectId(instanceId int, projectId int) []Project {
	query := "SELECT id, instance_id, instance_project_id, title, created_at, updated_at FROM projects WHERE instance_id = ? AND id = ?"
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(instanceId, projectId)
	checkErr(err)

	defer rows.Close()
	var result []Project
	for rows.Next() {
		err := rows.Scan(&id, &instanceId, &projectId, &title, &createdAt, &updatedAt)
		checkErr(err)
		result = append(result, Project{
			Id:                id,
			InstanceId:        instanceId,
			InstanceProjectId: projectId,
			Title:             title,
			UpdatedAt:         updatedAt,
			CreatedAt:         createdAt,
		})
	}
	return result
}

func GetProjectByInstanceIdAndInstanceProjectId(instanceId int, instanceProjectId int) []Project {
	query := "SELECT id, instance_id, instance_project_id, title, created_at, updated_at FROM projects WHERE instance_id = ? AND instance_project_id = ?"
	stmt, err := DB.Prepare(query)
	checkErr(err)
	defer stmt.Close()

	rows, err := stmt.Query(instanceId, instanceProjectId)
	checkErr(err)

	defer rows.Close()
	var result []Project
	for rows.Next() {
		err := rows.Scan(&id, &instanceId, &instanceProjectId, &title, &createdAt, &updatedAt)
		checkErr(err)
		result = append(result, Project{
			Id:                id,
			InstanceId:        instanceId,
			InstanceProjectId: instanceProjectId,
			Title:             title,
			UpdatedAt:         updatedAt,
			CreatedAt:         createdAt,
		})
	}
	return result
}

func CreateProject(instanceId int, instanceProjectId int, title string) int {
	projects := GetProjectByInstanceIdAndInstanceProjectId(instanceId, instanceProjectId)
	if len(projects) > 0 {
		return projects[0].Id
	}

	stmt, err := DB.Prepare("INSERT projects SET instance_id = ?, instance_project_id = ?, title = ?, created_at = ?")
	defer stmt.Close()
	checkErr(err)

	res, err := stmt.Exec(instanceId, instanceProjectId, title, time.Now().Unix())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return int(id)
}
