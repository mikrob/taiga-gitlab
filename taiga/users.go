package taiga

import "fmt"

// UsersService handles communication with the user related methods of
// the Taiga API.
type UsersService struct {
	client *Client
}

// User represents a Taiga user
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	IsActive bool   `json:"is_active"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	LoginType string `json:"type"`
}

// Login represents the auth response
type Login struct {
	Token string `json:"auth_token"`
}

// CreateUserOptions represents the CreateUser() options
type CreateUserOptions struct {
	Username string
	Password string
}

// CurrentUser retrieve current logged user
func (s *UsersService) CurrentUser() (*User, *Response, error) {
	req, err := s.client.NewRequest("GET", "users/me", nil)
	if err != nil {
		return nil, nil, err
	}
	usr := new(User)
	resp, err := s.client.Do(req, usr)
	if err != nil {
		return nil, resp, err
	}

	return usr, resp, err
}

// ListUsers lists Taiga users
func (s *UsersService) ListUsers() ([]*User, *Response, error) {
	req, err := s.client.NewRequest("GET", "users", nil)
	if err != nil {
		return nil, nil, err
	}
	var u []*User
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}

// GetUser lists Taiga users
func (s *UsersService) GetUser(uid int) (*User, *Response, error) {
	uri := fmt.Sprintf("users/%d", uid)
	req, err := s.client.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, nil, err
	}
	var u *User
	resp, err := s.client.Do(req, &u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, err
}

// Login retrieve current logged user
func (s *UsersService) Login() (*Login, *Response, error) {
	loginRequest := LoginRequest{
		Username:  s.client.username,
		Password:  s.client.password,
		LoginType: loginType,
	}
	req, err := s.client.NewRequest("POST", "auth", loginRequest)
	if err != nil {
		return nil, nil, err
	}
	login := new(Login)
	resp, err := s.client.Do(req, login)
	if err != nil {
		return nil, resp, err
	}
	s.client.Token = login.Token
	return login, resp, err
}

// FindUserByUsername fetch users by name
func (s *UsersService) FindUserByUsername(name string) (*User, *Response, error) {
	users, resp, err := s.ListUsers()
	if err != nil {
		return nil, resp, err
	}
	for _, user := range users {
		if user.Username == name {
			return s.GetUser(user.ID)
		}
	}
	return nil, resp, err
}

// CreateUser creates a new Taiga user
// TODO this API does not exists
// func (s *UsersService) CreateUser(opt *CreateUserOptions) (*User, *Response, error) {
// 	s.client.SetBaseURL()
// 	req, err := s.client.NewRequest("POST", "users", opt)
// 	if err != nil {
// 		return nil, nil, err
// 	}
// 	u := new(User)
// 	resp, err := s.client.Do(req, u)
// 	if err != nil {
// 		return nil, resp, err
// 	}
// 	return u, resp, err
// }
