package model

import (
	"database/sql"
	"log"
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
