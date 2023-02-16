Todoist API in Go
===============
This library supports most if not all of the `developer.todoist.com/rest/v2` REST
calls.

Todoist API in Go [![Go Reference](https://pkg.go.dev/badge/github.com/volyanyk/todoist.svg)](https://pkg.go.dev/github.com/volyanyk/todoist)
===============

This is not original Todoist library for Go.

This library supports most if not all of the `developer.todoist.com/rest/v2` REST
calls.

## Project Status
There is currently no major version released.


## Installing

### *go get*

    $ go get -u github.com/volyanyk/todoist

## Example

### Getting all projects for the authenticated user

```golang
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

```
## Contributing

You are more than welcome to contribute to this project.  Fork and
make a Pull Request, or create an Issue if you see any problem.

Before making any Pull Request please run the following:

```
make pr-prep
```

This will check/update code formatting, linting and then run all tests

## License

BSD 2 Clause license