package todoist

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetActiveTasks(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/tasks", getTasks)
	expectedTasks := getTestTasks()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	request := GetActiveTasksRequest{
		ProjectId: "1",
		SectionId: "1",
		Label:     "1",
		Filter:    "1",
		Lang:      "2",
		Ids:       []string{"1", "2"},
	}
	tasks, err := api.GetActiveTasks(request)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedTasks, *tasks) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetActiveTaskById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/tasks/"+id, getTaskById(id))
	expectedTask := getTestTaskWithId(id)

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	task, err := api.GetActiveTaskById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedTask, *task) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestAddTask(t *testing.T) {
	http.HandleFunc("/tasks", postTestTask)
	once.Do(startServer)
	expectedTask := getTestTaskWithId("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	sectionID := "4"
	parentID := "5"
	assigneeID := "6"
	request1 := AddTaskRequest{
		Content:     "1",
		Description: "2",
		ProjectId:   "3",
		SectionId:   &sectionID,
		ParentId:    &parentID,
		Labels:      []string{"1", "2"},
		Order:       1,
		Priority:    2,
		AssigneeId:  &assigneeID,
		DueString:   "7",
		DueDate:     "8",
		DueDatetime: "9",
		DueLang:     "0",
	}
	task, err := api.AddTask(request1)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedTask, *task) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func TestUpdateTask(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/tasks/1", addTestTaskById("1"))
	once.Do(startServer)
	expectedTask := getTestTaskWithId("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	request := UpdateTaskRequest{
		Content:     "",
		Description: "",
	}
	task, err := api.UpdateTask("1", request)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedTask, *task) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func TestDeleteTaskById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/tasks/"+id, getOkResponse)
	expectedResponse := getTestOkResponse()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.DeleteTaskById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestCloseTask(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/tasks/1/close", getOkResponse)
	once.Do(startServer)
	expectedResponse := getTestOkResponse()

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.CloseTask("1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}
func TestReopenTask(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/tasks/1/reopen", getOkResponse)
	once.Do(startServer)
	expectedResponse := getTestOkResponse()

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	response, err := api.ReopenTask("1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func getTestTasks() []Task {
	return []Task{
		getTestTaskWithId("12345"),
		getTestTaskWithId("23456"),
	}
}

func getTestTaskWithId(id string) Task {
	return Task{
		CreatorId:    "",
		CreatedAt:    "",
		AssigneeId:   nil,
		AssignerId:   nil,
		CommentCount: 0,
		IsCompleted:  false,
		Content:      "",
		Description:  "",
		Due:          nil,
		Id:           id,
		Labels:       nil,
		Order:        0,
		Priority:     0,
		ProjectId:    "",
		SectionId:    "",
		ParentId:     "",
		Url:          "",
	}
}

func getTestOkResponse() TodoistResponse {
	return TodoistResponse{
		Ok:    true,
		Error: "",
	}
}

func getTasks(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(
		getTestTasks(),
	)
	_, err := rw.Write(response)
	if err != nil {
		return
	}
}
func getTaskById(id string) func(rw http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(
		getTestTaskWithId(id),
	)

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}

func postTestTask(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(getTestTaskWithId("1"))
	_, err := w.Write(response)
	if err != nil {
		return
	}
}

func addTestTaskById(id string) func(rw http.ResponseWriter, r *http.Request) {
	response, _ := json.Marshal(getTestTaskWithId(id))

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}
func getOkResponse(rw http.ResponseWriter, _ *http.Request) {
	response, _ := json.Marshal(TodoistResponse{
		Ok:    true,
		Error: "",
	})

	_, err := rw.Write(response)
	if err != nil {
		return
	}
}
