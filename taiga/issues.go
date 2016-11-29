package taiga

// IssuesService handles communication with the issues related methods of
// the Taiga API.
type IssuesService struct {
	client *Client
}

// Issue represent a Taiga issue
type Issue struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	ProjectID   int    `json:"project"`
	Description string `json:"description"`
	Status      int    `json:"status"`
}

// CreateIssueOptions represents the CreateIssue() options
type CreateIssueOptions struct {
	Subject     string `json:"subject"`
	ProjectID   int    `json:"project"`
	Description string `json:"description"`
	Status      int    `json:"status"`
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
