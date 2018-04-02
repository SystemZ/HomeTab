package main

import (
	_ "gitlab.systemz.pl/systemz/tasktab/config"
	"github.com/robfig/cron"
	"log"
)

var (
	scheduleRunning = false
)

func schedule() {
	c := cron.New()
	c.AddFunc("*/30 * * * * *", func() {
		// run only one instance at a time
		if scheduleRunning {
			log.Println("Cron tasks already running, aborting...")
			return
		}
		log.Println("Starting cron task...")
		scheduleRunning = true
		getTasksForAllGroups()
		scheduleRunning = false
	})
	c.Start()
}

func main() {
	go schedule()
	httpStart()

	/*
	integrations.GmailGetNewTokenStep1()
	token := integrations.GmailGetNewTokenStep2("")
	res, _ := json.Marshal(token)
	log.Printf("%s", res)
	*/
}
