package taiga

//MilestonesService handles communication with the milestones related methods of
// the Taiga API.
type MilestonesService struct {
	client *Client
}

// Milestone represent a Taiga milestone
type Milestone struct {
	ID              int    `json:"id"`
	Name            string `json:"name"`
	ProjectID       int    `json:"project"`
	EstimatedStart  string `json:"estimated_start"`
	EstimatedFinish string `json:"estimated_finish"`
}

// CreateMilestoneOptions represents the CreateMilestone() options
type CreateMilestoneOptions struct {
	Name            string `json:"name"`
	ProjectID       int    `json:"project"`
	EstimatedStart  string `json:"estimated_start"`
	EstimatedFinish string `json:"estimated_finish"`
}

// CreateMilestone creates a new project issue.
func (s *MilestonesService) CreateMilestone(opt *CreateMilestoneOptions) (*Milestone, *Response, error) {
	req, err := s.client.NewRequest("POST", "milestones", opt)
	if err != nil {
		return nil, nil, err
	}

	m := new(Milestone)
	resp, err := s.client.Do(req, m)
	if err != nil {
		return nil, resp, err
	}
	return m, resp, err
}

// ListMilestones list milestones
func (s *MilestonesService) ListMilestones() ([]*Milestone, *Response, error) {
	req, err := s.client.NewRequest("GET", "milestones", nil)
	if err != nil {
		return nil, nil, err
	}
	var m []*Milestone
	resp, err := s.client.Do(req, &m)
	if err != nil {
		return nil, resp, err
	}
	return m, resp, err
}

//FindMilestoneByName search issues by pattern matching issue name
func (s *MilestonesService) FindMilestoneByName(name string, pid int) ([]*Milestone, *Response, error) {
	var matchingMilestone []*Milestone
	milestones, resp, err := s.ListMilestones()
	if err != nil {
		return nil, resp, err
	}
	for _, milestone := range milestones {
		if milestone.Name == name {
			if pid > 0 {
				if milestone.ProjectID == pid {
					matchingMilestone = append(matchingMilestone, milestone)
				}
			} else {
				matchingMilestone = append(matchingMilestone, milestone)
			}
		}
	}
	return matchingMilestone, resp, err
}
