package taiga

import (
	"fmt"
	"regexp"
	"time"
)

// IssuesService handles communication with the issues related methods of
// the Taiga API.
type IssuesService struct {
	client *Client
}

// Issue represent a Taiga issue
type Issue struct {
	ID           int       `json:"id"`
	Subject      string    `json:"subject"`
	ProjectID    int       `json:"project"`
	Description  string    `json:"description"`
	Status       int       `json:"status"`
	Assigne      int       `json:"assigned_to,omitempty"`
	Milestone    int       `json:"milestone,omitempty"`
	OwnerID      int       `json:"owner"`
	Version      int       `json:"version"`
	LastModified time.Time `json:"modified_date"`
}

// CreateIssueOptions represents the CreateIssue() options
type CreateIssueOptions struct {
	Subject     string   `json:"subject"`
	ProjectID   int      `json:"project"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Tags        []string `json:"tags"`
	Assigne     int      `json:"assigned_to,omitempty"`
	Milestone   int      `json:"milestone,omitempty"`
}

// CreateCommentIssueOptions represents the CreateCommentIssue() options
type CreateCommentIssueOptions struct {
	Comment string `json:"comment"`
	Version int    `json:"version"`
}

// IssueStatus represents a Taiga issue status
type IssueStatus struct {
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	ProjectID int    `json:"project"`
	Slug      string `json:"slug"`
}

// CreateIssue creates a new project issue.
func (s *IssuesService) CreateIssue(opt *CreateIssueOptions) (*Issue, *Response, error) {
	req, err := s.client.NewRequest("POST", "issues", opt)
	if err != nil {
		return nil, nil, err
	}
	i := new(Issue)
	resp, err := s.client.Do(req, i)
	if err != nil {
		return nil, resp, err
	}
	return i, resp, err
}

// ListIssues lists issues
func (s *IssuesService) ListIssues() ([]*Issue, *Response, error) {
	req, err := s.client.NewRequest("GET", "issues", nil)
	if err != nil {
		return nil, nil, err
	}
	var i []*Issue
	resp, err := s.client.Do(req, &i)
	if err != nil {
		return nil, resp, err
	}
	return i, resp, err
}

// ListIssueStatuses lists issue status for a given project id
func (s *IssuesService) ListIssueStatuses() ([]*IssueStatus, *Response, error) {
	req, err := s.client.NewRequest("GET", "issue-statuses", nil)
	if err != nil {
		return nil, nil, err
	}
	var i []*IssueStatus
	resp, err := s.client.Do(req, &i)
	if err != nil {
		return nil, resp, err
	}
	return i, resp, err
}

//FindIssueByRegexName search issues by pattern matching issue name
func (s *IssuesService) FindIssueByRegexName(pattern string) ([]*Issue, *Response, error) {
	re := regexp.MustCompile(pattern)
	issues, resp, err := s.ListIssues()
	var matchingIssue []*Issue
	for _, issue := range issues {
		if re.FindString(issue.Subject) != "" {
			matchingIssue = append(matchingIssue, issue)
		}
	}
	return matchingIssue, resp, err
}

// CreateCommentIssue creates a new comment in project issue.
func (s *IssuesService) CreateCommentIssue(issueID int, opt *CreateCommentIssueOptions) (*Issue, *Response, error) {
	url := fmt.Sprintf("issues/%d", issueID)
	req, err := s.client.NewRequest("PATCH", url, opt)
	if err != nil {
		return nil, nil, err
	}
	i := new(Issue)
	resp, err := s.client.Do(req, i)
	if err != nil {
		return nil, resp, err
	}
	return i, resp, err
}

// GetIssueHistory retrieve user story history
func (s *UserstoriesService) GetIssueHistory(userstoryID int) ([]*HistoryEntry, *Response, error) {
	url := fmt.Sprintf("history/issue/%d", userstoryID)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var historyEntries []*HistoryEntry
	resp, err := s.client.Do(req, &historyEntries)
	if err != nil {
		return nil, resp, err
	}
	return historyEntries, resp, err
}
