package golang_todoist_api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

const (
	APIURL = "https://api.todoist.com/rest/v2/"
)

type Client struct {
	token      string
	endpoint   string
	debug      bool
	log        ilogger
	httpclient httpClient
}

type TodoistResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

func (api *Client) Debugf(format string, v ...interface{}) {
	if api.debug {
		err := api.log.Output(2, fmt.Sprintf(format, v...))
		if err != nil {
			return
		}
	}
}
func (api *Client) Debugln(v ...interface{}) {
	if api.debug {
		err := api.log.Output(2, fmt.Sprintln(v...))
		if err != nil {
			return
		}
	}
}

func (api *Client) Debug() bool {
	return api.debug
}

type httpClient interface {
	Do(*http.Request) (*http.Response, error)
}

func (api *Client) getMethod(ctx context.Context, path string, token string, values url.Values, intf interface{}) error {
	return performGet(ctx, api.httpclient, api.endpoint+path, token, values, intf, api)
}

func performPost(ctx context.Context, client httpClient, endpoint, token string, json []byte, intf interface{}, d Debug) error {
	reqBody := bytes.NewBuffer(json)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-Request-ID", uuid.New().String())
	req.Header.Set("Content-Type", "application/json")

	return perform(ctx, client, req, newJSONParser(intf), d)
}
func performPostWithoutResponse(ctx context.Context, client httpClient, endpoint, token string, intf interface{}, d Debug) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	return perform(ctx, client, req, newJSONParser(intf), d)
}

func performGet(ctx context.Context, client httpClient, endpoint, token string, values url.Values, intf interface{}, d Debug) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	req.URL.RawQuery = values.Encode()

	return perform(ctx, client, req, newJSONParser(intf), d)
}
func performDelete(ctx context.Context, client httpClient, endpoint, token string, intf interface{}, d Debug) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	return perform(ctx, client, req, newJSONParser(intf), d)
}

func newJSONParser(dst interface{}) responseParser {
	return func(resp *http.Response) error {
		if dst == nil {
			return nil
		}
		return json.NewDecoder(resp.Body).Decode(dst)
	}
}
