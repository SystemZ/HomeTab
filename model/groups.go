package model

import (
	_ "github.com/go-sql-driver/mysql"
)

func GetAllGroupsIds() []int {
	rows, err := DB.Query("SELECT id FROM groups")
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
