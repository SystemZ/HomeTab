package main

import (
	"github.com/icrowley/fake"
	"gitlab.com/systemz/tasktab/model"
	"gitlab.com/systemz/tasktab/web"
)

func main() {
	model.InitMysql()
	model.InitRedis()

	//model.CreateUser("fake", fake.EmailAddress(), "pass")
	model.CreateUser(fake.UserName(), fake.EmailAddress(), fake.SimplePassword())
	groupId := model.CreateGroup("Test group")
	model.CreateProject("Test project", groupId)

	web.StartWebInterface()
}
