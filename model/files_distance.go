package model

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DistanceRank struct {
	A        int
	B        int
	Distance int
	Sha256   string
}

const DistanceIdDelimiter string = "-"

func DistanceTopFind(db *sql.DB, a int) (found bool, rank map[int]DistanceRank) {
	//FIXME placeholder didn't work :(
	//rows, err := db.Query(`SELECT a_b, distance FROM files_distance WHERE a_b LIKE "?-%" ORDER BY distance ASC LIMIT 10`, a)
	aStr := strconv.Itoa(a)
	rows, err := db.Query(`SELECT a_b, distance FROM files_distance WHERE a_b LIKE "` + aStr + DistanceIdDelimiter + "%" + `" ORDER BY distance ASC LIMIT 10`)
	checkErr(err)
	defer rows.Close()

	found = false
	var aB string
	var distance int
	rank = make(map[int]DistanceRank)

	i := 0
	for rows.Next() {
		err = rows.Scan(&aB, &distance)
		checkErr(err)

		res := strings.Split(aB, DistanceIdDelimiter)
		a, _ := strconv.Atoi(res[0])
		b, _ := strconv.Atoi(res[1])
		_, file := FindById(db, b)

		rank[i] = DistanceRank{a, b, distance, file.Sha256}

		found = true
		i++
	}
	return found, rank
}

func DistanceTopFindSimilar(db *sql.DB, a int) (found bool, rank map[int]DistanceRank) {
	//FIXME placeholder didn't work :(
	//rows, err := db.Query(`SELECT a_b, distance FROM files_distance WHERE a_b LIKE "?-%" ORDER BY distance ASC LIMIT 10`, a)
	aStr := strconv.Itoa(a)
	rows, err := db.Query(`SELECT a_b, distance FROM files_distance WHERE a_b LIKE "` + aStr + DistanceIdDelimiter + "%" + `" AND distance <= 20 ORDER BY distance ASC LIMIT 10`)
	checkErr(err)
	defer rows.Close()

	found = false
	var aB string
	var distance int
	rank = make(map[int]DistanceRank)

	i := 0
	for rows.Next() {
		err = rows.Scan(&aB, &distance)
		checkErr(err)

		res := strings.Split(aB, DistanceIdDelimiter)
		a, _ := strconv.Atoi(res[0])
		b, _ := strconv.Atoi(res[1])
		_, file := FindById(db, b)

		rank[i] = DistanceRank{a, b, distance, file.Sha256}

		found = true
		i++
	}
	return found, rank
}

func DistanceInsertPrepare(db *sql.DB) (tx *sql.Tx, stmt *sql.Stmt) {
	tx, err := db.Begin()
	stmt, err = tx.Prepare("INSERT OR IGNORE INTO files_distance (a_b, distance) VALUES(?,?)")
	if err != nil {
		fmt.Println(err)
	}
	checkErr(err)
	return tx, stmt
}

func DistanceInsert(stmt *sql.Stmt, a int, b int, dist int) {
	aStr := strconv.Itoa(a)
	bStr := strconv.Itoa(b)
	a_b := aStr + DistanceIdDelimiter + bStr

	_, err := stmt.Exec(a_b, dist)
	checkErr(err)
}

func DistanceInsertEnd(tx *sql.Tx) {
	tx.Commit()
}
