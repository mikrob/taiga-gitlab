package taiga

import (
	"fmt"
	"regexp"
	"time"
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
	ID             int            `json:"id"`
	Subject        string         `json:"subject"`
	ProjectID      int            `json:"project"`
	Description    string         `json:"description"`
	Status         int            `json:"status"`
	Milestone      int            `json:"milestone,omitempty"`
	OwnerID        int            `json:"owner"`
	Assigne        int            `json:"assigned_to,omitempty"`
	Version        int            `json:"version"`
	Points         map[string]int `json:"points"`
	LastModified   time.Time      `json:"modified_date"`
	Ref            int            `json:"ref"`
	PointsPerRoles map[Role]Point
	AssignedUser   string
	TotalPoint     float64
	ElapsedTime    float64
	Overtaking     bool
	Undertaking    bool
	RightTime      bool
	Color          string
	TaskList       []*Task
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

func internalListUserStories(client *Client) ([]*Userstory, *Response, error) {
	req, err := client.NewRequest("GET", "userstories", nil)
	if err != nil {
		return nil, nil, err
	}
	var u []*Userstory
	resp, err := client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}

// ListUserstories lists user stories
func (s *UserstoriesService) ListUserstories() ([]*Userstory, *Response, error) {
	return internalListUserStories(s.client)
}

// ListUserstoriesForMilestone lists user stories
func (s *UserstoriesService) ListUserstoriesForMilestone(milestone Milestone) ([]*Userstory, *Response, error) {
	u, resp, err := internalListUserStories(s.client)
	var result []*Userstory
	for _, us := range u {
		if us.Milestone == milestone.ID {
			result = append(result, us)
		}
	}
	return result, resp, err
}

// ListUserstoriesForUser lists user stories
func (s *UserstoriesService) ListUserstoriesForUser(user User) ([]*Userstory, *Response, error) {
	u, resp, err := internalListUserStories(s.client)
	var result []*Userstory
	for _, us := range u {
		if us.Assigne == user.ID {
			result = append(result, us)
		}
	}
	return result, resp, err
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

//HistoryValues represent history values
type HistoryValues struct {
	//Status map[string]string `json:"status"`
	Status []string `json:"status"`
}

//HistroyDiff represent history diff
type HistroyDiff struct {
	Status []int `json:"status"`
}

//HistoryEntry represent an history entry
type HistoryEntry struct {
	Comment          string        `json:"comment"`
	ID               string        `json:"id"`
	Type             int           `json:"type"`
	HistoryValueList HistoryValues `json:"values_diff"`
	Diff             HistroyDiff   `json:"diff"`
	CreatedAt        time.Time     `json:"created_at"`
}

// GetUserStoryHistory retrieve user story history
func (s *UserstoriesService) GetUserStoryHistory(userstoryID int) ([]*HistoryEntry, *Response, error) {
	url := fmt.Sprintf("history/userstory/%d", userstoryID)
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

//UserStoryCustomAttribute represent userstory custom attributes
type UserStoryCustomAttribute struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

//GetUserStoryCustomAttributes retrieve custom attributes available for US
func (s *UserstoriesService) GetUserStoryCustomAttributes() ([]*UserStoryCustomAttribute, *Response, error) {
	req, err := s.client.NewRequest("GET", "userstory-custom-attributes", nil)
	if err != nil {
		return nil, nil, err
	}
	var customAttributes []*UserStoryCustomAttribute
	resp, err := s.client.Do(req, &customAttributes)
	if err != nil {
		return nil, resp, err
	}
	return customAttributes, resp, err
}

//CustomAttributeValues represent a custom attribute value
type CustomAttributeValues struct {
	Values      map[string]string `json:"attributes_values"`
	UserStoryID int               `json:"user_story"`
}

//GetUserStoryCustomAttributeValue return value for a custom attribute
func (s *UserstoriesService) GetUserStoryCustomAttributeValue(userStoryID int) (*CustomAttributeValues, *Response, error) {
	url := fmt.Sprintf("userstories/custom-attributes-values/%d", userStoryID)
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	customAttributesValues := new(CustomAttributeValues)
	resp, err := s.client.Do(req, customAttributesValues)
	if err != nil {
		return nil, nil, err
	}
	return customAttributesValues, resp, err
}
