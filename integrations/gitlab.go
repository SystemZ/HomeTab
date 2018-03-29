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
