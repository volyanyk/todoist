package golang_todoist_api

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type responseParser func(*http.Response) error

type TodoistErrorResponse struct {
	Err string
}

func (t TodoistErrorResponse) Error() string { return t.Err }

func checkStatusCode(resp *http.Response, d Debug) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitedError{time.Duration(retry) * time.Second}
	}

	if resp.StatusCode != http.StatusOK {
		err := logResponse(resp, d)
		if err != nil {
			return err
		}
		return StatusCodeError{Code: resp.StatusCode, Status: resp.Status}
	}

	return nil
}

func doPost(ctx context.Context, client httpClient, req *http.Request, parser responseParser, d Debug) error {
	resp, err := client.Do(req)

	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(resp.Body)

	err = checkStatusCode(resp, d)

	return parser(resp)
}

func postJSON(ctx context.Context, client httpClient, endpoint, token string, json []byte, intf interface{}, d Debug) error {
	reqBody := bytes.NewBuffer(json)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, reqBody)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Set("X-Request-ID", uuid.New().String())

	return doPost(ctx, client, req, newJSONParser(intf), d)
}

func logResponse(resp *http.Response, d Debug) error {
	if d.Debug() {
		text, err := httputil.DumpResponse(resp, true)
		if err != nil {
			return err
		}
		d.Debugln(string(text))
	}

	return nil
}

type RateLimitedError struct {
	RetryAfter time.Duration
}

func (e *RateLimitedError) Error() string {
	return fmt.Sprintf("rate limit exceeded, retry after %s", e.RetryAfter)
}

func (e *RateLimitedError) Retryable() bool {
	return true
}

func (t TodoistResponse) Err() error {
	if t.Ok {
		return nil
	}

	if strings.TrimSpace(t.Error) == "" {
		return nil
	}

	return TodoistErrorResponse{Err: t.Error}
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
