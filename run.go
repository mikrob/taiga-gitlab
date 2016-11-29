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
	// fetch useful issue status
	issueStatuses, _, err := taigaClient.Issues.ListIssueStatuses()
	if err != nil {
		panic(err.Error())
	}
	issueStatusClosed := new(taiga.IssueStatus)
	issueStatusNew := new(taiga.IssueStatus)
	issueStatusInprogress := new(taiga.IssueStatus)
	for _, issueStatus := range issueStatuses {
		if issueStatus.ProjectID == taigaProject.ID {
			switch issueStatus.Slug {
			case "closed":
				issueStatusClosed = issueStatus
			case "new":
				issueStatusNew = issueStatus
			case "in-progress":
				issueStatusInprogress = issueStatus
			}

		}
	}
	// fetch useful user story status
	userstoryStatuses, _, err := taigaClient.Issues.ListUserstoryStatuses()
	if err != nil {
		panic(err.Error())
	}
	userstoryStatusDone := new(taiga.UserstoryStatus)
	userstoryStatusNew := new(taiga.UserstoryStatus)
	userstoryStatusInprogress := new(taiga.UserstoryStatus)
	for _, userstoryStatus := range userstoryStatuses {
		if userstoryStatus.ProjectID == taigaProject.ID {
			switch userstoryStatus.Slug {
			case "done":
				userstoryStatusDone = userstoryStatus
			case "new":
				userstoryStatusNew = userstoryStatus
			case "in-progress":
				userstoryStatusInprogress = userstoryStatus
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
	userstoryStatus := new(taiga.UserstoryStatus)
	var objectToCreate string

	for _, issue := range issues {
		var tags []string
		tags = append(tags, projectName)
		issueSubjectPrefix := fmt.Sprintf("gitlab/%s/%d", projectName, issue.IID)
		issueSubject := fmt.Sprintf("%s %s", issueSubjectPrefix, issue.Title)
		for _, label := range issue.Labels {
			switch label {
			default:
				objectToCreate = "issue"
			case "functionnal":
				objectToCreate = "userstory"
			}
		}

		milestone := new(taiga.Milestone)
		if issue.Milestone != nil {
			m, _, _ := taigaClient.Milestones.FindMilestoneByName(issue.Milestone.Title, taigaProject.ID)
			if len(m) == 1 {
				milestone = m[0]
			} else {
				milestoneStart := issue.Milestone.StartDate
				milestoneFinish := issue.Milestone.DueDate
				if issue.Milestone.StartDate == "" {
					milestoneStart = fmt.Sprintf("%d-%02d-%02d",
						issue.Milestone.CreatedAt.Year(),
						issue.Milestone.CreatedAt.Month(),
						issue.Milestone.CreatedAt.Day())

				}
				if issue.Milestone.DueDate == "" {
					milestoneFinish = fmt.Sprintf("%d-%02d-%02d",
						issue.Milestone.CreatedAt.Year(),
						issue.Milestone.CreatedAt.Month(),
						issue.Milestone.CreatedAt.Day())
				}
				opt := &taiga.CreateMilestoneOptions{
					ProjectID:       taigaProject.ID,
					Name:            issue.Milestone.Title,
					EstimatedStart:  milestoneStart,
					EstimatedFinish: milestoneFinish,
				}
				m, _, err := taigaClient.Milestones.CreateMilestone(opt)
				if err != nil {
					log.Fatal("Cannot create milestone ", fmt.Sprintf("%+v", opt), err.Error())
				}
				milestone = m
			}
		}
		if objectToCreate == "issue" {
			switch {
			case issue.State == "closed":
				issueStatus = issueStatusClosed
			case issue.Assignee.ID > 0:
				issueStatus = issueStatusInprogress
			default:
				issueStatus = issueStatusNew
			}

			i := &taiga.CreateIssueOptions{
				Subject:     issueSubject,
				ProjectID:   taigaProject.ID,
				Description: fmt.Sprintf("Gitlab issue: %s/%s/issues/%d\n\n%s", gitlabURL, projectName, issue.IID, issue.Description),
				Status:      issueStatus.ID,
				Tags:        tags,
			}
			if milestone.ID > 0 {
				i.Milestone = milestone.ID
			}
			searchIssues, _, _ := taigaClient.Issues.FindIssueByRegexName(issueSubjectPrefix)
			if len(searchIssues) == 0 {
				issue, _, err := taigaClient.Issues.CreateIssue(i)
				if err != nil {
					log.Print(err.Error())
					continue
				}
				fmt.Println("Created issue", issue.ID, issueStatus.Name)
			}
		} else if objectToCreate == "userstory" {
			switch {
			case issue.State == "closed":
				userstoryStatus = userstoryStatusDone
			case issue.Assignee.ID > 0:
				userstoryStatus = userstoryStatusInprogress
			default:
				userstoryStatus = userstoryStatusNew
			}
			u := &taiga.CreateUserstoryOptions{
				Subject:     issueSubject,
				ProjectID:   taigaProject.ID,
				Description: fmt.Sprintf("Gitlab issue: %s/%s/issues/%d\n\n%s", gitlabURL, projectName, issue.IID, issue.Description),
				Status:      userstoryStatus.ID,
				Tags:        tags,
			}
			if milestone.ID > 0 {
				u.Milestone = milestone.ID
			}
			searchUserstories, _, _ := taigaClient.Userstories.FindUserstoryByRegexName(issueSubjectPrefix)
			if len(searchUserstories) == 0 {
				userstory, _, err := taigaClient.Userstories.CreateUserstory(u)
				if err != nil {
					log.Print(err.Error())
					continue
				}
				fmt.Println("Created user story", userstory.ID, userstoryStatus.Name)
			}
		}
	}
}
