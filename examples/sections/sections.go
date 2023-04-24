package main

import (
	"fmt"

	"github.com/volyanyk/todoist"
)

func main() {
	api := todoist.New("TOKEN")
	section := todoist.SectionParameters{
		ProjectId: "2234026034",
		Name:      "Test1",
		Order:     0,
	}
	sections, err := api.AddSection(&section)
	//sections, err := api.UpdateSection("121915105", "test")
	if err != nil {
		return
	}

	fmt.Printf("%v", sections)

}
