package main

import (
	"github.com/robfig/cron"
)

func schedule() {
	c := cron.New()
	//c.AddFunc("*/15 * * * * *", func() { getTasksForAllGroups() })
	c.AddFunc("0 */5 * * * *", func() { getTasksForAllGroups() })
	c.Start()
}

func main() {
	go schedule()
	httpStart()
}
