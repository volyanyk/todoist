package golang_todoist_api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"strconv"
	"time"
)

type responseParser func(*http.Response) error

func checkStatusCode(resp *http.Response, d Debug) error {
	if resp.StatusCode == http.StatusTooManyRequests {
		retry, err := strconv.ParseInt(resp.Header.Get("Retry-After"), 10, 64)
		if err != nil {
			return err
		}
		return &RateLimitedError{time.Duration(retry) * time.Second}
	}

	if resp.StatusCode != http.StatusOK {
		logResponse(resp, d)
		return StatusCodeError{Code: resp.StatusCode, Status: resp.Status}
	}

	return nil
}

func doPost(ctx context.Context, client httpClient, req *http.Request, parser responseParser, d Debug) error {
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = checkStatusCode(resp, d)
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
