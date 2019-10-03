package main

import (
	"github.com/icrowley/fake"
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/web"
	"time"
)

func main() {
	model.InitMysql()
	model.InitRedis()

	model.CreateUser(fake.UserName(), fake.EmailAddress(), fake.SimplePassword())
	groupId := model.CreateGroup("Test group")
	projectId := model.CreateProject("Test project", groupId)
	now := time.Now()
	task := model.Task{
		Subject:          fake.Sentence(),
		ProjectId:        projectId,
		AssignedUserId:   0,
		Repeating:        0,
		NeverEnding:      0,
		RepeatUnit:       "",
		RepeatMin:        0,
		RepeatBest:       0,
		RepeatMax:        0,
		RepeatFrom:       &now,
		EstimateS:        0,
		MasterTaskId:     0,
		SeparateChildren: 0,
	}
	model.CreateTask(task)

	web.StartWebInterface()
}
