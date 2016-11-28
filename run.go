package main

import (
	"fmt"
	"log"

	gitlab "github.com/xanzy/go-gitlab"

	"taiga-gitlab/taiga"
)

func main() {
	taigaUsername := "admin"
	taigaPassword := "secret"
	taigaURL := "http://192.168.99.102"
	taigaClient := taiga.NewClient(nil, taigaUsername, taigaPassword)
	taigaProjectName := "test"
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", taigaURL))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		panic(err.Error())
	}
	taigaProject, _, err := taigaClient.Projects.GetProjectByName(taigaProjectName)
	if err != nil {
		panic(err.Error())
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
	for _, issue := range issues {
		fmt.Println("Creating ", issue.Title)
		i := &taiga.CreateIssueOptions{
			Subject:   issue.Title,
			ProjectID: taigaProject.ID,
			//TypeID:    1,
		}
		issue, _, err := taigaClient.Issues.CreateIssue(i)
		if err != nil {
			log.Print(err.Error())
			continue
		}
		fmt.Println("Created issue", issue.ID)

	}
}
