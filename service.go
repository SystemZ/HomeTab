package main

import (
	"gitlab.systemz.pl/systemz/tasktab/integrations"
	"gitlab.systemz.pl/systemz/tasktab/model"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"log"
)

func getTasksForAllGroups() {
	log.Println("Importing tasks for all groups")
	groups := model.GetAllGroupsIds()
	for _, id := range groups {
		getTasksForGroup(id)
	}
}

func getTasksForGroup(groupId int) {
	log.Printf("Importing tasks for group ID %v", groupId)
	accessIds := model.GetAllInstancesAccessIds()
	for _, accessId := range accessIds {
		credentials := model.GetInstanceByAccessId(accessId)
		go GetTasksForCredential(credentials, accessId, groupId)
	}
}

func GetTasksForCredential(credentials types.Credentials, accessId int, groupId int) {
	switch credentials.TypeId {
	case 1:
		log.Printf("Processing GitLab issues for credentials #%v", accessId)
		tasks := integrations.GetAllTasksAssignedToAidFromGitLab(credentials)
		for _, task := range tasks {
			model.ImportGitlabTask(task, credentials.InstanceId, groupId)
		}
	case 2:
		log.Printf("Processing GitHub issues for credentials #%v", accessId)
		tasks := integrations.GetAllIssuesAssignedToGitHubUser(credentials)
		for _, task := range tasks {
			model.ImportGithubTask(task, credentials.InstanceId, groupId)
		}
	case 3:
		log.Printf("Processing Gmail messages for credentials #%v", accessId)
		tasks := integrations.GmailGetInboxUnreadMessages(credentials)
		for _, task := range tasks.Messages {
			log.Printf("%v", task)
			t := integrations.GmailGetMessage(credentials, task.Id)
			model.ImportGmailTask(t, credentials.InstanceId, groupId)
		}
	default:
		log.Printf("%s: %v", "Unknown instance typeID", credentials.TypeId)
	}
}
