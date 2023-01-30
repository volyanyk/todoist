package golang_todoist_api

import (
	"log"
	"net/http"
	"os"
)

type Option func(*Client)

func OptionAPIURL(u string) func(*Client) {
	return func(c *Client) { c.endpoint = u }
}

type Project struct {
	ID             string `json:"id"`
	ParentId       string `json:"parent_id"`
	Order          int    `json:"order"`
	Color          string `json:"color"`
	Name           string `json:"name"`
	CommentCount   int    `json:"comment_count"`
	IsShared       bool   `json:"is_shared"`
	IsFavorite     bool   `json:"is_favorite"`
	IsInboxProject bool   `json:"is_inbox_project"`
	IsTeamInbox    bool   `json:"is_team_inbox"`
	Url            string `json:"url"`
	ViewStyle      string `json:"view_style"`
}

func New(token string, options ...Option) *Client {
	s := &Client{
		token:      token,
		endpoint:   APIURL,
		httpclient: &http.Client{},
		log:        log.New(os.Stderr, "volyanyk/golang-todoist-api", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}
