package taiga

import "fmt"

// PointsService handles communication with the points related methods of
// the Taiga API.
type PointsService struct {
	client *Client
}

// Point represent a Taiga point
type Point struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	ProjectID int     `json:"project"`
	Order     int     `json:"order"`
	Value     float64 `json:"value"`
}

// CreatePointOptions represents the CreatePoint() options
type CreatePointOptions struct {
	Name      string  `json:"name"`
	ProjectID int     `json:"project"`
	Order     int     `json:"order,omitempty"`
	Value     float64 `json:"value"`
}

// ListPointsOptions represents ListPoints() options
type ListPointsOptions struct {
	ProjectID int
}

// CreatePoint creates a new project issue.
func (s *PointsService) CreatePoint(opt *CreatePointOptions) (*Point, *Response, error) {
	req, err := s.client.NewRequest("POST", "points", opt)
	if err != nil {
		return nil, nil, err
	}

	p := new(Point)
	resp, err := s.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, err
}

// ListPoints lists points
func (s *PointsService) ListPoints(opt *ListPointsOptions) ([]*Point, *Response, error) {
	url := "points"
	if opt.ProjectID > 0 {
		url = fmt.Sprintf("%s?project=%d", url, opt.ProjectID)
	}
	req, err := s.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	var p []*Point
	resp, err := s.client.Do(req, &p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, err
}
