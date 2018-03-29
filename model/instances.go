package model

import (
	_ "github.com/go-sql-driver/mysql"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"log"
)

// returns row ID
func CreateInstance(typeId int, url string, token string) int64 {
	stmt, err := DB.Prepare("INSERT instances SET type_id=?, url=?, token=?")
	checkErr(err)

	res, err := stmt.Exec(typeId, url, token)
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}

func GetAllInstancesIds() []int {
	rows, err := DB.Query("SELECT id FROM instances")
	checkErr(err)
	defer rows.Close()
	var result []int
	for rows.Next() {
		err := rows.Scan(&id)
		checkErr(err)
		result = append(result, id)
	}
	return result
}

func GetAllInstancesAccessIds() []int {
	rows, err := DB.Query("SELECT id FROM instances_access")
	checkErr(err)
	defer rows.Close()
	var result []int
	for rows.Next() {
		err := rows.Scan(&id)
		checkErr(err)
		result = append(result, id)
	}
	return result
}

func GetInstanceByAccessId(id int) types.Credentials {
	//rows, err := DB.Query("SELECT url, (SELECT token FROM instances_access WHERE instances_access.instance_id = instances.id) AS token FROM instances WHERE id = ? LIMIT 1", id)
	rows, err := DB.Query("SELECT (SELECT url FROM instances WHERE instances.id = instances_access.instance_id) AS url, (SELECT type_id FROM instances WHERE instances.id = instances_access.instance_id) AS type_id, token, instance_user_id, instance_id  FROM instances_access WHERE id = ?", id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&url, &typeId, &token, &id, &instanceId)
		checkErr(err)
	}
	err = rows.Err()
	checkErr(err)

	return types.Credentials{UserIdOnInstance: id, InstanceId: instanceId, Url: url, Token: token, TypeId: typeId}
}

type Instance struct {
	InstanceId int
	Url        string
	TypeId     int
}

func GetInstanceById(id int) Instance {
	//rows, err := DB.Query("SELECT url, (SELECT token FROM instances_access WHERE instances_access.instance_id = instances.id) AS token FROM instances WHERE id = ? LIMIT 1", id)
	rows, err := DB.Query("SELECT id, type_id, url FROM instances WHERE id = ? LIMIT 1", id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&id, &typeId, &url)
		checkErr(err)
	}
	err = rows.Err()
	checkErr(err)

	return Instance{InstanceId: id, TypeId: typeId, Url: url}
}

func GetCredentialByInstanceId(id int) types.Credentials {
	rows, err := DB.Query("SELECT (SELECT url FROM instances WHERE instances.id = instances_access.instance_id) AS url, (SELECT type_id FROM instances WHERE instances.id = instances_access.instance_id) AS type_id, token, instance_user_id, instance_id  FROM instances_access WHERE instance_id = ?", id)
	defer rows.Close()
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		err := rows.Scan(&url, &typeId, &token, &id, &instanceId)
		checkErr(err)
	}
	err = rows.Err()
	checkErr(err)

	return types.Credentials{UserIdOnInstance: id, InstanceId: instanceId, Url: url, Token: token, TypeId: typeId}
}
