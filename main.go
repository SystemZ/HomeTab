package main

import (
	"github.com/robfig/cron"
)

func schedule() {
	c := cron.New()
	c.AddFunc("*/25 * * * * *", func() { getTasksForAllGroups() })
	c.AddFunc("*/20 * * * * *", func() { UpdateTasksForInstance(3) })
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
