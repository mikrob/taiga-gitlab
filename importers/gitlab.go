package importers

import (
	"fmt"
	"log"
	"taiga-gitlab/taiga"

	"github.com/urfave/cli"
	gitlab "github.com/xanzy/go-gitlab"
)

// Proxy bridges taiga and gitlab client
type Proxy struct {
	taiga        *taiga.Client
	gitlab       *gitlab.Client
	taigaProject *taiga.Project
}

// TaskStatusEquiv returns a task status ID based on userstory status ID
func TaskStatusEquiv(userstoryStatusID int) int {
	switch userstoryStatusID {
	case 1: // new
		return 1 // new
	case 2: //ready
		return 1 // new
	case 3: // in progress
		return 2 // in progress
	case 4: // ready for test
		return 3
	case 5: // done
		return 4 // closed
	case 6: // archived
		return 4 // closed
	default:
		return 1 // new
	}
}

// ImportGitlabUser sync Gitlab user to Taiga
func (p *Proxy) ImportGitlabUser(gitlabUser *gitlab.User) (*taiga.User, error) {
	taigaUser, _, err := p.taiga.Users.FindUserByUsername(gitlabUser.Username)
	if err != nil {
		return nil, err
	}
	if taigaUser == nil {
		return nil, fmt.Errorf("Please create Taiga following user:\nusername: %s\nemail: %s\nname: %s\n",
			gitlabUser.Username,
			gitlabUser.Email,
			gitlabUser.Name)
	}
	// ensure user is member of taiga project

	m, _, err := p.taiga.Memberships.GetUserInProjectMembership(taigaUser.ID, p.taigaProject.ID)
	if err != nil {
		log.Fatal(err.Error())
	}
	if m == nil {
		createMembershipOpts := &taiga.CreateMembershipOptions{
			RoleID:    1,
			Email:     taigaUser.Email,
			ProjectID: p.taigaProject.ID,
		}
		_, _, err := p.taiga.Memberships.CreateMembership(createMembershipOpts)
		if err != nil {
			return nil, fmt.Errorf("Cannot create membership for %s, %s", taigaUser.Username, err.Error())
		}
	}

	return taigaUser, nil
}

