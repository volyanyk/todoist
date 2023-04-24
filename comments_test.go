package todoist

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

func TestGetAllCommentsByProjectId(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/comments", getCommentsForProjectId)
	expectedComments := getTestCommentsByProjectId("1")

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	sections, err := api.GetAllCommentsByProjectId("1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedComments, *sections) {
		t.Fatal(ErrIncorrectResponse)
	}
}
func TestGetAllCommentsByTaskId(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/comments", getCommentsForTaskId)
	expectedComments := getTestCommentsByTaskId("1")

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	sections, err := api.GetAllCommentsByTaskId("1")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedComments, *sections) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestGetCommentById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/comments/"+id, getCommentById(id))
	expectedComment := getTestCommentById(id)

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	comment, err := api.GetCommentById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedComment, *comment) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func TestAddComment(t *testing.T) {
	http.HandleFunc("/comments", addComment)
	once.Do(startServer)
	expectedProject := getTestCommentById("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))
	param := NewCommentParameters{
		TaskId:  "",
		Content: "",
		Attachment: Attachment{
			ResourceType: "",
			FileUrl:      "",
			FileType:     "",
			FileName:     "",
		},
	}
	project, err := api.AddComment(&param)
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

func TestUpdateComment(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	http.HandleFunc("/comments/1", updateCommentById("1"))
	once.Do(startServer)
	expectedComment := getTestCommentById("1")

	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	comment, err := api.UpdateComment("1", "name")
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedComment, *comment) {
		t.Fatal(ErrIncorrectResponse)
	}
	if err != nil {
		t.Errorf("Failed, but should have succeeded")
	}
}

func TestDeleteCommentById(t *testing.T) {
	http.DefaultServeMux = new(http.ServeMux)
	var id = "1"
	http.HandleFunc("/comments/"+id, getDeleteCommentByIdResponse)
	expectedResponse := getTestDeleteCommentByIdResponse()

	once.Do(startServer)
	api := New(validToken, OptionAPIURL("http://"+serverAddr+"/"))

	response, err := api.DeleteCommentById(id)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
		return
	}
	if !reflect.DeepEqual(expectedResponse, *response) {
		t.Fatal(ErrIncorrectResponse)
	}
}

func getTestDeleteCommentByIdResponse() TodoistResponse {
	return TodoistResponse{
		Ok:    true,
		Error: "",
	}
}

func getDeleteCommentByIdResponse(writer http.ResponseWriter, _ *http.Request) {
	response, _ := json.Marshal(TodoistResponse{
		Ok:    true,
		Error: "",
	})

	_, err := writer.Write(response)
	if err != nil {
		return
	}
}

func getCommentsForProjectId(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(
		getTestCommentsByProjectId("1"),
	)

	_, err := rw.Write(response)
	if err != nil {
		return
	}

}
func getCommentsForTaskId(rw http.ResponseWriter, _ *http.Request) {
	rw.Header().Set("Content-Type", "application/json")

	response, _ := json.Marshal(
		getTestCommentsByTaskId("1"),
	)

	_, err := rw.Write(response)
	if err != nil {
		return
	}

}

func getTestCommentsByProjectId(projectId string) []Comment {
	return []Comment{
		getTestCommentByProjectId(projectId),
		getTestCommentByProjectId(projectId),
	}
}
func getTestCommentsByTaskId(taskId string) []Comment {
	return []Comment{
		getTestCommentByTaskId(taskId),
		getTestCommentByTaskId(taskId),
	}
}
func getTestCommentByProjectId(projectId string) Comment {
	return Comment{
		Content:    "",
		Id:         "",
		PostedAt:   "",
		ProjectId:  projectId,
		TaskId:     "",
		Attachment: nil,
	}
}
func getTestCommentByTaskId(taskId string) Comment {
	return Comment{
		Content:    "",
		Id:         "",
		PostedAt:   "",
		ProjectId:  "",
		TaskId:     taskId,
		Attachment: nil,
	}
}

func getCommentById(id string) func(http.ResponseWriter, *http.Request) {
	response, _ := json.Marshal(
		getTestCommentById(id),
	)

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}

func getTestCommentById(id string) Comment {
	return Comment{
		Content:    "",
		Id:         id,
		PostedAt:   "",
		ProjectId:  "",
		TaskId:     "",
		Attachment: nil,
	}
}

func addComment(writer http.ResponseWriter, _ *http.Request) {
	writer.Header().Set("Content-Type", "application/json")
	response, _ := json.Marshal(getTestCommentById("1"))
	_, err := writer.Write(response)
	if err != nil {
		return
	}
}

func updateCommentById(id string) func(http.ResponseWriter, *http.Request) {
	response, _ := json.Marshal(getTestCommentById(id))

	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")
		_, err := rw.Write(response)
		if err != nil {
			return
		}
	}
}
