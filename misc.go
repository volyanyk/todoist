package todoist

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"time"
)

type responseParser func(*http.Response) error

type ErrorResponse struct {
	Err string
}

func (t ErrorResponse) Error() string { return t.Err }

func checkStatusCode(method string, resp *http.Response, d Debug) error {
	if resp.StatusCode == http.StatusUnauthorized {
		return StatusCodeError{Code: resp.StatusCode, Status: resp.Status}
	}
	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitedError{time.Duration(retry) * time.Second}
	}
	if method == http.MethodDelete && resp.StatusCode == http.StatusOK {
		resp.Status = http.StatusText(http.StatusNoContent)
		resp.StatusCode = http.StatusNoContent
		return nil
	}
	if method == http.MethodPost && ((resp.StatusCode != http.StatusOK) && (resp.StatusCode != http.StatusNoContent)) {
		err := logResponse(resp, d)
		if err != nil {
			return err
		}
		return StatusCodeError{Code: resp.StatusCode, Status: resp.Status}

	} else {
		if resp.StatusCode != http.StatusOK {
			err := logResponse(resp, d)
			if err != nil {
				return err
			}
			return StatusCodeError{Code: resp.StatusCode, Status: resp.Status}
		}
	}

	return nil
}

func perform(client httpClient, req *http.Request, parser responseParser, d Debug) error {
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

	err = checkStatusCode(req.Method, resp, d)
	if err != nil {
		return err
	}

	return parser(resp)
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

	return ErrorResponse{Err: t.Error}
}

func New(token string, options ...Option) *Client {
	s := &Client{
		token:      token,
		endpoint:   APIURL,
		httpclient: &http.Client{},
		log:        log.New(os.Stderr, "volyanyk/todoist", log.LstdFlags|log.Lshortfile),
	}

	for _, opt := range options {
		opt(s)
	}

	return s
}

func OptionAPIURL(u string) func(*Client) {
	return func(c *Client) { c.endpoint = u }
}

type Option func(*Client)