// ImportGitlab2Taiga imports Gitlab issues, milestones to Taiga
func ImportGitlab2Taiga(c *cli.Context) error {
	requiredFlagsStrings := []string{
		"taiga-url", "taiga-user", "taiga-password", "taiga-project",
		"gitlab-url", "gitlab-token", "gitlab-project",
	}
	for _, flag := range requiredFlagsStrings {
		if c.String(flag) == "" {
			return cli.NewExitError(fmt.Sprintf("%s undefined", flag), 1)
		}
	}
	taigaUsername := c.String("taiga-user")
	taigaPassword := c.String("taiga-password")
	taigaURL := c.String("taiga-url")
	taigaClient := taiga.NewClient(nil, taigaUsername, taigaPassword)
	taigaProjectName := c.String("taiga-project")
	taigaClient.SetBaseURL(fmt.Sprintf("%s/api/v1", taigaURL))
	_, _, err := taigaClient.Users.Login()
	if err != nil {
		panic(err.Error())
	}
	taigaProject, _, err := taigaClient.Projects.GetProjectByName(taigaProjectName)
	if err != nil {
		panic(err.Error())
	}
	if taigaProject == nil {
		log.Fatal("No such project ", taigaProjectName)
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

	gitlabToken := c.String("gitlab-token")
	gitlabURL := c.String("gitlab-url")
	projectName := c.String("gitlab-project")
	git := gitlab.NewClient(nil, gitlabToken)
	git.SetBaseURL(fmt.Sprintf("%s/api/v3", gitlabURL))
	project, _, err := git.Projects.GetProject(projectName)
	if err != nil {
		panic(err.Error())
	}

	listOpts := &gitlab.ListOptions{
		PerPage: 100,
	}
	issuesOptions := &gitlab.ListProjectIssuesOptions{
		ListOptions: *listOpts,
	}
	issues, _, err := git.Issues.ListProjectIssues(project.ID, issuesOptions)
	if err != nil {
		panic(err.Error())
	}
	issueStatus := new(taiga.IssueStatus)
	userstoryStatus := new(taiga.UserstoryStatus)
	//var objectToCreate string
	objectToCreate := "issue"
	z := &Proxy{taiga: taigaClient, gitlab: git, taigaProject: taigaProject}
	for _, issue := range issues {
		// sync author
		issueAuthor, _, err := git.Users.GetUser(issue.Author.ID)
		if err != nil {
			log.Fatalf("unable to found Gitlab user %s", issue.Author.Name)
		}
		if c.Bool("taiga-skip-user") == false {
			_, err = z.ImportGitlabUser(issueAuthor)
			if err != nil {
				log.Fatalf("Cannot sync user %s from Gitlab to Taiga: %s", issueAuthor.Username, err.Error())
			}
		}
		// sync assignee
		issueAssigneTaiga := new(taiga.User)
		if issue.Assignee.ID > 0 && c.Bool("taiga-skip-user") == false {
			issueAssigneGitlab, _, _err := git.Users.GetUser(issue.Assignee.ID)
			if _err != nil {
				log.Fatalf("unable to found Gitlab user %s", issueAssigneGitlab.Name)
			}
			issueAssigneTaiga, err = z.ImportGitlabUser(issueAssigneGitlab)
			if err != nil {
				log.Fatalf("Cannot sync user %s from Gitlab to Taiga: %s", issueAssigneGitlab.Username, err.Error())
			}
		}

		// sync creator
		var tags []string
		tags = append(tags, projectName)
		issueSubjectPrefix := fmt.Sprintf("GITLAB-%d %s", issue.ID, projectName)
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
			if issueAssigneTaiga.ID > 0 {
				i.Assigne = issueAssigneTaiga.ID
			}
			searchIssues, _, _ := taigaClient.Issues.FindIssueByRegexName(issueSubjectPrefix)
			if len(searchIssues) == 0 {
				taigaIssue, _, err := taigaClient.Issues.CreateIssue(i)
				if err != nil {
					log.Fatalf("Cannot create issue %s", err.Error())
				}
				log.Println("Created issue", taigaIssue.ID, taigaIssue.Subject)
				listNotesOpts := gitlab.ListIssueNotesOptions{}
				notes, _, _ := git.Notes.ListIssueNotes(project.ID, issue.ID, &listNotesOpts)
				for _, note := range notes {
					commentOpts := &taiga.CreateCommentIssueOptions{
						Comment: fmt.Sprintf("Author: %s\n\n%s", note.Author.Name, note.Body),
						Version: taigaIssue.Version,
					}
					taigaIssuePatched, _, err := taigaClient.Issues.CreateCommentIssue(taigaIssue.ID, commentOpts)
					if err != nil {
						log.Fatal("Cannot create comment")
					}
					log.Printf("Create new comment %+v", taigaIssuePatched)
				}

				//	taigaClient.Issues.CreateCommentIssue(issue.ID, )
			} else {
				log.Printf("Gitlab issue found in Taiga %+v", searchIssues)
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
			if issueAssigneTaiga.ID > 0 {
				u.Assigne = issueAssigneTaiga.ID
			}
			var relatedUserStory *taiga.Userstory
			existingUserstory, _, _ := taigaClient.Userstories.FindUserstoryByRegexName(issueSubjectPrefix)
			if existingUserstory == nil {
				userstory, _, err := taigaClient.Userstories.CreateUserstory(u)
				if err != nil {
					log.Fatal(err.Error())
				}
				fmt.Println("Created user story", userstory.ID, userstoryStatus.Name)
				relatedUserStory = userstory
				listNotesOpts := gitlab.ListIssueNotesOptions{}
				notes, _, _ := git.Notes.ListIssueNotes(project.ID, issue.ID, &listNotesOpts)
				for _, note := range notes {
					commentOpts := &taiga.CreateCommentUserstoryOptions{
						Comment: fmt.Sprintf("Author: %s\n\n%s", note.Author.Name, note.Body),
						Version: userstory.Version,
					}
					taigaUserstoryPatched, _, err := taigaClient.Issues.CreateCommentUserstory(userstory.ID, commentOpts)
					if err != nil {
						log.Fatal("Cannot create comment")
					}
					log.Printf("Create new comment %+v", taigaUserstoryPatched)
				}
			} else {
				log.Printf("Gitlab issue found in Taiga %+v", existingUserstory)
				relatedUserStory = existingUserstory
			}
			// Create associated task
			if c.Bool("create-task") == true {
				searchTask, _, err := taigaClient.Tasks.FindTaskByRegexName(issueSubjectPrefix)
				if err != nil {
					log.Fatal("Cannot search task")
				}
				if searchTask == nil {
					taskStatusID := TaskStatusEquiv(relatedUserStory.Status)
					newTaskOpts := &taiga.CreateTaskOptions{
						Subject:     relatedUserStory.Subject,
						ProjectID:   relatedUserStory.ProjectID,
						UserstoryID: relatedUserStory.ID,
						Status:      taskStatusID,
					}
					if relatedUserStory.Milestone > 0 {
						newTaskOpts.Milestone = relatedUserStory.Milestone
					}
					if relatedUserStory.Assigne > 0 {
						newTaskOpts.Assigne = relatedUserStory.Assigne
					}
					newTask, _, err := taigaClient.Tasks.CreateTask(newTaskOpts)
					if err != nil {
						log.Fatal("Cannot create task")
					}
					log.Println("Create new task ", newTask.Subject)
				} else {
					log.Println("Associated task already exists")
				}
			}
		}
	}
	return nil
}
