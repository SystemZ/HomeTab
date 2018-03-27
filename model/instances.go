package model

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"log"
)

var DB *sql.DB

var (
	id         int
	url        string
	token      string
	instanceId int
	typeId     int
)

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func init() {
	//fmt.Println("init model")
	var err error
	DB, err = sql.Open("mysql", "dev:dev@/dev?charset=utf8")
	//user:password@/dbname
	if err != nil {
		log.Fatalln("open db fail:", err)
	}

	////DB.SetMaxIdleConns(g.Config().Database.Idle)
	//
	//err = DB.Ping()
	//if err != nil {
	//	log.Fatalln("ping db fail:", err)
	//}
}

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
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&url, &typeId, &token, &id, &instanceId)
		checkErr(err)
	}
	err = rows.Err()
	checkErr(err)

	return types.Credentials{id, instanceId, url, token, typeId}
}

func GetCredentialByInstanceId(id int) types.Credentials {
	rows, err := DB.Query("SELECT (SELECT url FROM instances WHERE instances.id = instances_access.instance_id) AS url, (SELECT type_id FROM instances WHERE instances.id = instances_access.instance_id) AS type_id, token, instance_user_id, instance_id  FROM instances_access WHERE instance_id = ?", id)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&url, &typeId, &token, &id, &instanceId)
		checkErr(err)
	}
	err = rows.Err()
	checkErr(err)

	return types.Credentials{id, instanceId, url, token, typeId}
}
