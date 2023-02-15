package todoist

import (
	"context"
	"encoding/json"
	"net/url"
)

type LabelsResponse struct {
	Labels []Label
	TodoistResponse
}
type LabelResponse struct {
	Label Label
	TodoistResponse
}

type Label struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Color      string `json:"color"`
	Order      int    `json:"order"`
	IsFavorite bool   `json:"is_favorite"`
}

type LabelRequest struct {
	Name       string `json:"name"`        // Required / Optional
	Color      string `json:"color"`       // Optional
	Order      int    `json:"order"`       // Optional
	IsFavorite bool   `json:"is_favorite"` // Optional
}

func (api *Client) GetLabels() (*[]Label, error) {
	return api.GetLabelsContext(context.Background())
}
func (api *Client) GetSharedLabels() (*[]Label, error) {
	return api.GetSharedLabelsContext(context.Background())
}
func (api *Client) AddLabel(request LabelRequest) (*Label, error) {
	return api.AddLabelContext(request, context.Background())
}
func (api *Client) UpdateLabel(id string, request LabelRequest) (*Label, error) {
	return api.UpdateLabelContext(id, request, context.Background())
}
func (api *Client) RenameLabel(oldName string, newName string) (*TodoistResponse, error) {
	return api.RenameLabelContext(oldName, newName, context.Background())
}
func (api *Client) RemoveSharedLabel(name string) (*TodoistResponse, error) {
	return api.RemoveSharedLabelContext(name, context.Background())
}
func (api *Client) GetLabelById(id string) (*Label, error) {
	return api.GetLabelByIdContext(id, context.Background())
}
func (api *Client) DeleteLabelById(id string) (*TodoistResponse, error) {
	return api.DeleteLabelByIdContext(id, context.Background())
}

func (api *Client) GetLabelsContext(context context.Context) (*[]Label, error) {
	response := &LabelsResponse{}

	err := api.getMethod(context,
		"labels",
		api.token,
		url.Values{},
		&response.Labels)

	return &response.Labels, err
}
func (api *Client) GetSharedLabelsContext(context context.Context) (*[]Label, error) {
	response := &LabelsResponse{}

	err := api.getMethod(context,
		"labels/shared",
		api.token,
		url.Values{},
		&response.Labels)

	return &response.Labels, err
}

func (api *Client) AddLabelContext(addLabelRequest LabelRequest, context context.Context) (*Label, error) {
	response := &LabelResponse{}

	request, err := json.Marshal(addLabelRequest)

	err = performPost(context, api.httpclient, api.endpoint+"labels", api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Label, nil
}

func (api *Client) UpdateLabelContext(id string, updateLabelRequest LabelRequest, context context.Context) (*Label, error) {
	response := &LabelResponse{}
	request, _ := json.Marshal(updateLabelRequest)
	err := performPost(context, api.httpclient, api.endpoint+"labels/"+id, api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return &response.Label, nil
}
func (api *Client) RenameLabelContext(oldName string, newName string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}
	request, _ := json.Marshal(map[string]string{
		"name":     oldName,
		"new_name": newName,
	})
	err := performPost(context, api.httpclient, api.endpoint+"labels/shared/rename", api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return response, nil
}
func (api *Client) RemoveSharedLabelContext(name string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}
	request, _ := json.Marshal(map[string]string{
		"name": name,
	})
	err := performPost(context, api.httpclient, api.endpoint+"labels/shared/remove", api.token, request, &response, api)

	if err != nil {
		return nil, err
	}
	if !response.Ok {
		return nil, response.Err()
	}

	return response, nil
}
func (api *Client) GetLabelByIdContext(id string, context context.Context) (*Label, error) {
	response := &LabelResponse{}

	err := api.getMethod(context,
		"labels/"+id,
		api.token,
		url.Values{},
		&response.Label)

	return &response.Label, err
}
func (api *Client) DeleteLabelByIdContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}

	err := performDelete(context,
		api.httpclient,
		api.endpoint+"labels/"+id,
		api.token,
		response,
		api)

	return response, err
}
