package main

import (
	"fmt"

	"github.com/volyanyk/todoist"
)

func main() {
	api := todoist.New("TOKEN")
	comments, err := api.GetAllCommentsByTaskId("6590043544")
	if err != nil {
		fmt.Printf("%v", err)
	}
	fmt.Printf("%v", comments)
}
