package todoist

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetSectionsByProjectId(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/sections", getSectionsByProjectId)
	expectedSections := getTestSectionsByProjectId("1")

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	sections, err := api.GetSectionsByProjectId("1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedSections, *sections) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetSectionById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/sections/"+id, getSectionById(id))
	expectedSection := getTestSectionById(id)

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	section, err := api.GetSectionById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedSection, *section) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestAddSection(t *testing.T) {
	http.HandleFunc("/sections", addSection)
	once.Do(startServer)
	expectedProject := getTestSectionByProjectIdAndName("0", "1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	param := SectionParameters{
		ProjectId: "0",
		Name:      "1",
	}
	project, err := api.AddSection(&param)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedProject, *project) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func TestUpdateSection(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/sections/1", updateSectionById("1"))
	once.Do(startServer)
	expectedSection := getTestSectionByProjectId("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	section, err := api.UpdateSection("1", "name")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedSection, *section) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func TestDeleteSectionById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/sections/"+id, getDeleteSectionByIdResponse)
	expectedResponse := getTestDeleteSectionByIdResponse()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.DeleteSectionById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func getTestDeleteSectionByIdResponse() TodoistResponse {
	return TodoistResponse{
		Ok:    true,
		Error: "",
	}
}

func getDeleteSectionByIdResponse(writer http.ResponseWriter, _ *http.Request) {
	response, _ := json.Marshal(TodoistResponse{
		Ok:    true,
		Error: "",
	})

	_, err := writer.Write(response)
	if err != nil {
		return
	}
}

func getSectionsByProjectId(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(
		getTestSectionsByProjectId("1"),
	)

	_, err := rw.Write(response)
	if err != nil {
		return
	}

}

func getTestSectionsByProjectId(projectId string) []Section {
	return []Section{
		getTestSectionByProjectId(projectId),
		getTestSectionByProjectId(projectId),
	}
}
func getTestSectionByProjectId(projectId string) Section {
	return getTestSectionByProjectIdAndName(projectId, "name")
}
func getTestSectionByProjectIdAndName(projectId string, name string) Section {
	return Section{
		ID:        "",
		ProjectId: projectId,
		Order:     0,
		Name:      name,
	}
}

func getSectionById(id string) func(http.ResponseWriter, *http.Request) {
	response, _ := json.Marshal(
		getTestSectionById(id),
	)

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}

func getTestSectionById(id string) Section {
	return Section{
		ID:        id,
		ProjectId: "",
		Order:     0,
		Name:      "",
	}
}

func addSection(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(SectionResponse{
		Section:         getTestSectionByProjectIdAndName("0", "1"),
		TodoistResponse: TodoistResponse{Ok: true},
	})
	_, err := writer.Write(response)
	if err != nil {
		return
	}
}

func updateSectionById(id string) func(http.ResponseWriter, *http.Request) {
	response, _ := json.Marshal(SectionResponse{
		Section:         getTestSectionByProjectId(id),
		TodoistResponse: TodoistResponse{Ok: true},
	})

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}
