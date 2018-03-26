package model

import (
	_ "github.com/go-sql-driver/mysql"
)

func GetAllUsersIds() []int {
	rows, err := DB.Query("SELECT id FROM users")
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
