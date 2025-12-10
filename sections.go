package todoist

import (
	"context"
	"encoding/json"
	"net/url"
)

type SectionsResponse struct {
	Sections []Section
	TodoistResponse
}
type SectionResponse struct {
	Section Section
	TodoistResponse
}

type Section struct {
	ID        string `json:"id"`
	ProjectId string `json:"project_id"`
	Order     *int   `json:"order"`
	Name      string `json:"name"`
}

type SectionParameters struct {
	ProjectId string `json:"project_id"`
	Name      string `json:"name"`
	Order     *int   `json:"order"`
}

func (api *Client) GetSectionsByProjectId(projectId string) (*[]Section, error) {
	return api.GetSectionsByProjectIdContext(projectId, context.Background())
}
func (api *Client) AddSection(param *SectionParameters) (*Section, error) {
	return api.AddSectionContext(param, context.Background())
}
func (api *Client) UpdateSection(projectId string, name string) (*Section, error) {
	return api.UpdateSectionContext(projectId, name, context.Background())
}
func (api *Client) GetSectionById(id string) (*Section, error) {
	return api.GetSectionByIdContext(id, context.Background())
}
func (api *Client) DeleteSectionById(id string) (*TodoistResponse, error) {
	return api.DeleteSectionByIdContext(id, context.Background())
}

func (api *Client) GetSectionsByProjectIdContext(projectId string, context context.Context) (*[]Section, error) {
	response := &SectionsResponse{}
	values := url.Values{
		"project_id": {projectId},
	}
	err := api.get(context,
		"sections",
		api.token,
		values,
		&response.Sections)

	return &response.Sections, err
}

func (api *Client) GetSectionCollaboratorsContext(id string, context context.Context) (*[]Collaborator, error) {
	response := &CollaboratorsResponse{}

	err := api.get(context,
		"sections/"+id+"/collaborators",
		api.token,
		url.Values{},
		&response.Collaborators)

	return &response.Collaborators, err
}

func (api *Client) AddSectionContext(params *SectionParameters, context context.Context) (*Section, error) {
	response := &SectionResponse{}
	request, _ := json.Marshal(params)
	err := api.post(context, "sections", api.token, request, &response.Section)

	if err != nil {
		return nil, err
	} else {
		return &response.Section, nil
	}

}

func (api *Client) UpdateSectionContext(sectionId string, name string, context context.Context) (*Section, error) {
	response := &SectionResponse{}
	request, _ := json.Marshal(map[string]string{
		"name": name,
	})
	err := api.post(context, "sections/"+sectionId, api.token, request, &response.Section)

	if err != nil {
		return nil, err
	} else {
		return &response.Section, nil
	}

}
func (api *Client) GetSectionByIdContext(id string, context context.Context) (*Section, error) {
	response := &SectionResponse{}

	err := api.get(context,
		"sections/"+id,
		api.token,
		url.Values{},
		&response.Section)

	return &response.Section, err
}
func (api *Client) DeleteSectionByIdContext(id string, context context.Context) (*TodoistResponse, error) {
	response := &TodoistResponse{}

	err := performDelete(context,
		api.httpclient,
		api.endpoint+"sections/"+id,
		api.token,
		response,
		api)

	return response, err
}
