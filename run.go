package main

import (
	"fmt"
	"log"

	gitlab "github.com/xanzy/go-gitlab"

	"taiga-gitlab/taiga"
)

func main() {
	taigaUsername := "admin"
	taigaPassword := "123123"
	taigaURL := "http://192.168.99.102"
	taigaClient := taiga.NewClient(nil, taigaUsername, taigaPassword)
	taigaProjectName := "ufancyme"
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", taigaURL))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		panic(err.Error())
	}
	taigaProject, _, err := taigaClient.Projects.GetProjectByName(taigaProjectName)
	if err != nil {
		panic(err.Error())
	}
	issueStatuses, _, err := taigaClient.Issues.ListIssueStatuses()
	if err != nil {
		panic(err.Error())
	}
	issueStatusClosed := new(taiga.IssueStatus)
	issueStatusNew := new(taiga.IssueStatus)
	for _, issueStatus := range issueStatuses {
		if issueStatus.ProjectID == taigaProject.ID {
			switch issueStatus.Slug {
			case "closed":
				issueStatusClosed = issueStatus
			case "new":
				issueStatusNew = issueStatus
			}

		}
	}
	fmt.Println("Project Name:", taigaProject.Name)

	gitlabToken := "1qVsgb99XFst2GRwBXxn"
	gitlabURL := "https://gitlab.botsunit.com"
	projectName := "boobs/payment"
	git := gitlab.NewClient(nil, gitlabToken)
	git.SetBaseURL(fmt.Sprintf("%s/api/v3", gitlabURL))
	project, _, err := git.Projects.GetProject(projectName)
	if err != nil {
		panic(err.Error())
	}
	issuesOptions := &gitlab.ListProjectIssuesOptions{}
	issues, _, err := git.Issues.ListProjectIssues(project.ID, issuesOptions)
	if err != nil {
		panic(err.Error())
	}
	issueStatus := new(taiga.IssueStatus)
	for _, issue := range issues {
		switch issue.State {
		case "closed":
			issueStatus = issueStatusClosed
		default:
			issueStatus = issueStatusNew
		}
		fmt.Println(issueStatus.Name)
		i := &taiga.CreateIssueOptions{
			Subject:     fmt.Sprintf("gitlab/%s/%d %s", projectName, issue.IID, issue.Title),
			ProjectID:   taigaProject.ID,
			Description: fmt.Sprintf("Gitlab issue: http://gitlab.botsunit.com/%s/issues/%d\n\n%s", projectName, issue.IID, issue.Description),
			Status:      issueStatus.ID,
		}
		issue, _, err := taigaClient.Issues.CreateIssue(i)
		if err != nil {
			log.Print(err.Error())
			continue
		}
		fmt.Println("Created issue", issue.ID, issueStatus.Name)

	}
}
