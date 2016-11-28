package taiga

// IssuesService handles communication with the issues related methods of
// the Taiga API.
type IssuesService struct {
	client *Client
}

// Issue represent a Taiga issue
type Issue struct {
	ID        int    `json:"id"`
	Subject   string `json:"subject"`
	ProjectID int    `json:"project"`
}

// CreateIssueOptions represents the CreateIssue() options
type CreateIssueOptions struct {
	Subject   string `json:"subject"`
	ProjectID int    `json:"project"`
	//TypeID    int    `json:"type"`
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
