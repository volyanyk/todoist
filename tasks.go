package todoist

import (
	"context"
	"encoding/json"
	"net/url"
	"strings"
)

type TasksResponse struct {
	Tasks []Task
	TodoistResponse
}
type TaskResponse struct {
	Task Task
	TodoistResponse
}
type Due struct {
	Date        string `json:"date"`
	IsRecurring bool   `json:"is_recurring"`
	Datetime    string `json:"datetime"`
	String      string `json:"string"`
	Timezone    string `json:"timezone"`
}

type Task struct {
	CreatorId    string   `json:"creator_id"`
	CreatedAt    string   `json:"created_at"`
	AssigneeId   string   `json:"assignee_id"`
	AssignerId   string   `json:"assigner_id"`
	CommentCount int      `json:"comment_count"`
	IsCompleted  bool     `json:"is_completed"`
	Content      string   `json:"content"`
	Description  string   `json:"description"`
	Due          Due      `json:"due"`
	Id           string   `json:"id"`
	Labels       []string `json:"labels"`
	Order        int      `json:"order"`
	Priority     int      `json:"priority"`
	ProjectId    string   `json:"project_id"`
	SectionId    string   `json:"section_id"`
	ParentId     string   `json:"parent_id"`
	Url          string   `json:"url"`
}
type AddTaskRequest struct {
	Content     string   `json:"content"`
	Description string   `json:"description"`
	ProjectId   string   `json:"project_id"`
	SectionId   string   `json:"section_id"`
	ParentId    string   `json:"parent_id"`
	Labels      []string `json:"labels"`
	Order       int      `json:"order"`
	Priority    int      `json:"priority"`
	AssigneeId  string   `json:"assignee_id"`
	DueString   string   `json:"due_string"`
	DueDate     string   `json:"due_date"`
	DueDatetime string   `json:"due_datetime"`
	DueLang     string   `json:"due_lang"`
}

type GetActiveTasksRequest struct {
	ProjectId string   `json:"project_id"` // Optional
	SectionId string   `json:"section_id"` // Optional
	Label     string   `json:"label"`      // Optional
	Filter    string   `json:"filter"`     // Optional
	Lang      string   `json:"lang"`       // Optional
	Ids       []string `json:"ids"`        // Optional
}
type UpdateTaskRequest struct {
	Content     string   `json:"content"`      // Optional
	Description string   `json:"description"`  // Optional
	Labels      []string `json:"labels"`       // Optional
	Priority    int      `json:"priority"`     // Optional
	DueString   string   `json:"due_string"`   // Optional
	DueDate     string   `json:"due_date"`     // Optional
	DueDatetime string   `json:"due_datetime"` // Optional
	DueLang     string   `json:"due_lang"`     // Optional
	AssigneeId  string   `json:"assignee_id"`  // Optional
}

func (api *Client) GetActiveTasks(getActiveTasksRequest GetActiveTasksRequest) (*[]Task, error) {
	return api.GetActiveTasksContext(getActiveTasksRequest, context.Background())
}
func (api *Client) AddTask(request AddTaskRequest) (*Task, error) {
	return api.AddTaskContext(request, context.Background())
}
func (api *Client) GetActiveTaskById(id string) (*Task, error) {
	return api.GetActiveTaskByIdContext(id, context.Background())
}
func (api *Client) UpdateTask(id string, updateTaskRequest UpdateTaskRequest) (*Task, error) {
	return api.UpdateTaskContext(id, updateTaskRequest, context.Background())
}
func (api *Client) CloseTask(id string) (*TodoistResponse, error) {
	return api.CloseTaskContext(id, context.Background())
}
func (api *Client) ReopenTask(id string) (*TodoistResponse, error) {
	return api.ReopenTaskContext(id, context.Background())
}
func (api *Client) DeleteTaskById(id string) (*TodoistResponse, error) {
	return api.DeleteTaskByIdContext(id, context.Background())
}

func (api *Client) GetActiveTasksContext(request GetActiveTasksRequest, context context.Context) (*[]Task, error) {
	response := &TasksResponse{}
	values := url.Values{
		"project_id": {request.ProjectId},
		"section_id": {request.SectionId},
		"label":      {request.Label},
		"filter":     {request.Filter},
		"lang":       {request.Lang},
	}
	if len(request.Ids) > 0 {
		values.Add("ids", strings.Join(request.Ids, ","))
	}

	err := api.getMethod(context,
		"tasks",
		api.token,
		values,
		&response.Tasks)

	return &response.Tasks, err
}

func (api *Client) AddTaskContext(addTaskRequest AddTaskRequest, context context.Context) (*Task, error) {
	response := &TaskResponse{}

	request, err := json.Marshal(addTaskRequest)

	err = performPost(context, api.httpclient, api.endpoint+"tasks", api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Task, nil
}
func (api *Client) GetActiveTaskByIdContext(id string, context context.Context) (*Task, error) {
	response := &TaskResponse{}

	err := api.getMethod(context,
		"tasks/"+id,
		api.token,
		url.Values{},
		&response.Task)

	return &response.Task, err
}
func (api *Client) UpdateTaskContext(id string, updateTaskRequest UpdateTaskRequest, context context.Context) (*Task, error) {
	response := &TaskResponse{}
	request, _ := json.Marshal(updateTaskRequest)
	err := performPost(context, api.httpclient, api.endpoint+"tasks/"+id, api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Task, nil
}
func (api *Client) CloseTaskContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}
	err := performPostWithoutResponse(context, api.httpclient, api.endpoint+"tasks/"+id+"/close", api.token, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return response, err
}
func (api *Client) ReopenTaskContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}
	err := performPostWithoutResponse(context, api.httpclient, api.endpoint+"tasks/"+id+"/reopen", api.token, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return response, err
}
func (api *Client) DeleteTaskByIdContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}

	err := performDelete(context,
		api.httpclient,
		api.endpoint+"tasks/"+id,
		api.token,
		response,
		api)

	return response, err
}
