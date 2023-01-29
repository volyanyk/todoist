package golang_todoist_api

import "net/http"

const (
	APIURL = "https://api.todoist.com/rest/v2/"
)

type Client struct {
	token         string
	appLevelToken string
	endpoint      string
	debug         bool
	log           ilogger
	httpclient    httpClient
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (api *Client) Debug() bool {
	return api.debug
}
