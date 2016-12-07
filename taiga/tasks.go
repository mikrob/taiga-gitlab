package taiga

import "regexp"

// TasksService handles communication with the tasks related methods of
// the Taiga API.
type TasksService struct {
	client *Client
}

// Task represent a Taiga task
type Task struct {
	ID          int    `json:"id"`
	Subject     string `json:"subject"`
	ProjectID   int    `json:"project"`
	UserstoryID int    `json:"user_story"`
	Status      int    `json:"status"`
	Assigne     int    `json:"assigned_to,omitempty"`
	Milestone   int    `json:"milestone,omitempty"`
}

// CreateTaskOptions represents the CreateTask() options
type CreateTaskOptions struct {
	Subject     string `json:"subject"`
	ProjectID   int    `json:"project"`
	UserstoryID int    `json:"user_story"`
	Status      int    `json:"status"`
	Assigne     int    `json:"assigned_to,omitempty"`
	Milestone   int    `json:"milestone,omitempty"`
}

// TaskStatus represents a Taiga task status
type TaskStatus struct {
	ID        int    `json:"id"`
	IsClosed  bool   `json:"is_closed"`
	Name      string `json:"name"`
	ProjectID int    `json:"project"`
	Slug      string `json:"slug"`
}

// CreateTask creates a new project task.
func (s *TasksService) CreateTask(opt *CreateTaskOptions) (*Task, *Response, error) {
	req, err := s.client.NewRequest("POST", "tasks", opt)
	if err != nil {
		return nil, nil, err
	}
	t := new(Task)
	resp, err := s.client.Do(req, t)
	if err != nil {
		return nil, resp, err
	}
	return t, resp, err
}

// ListTasks lists issues
func (s *TasksService) ListTasks() ([]*Task, *Response, error) {
	req, err := s.client.NewRequest("GET", "tasks", nil)
	if err != nil {
		return nil, nil, err
	}
	var t []*Task
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}
	return t, resp, err
}

// FindTaskByRegexName search issues by pattern matching task name
func (s *TasksService) FindTaskByRegexName(pattern string) (*Task, *Response, error) {
	re := regexp.MustCompile(pattern)
	tasks, resp, err := s.ListTasks()
	if err != nil {
		return nil, resp, err
	}
	for _, task := range tasks {
		if re.FindString(task.Subject) != "" {
			return task, resp, nil
		}
	}
	return nil, resp, err
}

// ListTaskStatuses lists issue status for a given project id
func (s *TasksService) ListTaskStatuses() ([]*TaskStatus, *Response, error) {
	req, err := s.client.NewRequest("GET", "task-statuses", nil)
	if err != nil {
		return nil, nil, err
	}
	var t []*TaskStatus
	resp, err := s.client.Do(req, &t)
	if err != nil {
		return nil, resp, err
	}
	return t, resp, err
}

// FindTaskStatusByRegexName search issues by pattern matching task status
func (s *TasksService) FindTaskStatusByRegexName(pattern string) (*TaskStatus, *Response, error) {
	re := regexp.MustCompile(pattern)
	statuses, resp, err := s.ListTaskStatuses()
	if err != nil {
		return nil, resp, err
	}
	for _, status := range statuses {
		if re.FindString(status.Name) != "" {
			return status, nil, nil
		}
	}
	return nil, nil, nil
}
