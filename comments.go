package golang_todoist_api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
)

type CommentsResponse struct {
	Comments []Comment
	TodoistResponse
}
type CommentResponse struct {
	Comment Comment
	TodoistResponse
}

type Comment struct {
	Content    string     `json:"content"`
	Id         string     `json:"id"`
	PostedAt   string     `json:"posted_at"`
	ProjectId  string     `json:"project_id"`
	TaskId     string     `json:"task_id"`
	Attachment Attachment `json:"attachment"`
}
type Attachment struct {
	ResourceType string `json:"resource_type"`
	FileUrl      string `json:"file_url"`
	FileType     string `json:"file_type"`
	FileName     string `json:"file_name"`
}

type NewCommentParameters struct {
	TaskId     string     `json:"task_id"`
	Content    string     `json:"content"`
	Attachment Attachment `json:"attachment"`
}

func (api *Client) GetAllCommentsByProjectId(projectId string) (*[]Comment, error) {
	return api.GetAllCommentsContext(projectId, "", context.Background())
}
func (api *Client) GetAllCommentsByTaskId(taskId string) (*[]Comment, error) {
	return api.GetAllCommentsContext("", taskId, context.Background())
}
func (api *Client) AddComment(param *NewCommentParameters) (*Comment, error) {
	return api.AddCommentContext(param, context.Background())
}
func (api *Client) UpdateComment(id string, content string) (*Comment, error) {
	return api.UpdateCommentContext(id, content, context.Background())
}
func (api *Client) GetCommentById(id string) (*Comment, error) {
	return api.GetCommentByIdContext(id, context.Background())
}
func (api *Client) DeleteCommentById(id string) (*TodoistResponse, error) {
	return api.DeleteCommentByIdContext(id, context.Background())
}

func (api *Client) GetAllCommentsContext(projectId string, taskId string, context context.Context) (*[]Comment, error) {
	response := &CommentsResponse{}
	values := url.Values{}
	if projectId != "" {
		values = url.Values{
			"project_id": {projectId},
		}
	} else if taskId != "" {
		values = url.Values{
			"task_id": {taskId},
		}
	} else {
		return nil, fmt.Errorf("task_id or project_id are missing on input")
	}

	err := api.getMethod(context,
		"comments",
		api.token,
		values,
		&response.Comments)

	return &response.Comments, err
}

func (api *Client) AddCommentContext(params *NewCommentParameters, context context.Context) (*Comment, error) {
	response := &CommentResponse{}
	request, _ := json.Marshal(params)
	err := performPost(context, api.httpclient, api.endpoint+"comments", api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Comment, nil
}

func (api *Client) UpdateCommentContext(id string, content string, context context.Context) (*Comment, error) {
	response := &CommentResponse{}
	request, _ := json.Marshal(map[string]string{
		"content": content,
	})
	err := performPost(context, api.httpclient, api.endpoint+"comments/"+id, api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Comment, nil
}
func (api *Client) GetCommentByIdContext(id string, context context.Context) (*Comment, error) {
	response := &CommentResponse{}

	err := api.getMethod(context,
		"comments/"+id,
		api.token,
		url.Values{},
		&response.Comment)

	return &response.Comment, err
}
func (api *Client) DeleteCommentByIdContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}

	err := performDelete(context,
		api.httpclient,
		api.endpoint+"comments/"+id,
		api.token,
		response,
		api)

	return response, err
}
