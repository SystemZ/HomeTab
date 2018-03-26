package integrations

import (
	"github.com/xanzy/go-gitlab"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"log"
	"time"
)

func GetAllTasksAssignedToAidFromGitLab(credentials types.Credentials) []*gitlab.Issue {
	lab := gitlab.NewClient(nil, credentials.Token)
	lab.SetBaseURL(credentials.Url)

	issues := lab.Issues
	sort := "desc"

	page := 1
	finished := false

	var result []*gitlab.Issue

	for ; finished == false; page++ {
		listOptions := gitlab.ListOptions{Page: page, PerPage: 100}
		issuesOptions := &gitlab.ListIssuesOptions{Sort: &sort, AssigneeID: &credentials.UserIdOnInstance, ListOptions: listOptions}
		iss, res, err := issues.ListIssues(issuesOptions)

		if err != nil {
			log.Panic("Error with gitlab ListIssues req")
		}

		if res.NextPage == 0 {
			finished = true
		}
		result = append(result, iss...)

		//FIXME configurable cooldown
		time.Sleep(time.Second * 1)
	}
	return result
}

func GetTasksFromProjectIdFromGitlab(projectId int, credentials types.Credentials) {
	lab := gitlab.NewClient(nil, credentials.Token)
	lab.SetBaseURL(credentials.Url)

	issues := lab.Issues
	sort := "desc"

	page := 1
	finished := false

	for ; finished == false; page++ {
		listOptions := gitlab.ListOptions{Page: page, PerPage: 50}
		projectIssuesOptions := &gitlab.ListProjectIssuesOptions{Sort: &sort, ListOptions: listOptions}
		iss, res, err := issues.ListProjectIssues(projectId, projectIssuesOptions)
		if err != nil {
			log.Println(err)
			log.Panic("Error with gitlab ListProjectIssues req")
		}

		for _, v := range iss {
			//log.Printf("key=%v, value=%v", k, v)
			log.Printf("%v", v.Title, v.WebURL)
		}

		if res.NextPage == 0 {
			finished = true
		}

		//FIXME configurable cooldown
		time.Sleep(time.Second * 1)
	}

}
