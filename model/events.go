package model

import "time"

// returns row ID
func CreateEvent(typeId int, userId int) int64 {
	stmt, err := DB.Prepare("INSERT events SET user_id=?, type_id=?, created_at=?")
	checkErr(err)

	res, err := stmt.Exec(userId, typeId, time.Now().Unix())
	checkErr(err)

	id, err := res.LastInsertId()
	checkErr(err)

	return id
}
