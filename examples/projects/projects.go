package main

import (
	"fmt"

	"github.com/volyanyk/todoist"
)

func main() {
	api := todoist.New("TOKEN")
	projects, err := api.GetProjects()
	if err != nil {
		return
	}

	fmt.Printf("%v", projects)

}
