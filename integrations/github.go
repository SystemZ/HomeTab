package integrations

import (
	"github.com/google/go-github/github"
	"gitlab.systemz.pl/systemz/tasktab/types"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"log"
	"time"
)

func GetAllIssuesAssignedToGitHubUser(credentials types.Credentials) []*github.Issue {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: credentials.Token},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	opt := &github.IssueListOptions{State: "all", ListOptions: github.ListOptions{PerPage: 100}}
	var allIssues []*github.Issue

	for {
		issues, resp, err := client.Issues.List(ctx, true, opt)
		if err != nil {
			log.Printf("%v", err)
		}
		allIssues = append(allIssues, issues...)
		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage

		//FIXME configurable cooldown
		time.Sleep(time.Second * 1)
	}

	return allIssues
}
