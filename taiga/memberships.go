package taiga

import "fmt"

// MembershipsService handles communication with the membership related methods of
// the Taiga API.
type MembershipsService struct {
	client *Client
}

// Membership represents a Taiga project membership
type Membership struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user"`
	Email     string `json:"email"`
	RoleID    int    `json:"role"`
	ProjectID int    `json:"project"`
}

// CreateMembershipOptions represents the CreateMembership() options
type CreateMembershipOptions struct {
	Email     string `json:"email"`
	RoleID    int    `json:"role"`
	ProjectID int    `json:"project"`
}

// ListMembershipOptions represents ListMemberships() options
type ListMembershipOptions struct {
	ProjectID int
}

// ListMemberships lists Taiga memberships
func (s *MembershipsService) ListMemberships(opt *ListMembershipOptions) ([]*Membership, *Response, error) {
	url := "memberships"
	if opt.ProjectID > 0 {
		url = fmt.Sprintf("memberships?project=%d", opt.ProjectID)
	}
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var m []*Membership
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}
	return m, resp, err
}

//GetUserInProjectMembership returns a user member of project
func (s *MembershipsService) GetUserInProjectMembership(userID int, projectID int) (*Membership, *Response, error) {
	opts := &ListMembershipOptions{
		ProjectID: projectID,
	}
	memberships, resp, err := s.ListMemberships(opts)
	if err != nil {
		return nil, resp, err
	}
	for _, m := range memberships {
		if m.ProjectID == projectID && m.UserID == userID {
			return m, resp, err
		}
	}
	return nil, resp, err
}

// CreateMembership creates a new project issue.
func (s *MembershipsService) CreateMembership(opt *CreateMembershipOptions) (*Membership, *Response, error) {
	req, err := s.client.NewRequest("POST", "memberships", opt)
	if err != nil {
		return nil, nil, err
	}

	m := new(Membership)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}
	return m, resp, err
}
