package main

import (
	"fmt"

	"github.com/volyanyk/todoist"
)

func main() {
	api := todoist.New("TOKEN")

	tasks, err := api.GetActiveTasks(todoist.GetActiveTasksRequest{
		ProjectId: "2171745492",
		Filter:    "(today|overdue)",
	})
	if err != nil {
		return
	}

	fmt.Printf("%v", tasks)

}
