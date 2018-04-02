package model

import (
	"database/sql"
	"log"
	"gitlab.systemz.pl/systemz/tasktab/config"
)

var (
	DB         *sql.DB
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
	var err error
	dbUrl := config.DB_USERNAME + ":" + config.DB_PASSWORD + "@tcp(" + config.DB_HOST + ":" + config.DB_PORT + ")/" + config.DB_NAME + "?charset=utf8"
	DB, err = sql.Open("mysql", dbUrl)
	if err != nil {
		log.Fatalln("open db fail:", err)
	}
	///DB.SetMaxIdleConns(g.Config().Database.Idle)

	err = DB.Ping()
	if err != nil {
		log.Panic("Ping to DB failed")
	}
	log.Printf("%v", "Connection to DB seems OK!")
}
