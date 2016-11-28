package main

import (
	"fmt"

	"taiga-gitlab/taiga"
)

func main() {
	taigaUsername := "admin"
	taigaPassword := "secret"
	taigaURL := "http://192.168.99.102"
	taiga := taiga.NewClient(nil, taigaUsername, taigaPassword)
	taiga.SetBaseURL(fmt.Sprintf("%s/api/v1", taigaURL))
	_, _, err := taiga.Users.Login()
	fmt.Println("token:", taiga.Token)
	if err != nil {
		panic(err.Error())
	}
	me, _, err := taiga.Users.CurrentUser()
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("%+v", me)

	// gitlabToken := "1qVsgb99XFst2GRwBXxn"
	// gitlabURL := "https://gitlab.botsunit.com"
	// projectName := "boobs/payment"
	// git := gitlab.NewClient(nil, gitlabToken)
	// git.SetBaseURL(fmt.Sprintf("%s/api/v3", gitlabURL))
	// project, _, err := git.Projects.GetProject(projectName)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// issuesOptions := &gitlab.ListProjectIssuesOptions{}
	// issues, _, err := git.Issues.ListProjectIssues(project.ID, issuesOptions)
	// if err != nil {
	// 	panic(err.Error())
	// }
	// fmt.Printf("%+v", issues)
}
