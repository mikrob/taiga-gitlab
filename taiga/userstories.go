package taiga

import (
	"fmt"
	"regexp"
)

// UserstoriesService handles communication with the user stories related methods of
// the Taiga API.
type UserstoriesService struct {
	client *Client
}

// CreateUserstoryPointOptions represents a Taiga user story point
type CreateUserstoryPointOptions map[string]float64

// Userstory represent a Taiga user story
type Userstory struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	ProjectID   int    `json:"project"`
	Description string `json:"description"`
	Status      int    `json:"status"`
	Milestone   int    `json:"milestone,omitempty"`
	OwnerID     int    `json:"owner"`
	Assigne     int    `json:"assigned_to,omitempty"`
	Version     int    `json:"version"`
}

// CreateUserstoryOptions represents the CreateUserstory() options
type CreateUserstoryOptions struct {
	Subject     string   `json:"subject"`
	ProjectID   int      `json:"project"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Tags        []string `json:"tags"`
	Milestone   int      `json:"milestone,omitempty"`
	Assigne     int      `json:"assigned_to,omitempty"`
	//	Points      *CreateUserstoryPointOptions `json:"points",omitempty`
}

// UserstoryStatus represents a Taiga user story status
type UserstoryStatus struct {
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	ProjectID int    `json:"project"`
	Slug      string `json:"slug"`
}

// CreateCommentUserstoryOptions represents the CreateCommentUserstory() options
type CreateCommentUserstoryOptions struct {
	Comment string `json:"comment"`
	Version int    `json:"version"`
}

// CreateUserstory creates a new project issue.
func (s *UserstoriesService) CreateUserstory(opt *CreateUserstoryOptions) (*Userstory, *Response, error) {
	req, err := s.client.NewRequest("POST", "userstories", opt)
	if err != nil {
		return nil, nil, err
	}

	u := new(Userstory)
	resp, err := s.client.Do(req, u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}

// ListUserstories lists user stories
func (s *UserstoriesService) ListUserstories() ([]*Userstory, *Response, error) {
	req, err := s.client.NewRequest("GET", "userstories", nil)
	if err != nil {
		return nil, nil, err
	}
	var u []*Userstory
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}

//FindUserstoryByRegexName search issues by pattern matching user stories name
func (s *UserstoriesService) FindUserstoryByRegexName(pattern string) (*Userstory, *Response, error) {
	re := regexp.MustCompile(pattern)
	userstories, _, err := s.ListUserstories()
	if err != nil {
		return nil, nil, err
	}
	for _, userstory := range userstories {
		if re.FindString(userstory.Subject) != "" {
			return userstory, nil, err
		}
	}
	return nil, nil, err
}

// ListUserstoryStatuses lists issue status for a given project id
func (s *IssuesService) ListUserstoryStatuses() ([]*UserstoryStatus, *Response, error) {
	req, err := s.client.NewRequest("GET", "userstory-statuses", nil)
	if err != nil {
		return nil, nil, err
	}
	var u []*UserstoryStatus
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}

// CreateCommentUserstory creates a new comment in project issue.
func (s *IssuesService) CreateCommentUserstory(userstoryID int, opt *CreateCommentUserstoryOptions) (*Userstory, *Response, error) {
	url := fmt.Sprintf("userstories/%d", userstoryID)
	req, err := s.client.NewRequest("PATCH", url, opt)
	if err != nil {
		return nil, nil, err
	}
	u := new(Userstory)
	resp, err := s.client.Do(req, u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}
