package todoist

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetLabels(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/labels", getLabels)
	expectedProjects := getTestLabels()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	projects, err := api.GetLabels()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedProjects, *projects) {
		t.Fatal(ErrIncorrectResponse)
	}
}
func TestGetSharedLabels(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/labels/shared", getLabels)
	expectedProjects := getTestLabels()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	projects, err := api.GetSharedLabels()
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedProjects, *projects) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetLabelById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/labels/"+id, getLabelById(id))
	expectedProject := getTestLabelWithId(id)

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	project, err := api.GetLabelById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedProject, *project) {
		t.Fatal(ErrIncorrectResponse)
	}
}
func TestPostLabel(t *testing.T) {
	http.HandleFunc("/labels", postTestLabel)
	once.Do(startServer)
	expectedProject := getTestLabelWithId("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	request1 := LabelRequest{
		Name:       "name",
		Color:      "",
		Order:      0,
		IsFavorite: false,
	}
	project, err := api.AddLabel(request1)
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
func TestPostLabelByName(t *testing.T) {
	http.HandleFunc("/labels/shared/rename", postLabelNewName())
	once.Do(startServer)
	expectedProject := getTestLabelByName()

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	project, err := api.RenameLabel("name0", "name1")
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
func TestPostSharedLabelByName(t *testing.T) {
	http.HandleFunc("/labels/shared/remove", postLabelNewName())
	once.Do(startServer)
	expectedProject := getTestLabelByName()

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	project, err := api.RemoveSharedLabel("name0")
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

func TestPostLabelById(t *testing.T) {
	http.HandleFunc("/projects/1", postTestProjectById("1"))
	once.Do(startServer)
	expectedProject := getTestProjectWithId("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	request := UpdateProjectRequest{
		Name: "name",
	}
	project, err := api.UpdateProject("1", request)
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

func TestDeleteLabelById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/projects/"+id, getDeleteLabelByIdResponse)
	expectedResponse := getTestDeleteLabelByIdResponse()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.DeleteProjectById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func getTestLabels() []Label {
	return []Label{
		getTestLabelWithId("12345"),
		getTestLabelWithId("23456"),
	}
}

func getTestLabelWithId(id string) Label {
	return Label{
		ID:         id,
		Name:       "",
		Color:      "",
		Order:      0,
		IsFavorite: false,
	}
}

func getTestLabelByName() TodoistResponse {
	return TodoistResponse{
		Ok:    true,
		Error: "",
	}
}

func getTestDeleteLabelByIdResponse() TodoistResponse {
	return TodoistResponse{
		Ok:    true,
		Error: "",
	}
}

func getLabels(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(
		getTestProjects(),
	)
	_, err := rw.Write(response)
	if err != nil {
		return
	}
}
func getLabelById(id string) func(rw http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(
		getTestLabelWithId(id),
	)

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}
func postLabelNewName() func(w http.ResponseWriter, r *http.Request) {

	response, _ := json.Marshal(TodoistResponse{
		Ok:    true,
		Error: "",
	})

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}

func postTestLabel(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(getTestLabelWithId("1"))
	_, err := w.Write(response)
	if err != nil {
		return
	}
}

func getDeleteLabelByIdResponse(rw http.ResponseWriter, _ *http.Request) {
	response, _ := json.Marshal(TodoistResponse{
		Ok:    true,
		Error: "",
	})

	_, err := rw.Write(response)
	if err != nil {
		return
	}
}
