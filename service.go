package main

import (
	"gitlab.systemz.pl/systemz/tasktab/integrations"
	"gitlab.systemz.pl/systemz/tasktab/model"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"log"
	"time"
	"sync"
)

func markAsDone(id int) {
	task := model.GetTaskById(id)
	if task.Done {
		log.Printf("Task already done: %v", id)
	}
	credential := model.GetCredentialByInstanceId(task.InstanceId)
	switch credential.TypeId {
	case 3:
		integrations.GmailMarkMessageAsDone(credential, task.InstanceTaskId)
	}
	model.SetAsDone(task)
}

func originUrl(task model.Task) string {
	switch task.Type {
	case "gmail":
		return "https://mail.google.com/mail/u/0/#inbox/" + task.InstanceTaskId
	}

	return ""
}

func getTasksForAllGroups() {
	log.Println("Importing tasks for all groups")
	groups := model.GetAllGroupsIds()
	for _, id := range groups {
		getTasksForGroup(id)
	}
}

func getTasksForGroup(groupId int) {
	log.Printf("Importing tasks for group ID %v", groupId)
	var wg sync.WaitGroup
	accessIds := model.GetAllInstancesAccessIds()
	for _, accessId := range accessIds {
		credentials := model.GetInstanceByAccessId(accessId)
		wg.Add(1)
		go GetTasksForCredential(credentials, accessId, groupId, &wg)
	}
	wg.Wait()
}

func contains(slice []string, wanted string) bool {
	for _, v := range slice {
		if v == wanted {
			return true
		}
	}
	return false
}

func UpdateTasksForInstance(instanceId int) {
	log.Printf("Updating tasks for instance ID %v", instanceId)
	//instanceIds := model.GetAllInstancesIds()
	//for _, instanceId := range instanceIds {
	credential := model.GetCredentialByInstanceId(instanceId)
	tasks := model.ListTasksForInstance(instanceId)
	for _, task := range tasks {
		//log.Printf("%v", task.InstanceTaskId)
		timestampNow := time.Now().Unix()
		// no update for 10 minutes, need to check
		if timestampNow-int64(task.CheckedAt) > 600 {
			switch credential.TypeId {
			case 3: //gmail
				log.Printf("Updating email with title: %v", task.Title)
				freshMsg := integrations.GmailGetMessage(credential, task.InstanceTaskId)
				if task.Done == contains(freshMsg.LabelIds, "INBOX") {
					if (task.Done) {
						model.SetAsNotDone(task)
					} else {
						model.SetAsDone(task)
					}
				}
			}
		}
		time.Sleep(time.Second)
	}
	//}
}

//func UpdateTasksForCredential(credentials types.Credentials, accessId int, groupId int) {
//	switch credentials.TypeId {
//	case 3:
//		log.Printf("Processing Gmail messages for credentials #%v", accessId)
//		tasks := integrations.GmailGetInboxUnreadMessages(credentials)
//		for _, task := range tasks.Messages {
//			log.Printf("%v", task)
//			t := integrations.GmailGetMessage(credentials, )
//			model.ImportGmailTask(t, credentials.InstanceId, groupId)
//		}
//	default:
//		log.Printf("%s: %v", "Unknown instance typeID", credentials.TypeId)
//	}
//}

func GetTasksForCredential(credentials types.Credentials, accessId int, groupId int, wg *sync.WaitGroup) {
	defer wg.Done()
	switch credentials.TypeId {
	case 1:
		log.Printf("Importing GitLab issues for credentials #%v", accessId)
		tasks := integrations.GetAllTasksAssignedToAidFromGitLab(credentials)
		for _, task := range tasks {
			model.ImportGitlabTask(task, credentials.InstanceId, groupId)
		}
	case 2:
		log.Printf("Importing GitHub issues for credentials #%v", accessId)
		tasks := integrations.GetAllIssuesAssignedToGitHubUser(credentials)
		for _, task := range tasks {
			model.ImportGithubTask(task, credentials.InstanceId, groupId)
		}
	case 3:
		log.Printf("Importing Gmail messages for credentials #%v", accessId)
		tasks := integrations.GmailGetInboxMessages(credentials)
		for _, task := range tasks.Messages {
			t := integrations.GmailGetMessage(credentials, task.Id)
			model.ImportGmailTask(t, credentials.InstanceId, groupId)
		}
	default:
		log.Printf("%s: %v", "Unknown instance typeID", credentials.TypeId)
	}
}
