package todoist

import (
	"context"
	"encoding/json"
	"net/url"
)

type ProjectsResponse struct {
	Projects []Project
	TodoistResponse
}
type ProjectResponse struct {
	Project Project
	TodoistResponse
}
type CollaboratorsResponse struct {
	Collaborators []Collaborator
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
type Collaborator struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type AddProjectRequest struct {
	Name       string `json:"name"`        // Required
	ParentId   string `json:"parent_id"`   // Optional
	Color      string `json:"color"`       // Optional
	IsFavorite bool   `json:"is_favorite"` // Optional
	ViewStyle  string `json:"view_style"`  // Optional
}
type UpdateProjectRequest struct {
	Name       string `json:"name"`        // Required
	Color      string `json:"color"`       // Optional
	IsFavorite bool   `json:"is_favorite"` // Optional
	ViewStyle  string `json:"view_style"`  // Optional
}

func (api *Client) GetProjects() (*[]Project, error) {
	return api.GetProjectsContext(context.Background())
}
func (api *Client) GetProjectCollaborators(id string) (*[]Collaborator, error) {
	return api.GetProjectCollaboratorsContext(id, context.Background())
}
func (api *Client) AddProject(request AddProjectRequest) (*Project, error) {
	return api.AddProjectContext(request, context.Background())
}
func (api *Client) UpdateProject(id string, updateProjectRequest UpdateProjectRequest) (*Project, error) {
	return api.UpdateProjectContext(id, updateProjectRequest, context.Background())
}
func (api *Client) GetProjectById(id string) (*Project, error) {
	return api.GetProjectByIdContext(id, context.Background())
}
func (api *Client) DeleteProjectById(id string) (*TodoistResponse, error) {
	return api.DeleteProjectByIdContext(id, context.Background())
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

func (api *Client) GetProjectCollaboratorsContext(id string, context context.Context) (*[]Collaborator, error) {
	response := &CollaboratorsResponse{}

	err := api.getMethod(context,
		"projects/"+id+"/collaborators",
		api.token,
		url.Values{},
		&response.Collaborators)

	return &response.Collaborators, err
}

func (api *Client) AddProjectContext(addProjectRequest AddProjectRequest, context context.Context) (*Project, error) {
	response := &ProjectResponse{}

	request, err := json.Marshal(addProjectRequest)

	err = performPost(context, api.httpclient, api.endpoint+"projects", api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Project, nil
}

func (api *Client) UpdateProjectContext(id string, updateProjectRequest UpdateProjectRequest, context context.Context) (*Project, error) {
	response := &ProjectResponse{}
	request, _ := json.Marshal(updateProjectRequest)
	err := performPost(context, api.httpclient, api.endpoint+"projects/"+id, api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Project, nil
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
func (api *Client) DeleteProjectByIdContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}

	err := performDelete(context,
		api.httpclient,
		api.endpoint+"projects/"+id,
		api.token,
		response,
		api)

	return response, err
}
