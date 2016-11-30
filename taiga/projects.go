package taiga

import "fmt"

// ProjectsService handles communication with the project related methods of
// the Taiga API.
type ProjectsService struct {
	client *Client
}

// Project represents a Taiga project
type Project struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Roles []*Role `json:"roles"`
}

// Role reprensents a Taiga project role
type Role struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Slug string `json:"slug"`
}

// ListProjects retrieve a Taiga project by its slug name
func (s *ProjectsService) ListProjects() ([]*Project, *Response, error) {
	req, err := s.client.NewRequest("GET", "projects", nil)
	if err != nil {
		return nil, nil, err
	}
	var p []*Project
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, err
}

// GetProject retrieve a Taiga project by its id
func (s *ProjectsService) GetProject(id int) (*Project, *Response, error) {
	u := fmt.Sprintf("projects/%d", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	p := new(Project)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}

	return p, resp, err
}

// GetProjectByName retrieve a Taiga project by its name
func (s *ProjectsService) GetProjectByName(name string) (*Project, *Response, error) {
	var p *Project
	projects, resp, err := s.ListProjects()
	if err != nil {
		return nil, resp, err
	}
	for _, project := range projects {
		if project.Name == name {
			return project, nil, nil
		}
	}
	return p, nil, nil
}
