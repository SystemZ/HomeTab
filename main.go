package main

import (
	"github.com/icrowley/fake"
	"gitlab.com/systemz/tasktab/model"
)

func main() {
	model.InitMysql()
	model.CreateUser(fake.UserName(), fake.EmailAddress(), fake.SimplePassword())
	groupId := model.CreateGroup("Test group")
	model.CreateProject("Test project", groupId)
}
