package taiga

import "regexp"

// UserstoriesService handles communication with the user stories related methods of
// the Taiga API.
type UserstoriesService struct {
	client *Client
}

// Userstory represent a Taiga user story
type Userstory struct {
	ID          int      `json:"id"`
	Subject     string   `json:"subject"`
	ProjectID   int      `json:"project"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Tags        []string `json:"tags"`
	Milestone   int      `json:"milestone,omitempty"`
}

// CreateUserstoryOptions represents the CreateUserstory() options
type CreateUserstoryOptions struct {
	Subject     string   `json:"subject"`
	ProjectID   int      `json:"project"`
	Description string   `json:"description"`
	Status      int      `json:"status"`
	Tags        []string `json:"tags"`
	Milestone   int      `json:"milestone,omitempty"`
}

// UserstoryStatus represents a Taiga user story status
type UserstoryStatus struct {
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	ProjectID int    `json:"project"`
	Slug      string `json:"slug"`
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
func (s *UserstoriesService) FindUserstoryByRegexName(pattern string) ([]*Userstory, *Response, error) {
	re := regexp.MustCompile(pattern)
	userstories, resp, err := s.ListUserstories()
	var matchingUserstories []*Userstory
	for _, userstory := range userstories {
		if re.FindString(userstory.Subject) != "" {
			matchingUserstories = append(matchingUserstories, userstory)
		}
	}
	return matchingUserstories, resp, err
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
