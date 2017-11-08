package main

import (
	"database/sql"
	"log"
	"strings"
)

type Tag struct {
	Id      int
	Fid     int
	Name    string
	Overall int
}

func dbTagList(db *sql.DB, fid int) (found bool, tags map[int]Tag) {
	found = false

	// query
	//rows, err := db.Query("SELECT id, name FROM tags WHERE fid = ? ", fid)
	rows, err := db.Query("SELECT id, name, (SELECT COUNT(id) FROM tags t2 WHERE t1.name = t2.name) as overall FROM tags t1 WHERE fid = ? ORDER BY overall DESC", fid)
	checkErr(err)

	i := 0
	for rows.Next() {
		var id int
		var name string
		var overall int
		err = rows.Scan(&id, &name, &overall)
		checkErr(err)

		if tags == nil {
			tags = make(map[int]Tag)
		}
		tags[i] = Tag{id, fid, name, overall}

		if !found {
			found = true
		}
		i++
	}
	rows.Close() //good habit to close
	return found, tags
}

func dbTagFind(db *sql.DB, name string, fid int) (found bool) {
	// query
	rows, err := db.Query("SELECT id FROM tags WHERE fid = ? AND name = ?", fid, name)
	checkErr(err)
	found = false

	for rows.Next() {
		found = true
		break
	}
	rows.Close() //good habit to close
	return found
}

func dbTagFindSert(db *sql.DB, name string, fid int) {
	found := dbTagFind(db, name, fid)
	if !found {
		dbTagInsert(db, name, fid)
	}
}

func dbTagInsert(db *sql.DB, name string, fid int) (id int64) {
	if len(strings.TrimSpace(name)) <= 1 {
		log.Print("Tag is empty or have less than 2 letters")
		return 0
	}

	stmt, err := db.Prepare("INSERT INTO tags(name, fid) VALUES(?,?)")
	checkErr(err)

	res, err := stmt.Exec(name, fid)
	checkErr(err)

	id, err = res.LastInsertId()
	checkErr(err)

	return id
}

func dbTagDelete(db *sql.DB, id int) {
	log.Printf("%v", id)
	stmt, err := db.Prepare("DELETE FROM tags WHERE id = ?")
	checkErr(err)

	_, err = stmt.Exec(id)
	//res, err = stmt.Exec(id)
	checkErr(err)

	//affect, err = res.RowsAffected()
	//checkErr(err)

	//db.Close()
}

type TagRank struct {
	Name    string
	Overall int
}

func dbTagRank(db *sql.DB) (found bool, tags map[int]TagRank) {
	found = false

	// query
	rows, err := db.Query("SELECT DISTINCT name, (SELECT COUNT(id) FROM tags t2 WHERE t1.name = t2.name) as overall FROM tags t1 ORDER BY overall DESC")
	checkErr(err)

	i := 0
	for rows.Next() {
		var name string
		var overall int
		err = rows.Scan(&name, &overall)
		checkErr(err)

		if tags == nil {
			tags = make(map[int]TagRank)
		}
		tags[i] = TagRank{name, overall}

		if !found {
			found = true
		}
		i++
	}
	rows.Close() //good habit to close
	return found, tags
}
