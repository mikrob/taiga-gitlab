package taiga

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/go-querystring/query"
)

const (
	libraryVersion = "0.1.0"
	userAgent      = "go-taiga"
	defaultBaseURL = "http://localhost:8000/api/v1"
	loginType      = "normal"
)

// Client manages communication with Taiga
type Client struct {
	client    *http.Client
	baseURL   *url.URL
	username  string
	password  string
	UserAgent string
	Token     string

	Users       *UsersService
	Projects    *ProjectsService
	Issues      *IssuesService
	Userstories *UserstoriesService
	Milestones  *MilestonesService
	Memberships *MembershipsService
	Points      *PointsService
}

// Response is a Taiga API response
type Response struct {
	*http.Response
	//Paginated bool
}

// SetBaseURL sets the base URL for API requests to a custom endpoint. urlStr
// should always be specified with a trailing slash.
func (c *Client) SetBaseURL(urlStr string) error {
	// Make sure the given URL end with a slash
	if !strings.HasSuffix(urlStr, "/") {
		urlStr += "/"
	}

	var err error
	c.baseURL, err = url.Parse(urlStr)

	return err
}

// NewClient returns a new Taiga API client. If a nil httpClient is
// provided, http.DefaultClient will be used. To use API methods which require
// authentication, provide a valid username and password.
func NewClient(httpClient *http.Client, username string, password string) *Client {
	return newClient(httpClient, username, password)
}

func newClient(httpClient *http.Client, username string, password string) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	c := &Client{client: httpClient, username: username, password: password, UserAgent: userAgent}
	if err := c.SetBaseURL(defaultBaseURL); err != nil {
		// should never happen since defaultBaseURL is our constant
		panic(err)
	}
	c.Users = &UsersService{client: c}
	c.Projects = &ProjectsService{client: c}
	c.Issues = &IssuesService{client: c}
	c.Userstories = &UserstoriesService{client: c}
	c.Milestones = &MilestonesService{client: c}
	c.Memberships = &MembershipsService{client: c}
	c.Points = &PointsService{client: c}
	return c
}

// NewRequest creates an API request. A relative URL path can be provided in
// urlStr, in which case it is resolved relative to the base URL of the Client.
/// Relative URL paths should always be specified without a preceding slash. If
// specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method, path string, opt interface{}) (*http.Request, error) {
	u := *c.baseURL
	// Set the encoded opaque data
	u.Opaque = c.baseURL.Path + path
	if opt != nil {
		q, err := query.Values(opt)
		if err != nil {
			return nil, err
		}
		u.RawQuery = q.Encode()

	}

	req := &http.Request{
		Method:     method,
		URL:        &u,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header:     make(http.Header),
		Host:       u.Host,
	}

	if method == "POST" || method == "PUT" || method == "PATCH" {
		bodyBytes, err := json.Marshal(opt)
		if err != nil {
			return nil, err
		}
		bodyReader := bytes.NewReader(bodyBytes)

		u.RawQuery = ""
		req.Body = ioutil.NopCloser(bodyReader)
		req.ContentLength = int64(bodyReader.Len())
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	return req, nil
}

// Do sends an API request and returns the API response. The API response is
// JSON decoded and stored in the value pointed to by v, or returned as an
// error if an API error has occurred. If v implements the io.Writer
// interface, the raw response body will be written to v, without attempting to
// first decode it.
func (c *Client) Do(req *http.Request, v interface{}) (*Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	response := newResponse(resp)

	err = CheckResponse(resp)
	if err != nil {
		// even though there was an error, we still return the response
		// in case the caller wants to inspect it further
		return response, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
		} else {
			err = json.NewDecoder(resp.Body).Decode(v)
		}
	}
	return response, err
}

// CheckResponse checks the API response for errors, and returns them if
// present.  A response is considered an error if it has a status code outside
// the 200 range.  API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse.  Any other
// response body will be silently ignored.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; 200 <= c && c <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// Error represents an error
type Error struct {
	Resource string `json:"resource"` // resource on which the error occurred
	Field    string `json:"field"`    // field on which the error occurred
	Code     string `json:"code"`     // validation error code
}

func (e *Error) Error() string {
	return fmt.Sprintf("%v error caused by %v field on %v resource",
		e.Code, e.Field, e.Resource)
}

// An ErrorResponse reports one or more errors caused by an API request.
//
// GitLab API docs:
// http://doc.gitlab.com/ce/api/README.html#data-validation-and-error-reporting
type ErrorResponse struct {
	Response *http.Response // HTTP response that caused this error
	Message  string         `json:"message"` // error message
	Errors   []Error        `json:"errors"`  // more detail on individual errors
}

func (r *ErrorResponse) Error() string {
	path, _ := url.QueryUnescape(r.Response.Request.URL.Opaque)
	ru := fmt.Sprintf("%s://%s%s", r.Response.Request.URL.Scheme, r.Response.Request.URL.Host, path)

	return fmt.Sprintf("%v %s: %d %v %+v",
		r.Response.Request.Method, ru, r.Response.StatusCode, r.Message, r.Errors)
}

// newResponse creats a new Response for the provided http.Response.
func newResponse(r *http.Response) *Response {
	response := &Response{Response: r}
	//response.populatePageValues()
	return response
}
