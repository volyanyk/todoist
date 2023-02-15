package main

import (
	"fmt"

	todoist "github.com/volyanyk/todoist"
)

func main() {
	api := todoist.New("TOKEN")
	comments, err := api.GetAllCommentsByProjectId("2171745488")
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", comments)
}
