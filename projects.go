package golang_todoist_api

import (
	"context"
	"log"
	"net/http"
	"net/url"
	"os"
)

type Option func(*Client)

type ProjectsResponse struct {
	Projects []Project
	TodoistResponse
}
type ProjectResponse struct {
	Project Project
	TodoistResponse
}

type Project struct {
	ID             string `json:"id"`
	ParentId       string `json:"parent_id"`
	Order          int    `json:"order"`
	Color          string `json:"color"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	Url            string `json:"url"`
	ViewStyle      string `json:"view_style"`
}

func OptionAPIURL(u string) func(*Client) {
	return func(c *Client) { c.endpoint = u }
}

func New(token string, options ...Option) *Client {
	s := &Client{
		token:      token,
		endpoint:   APIURL,
		httpclient: &http.Client{},
		log:        log.New(os.Stderr, "volyanyk/golang-todoist-api", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func (api *Client) GetProjects() (*[]Project, error) {
	return api.GetProjectsContext(context.Background())
}
func (api *Client) GetProjectById(id string) (*Project, error) {
	return api.GetProjectByIdContext(id, context.Background())
}

func (api *Client) GetProjectsContext(context context.Context) (*[]Project, error) {
	response := &ProjectsResponse{}

	err := api.getMethod(context,
		"projects",
		api.token,
		url.Values{},
		&response.Projects)

	return &response.Projects, err
}

func (api *Client) GetProjectByIdContext(id string, context context.Context) (*Project, error) {
	response := &ProjectResponse{}

	err := api.getMethod(context,
		"projects/"+id,
		api.token,
		url.Values{},
		&response.Project)

	return &response.Project, err
}
