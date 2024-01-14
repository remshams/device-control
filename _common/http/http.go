package dc_http

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/charmbracelet/log"
)

func RequestWithTimeout(method string, url string, body io.Reader, timeout *time.Duration) (*http.Request, *http.Client, context.CancelFunc, error) {
	defaultTimeout := 2 * time.Second
	requestTimeout := timeout
	if requestTimeout == nil {
		requestTimeout = &defaultTimeout
	}
	ctx, cancel := context.WithTimeout(context.Background(), *requestTimeout)
	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	return req, client, cancel, err
}

func PerformRequest(context string, path string, method string, body io.Reader, timeout *time.Duration) (*http.Response, error) {
	log.Debugf("%s: Performing request with path: %s and method: %s", context, path, method)
	req, client, cancel, err := RequestWithTimeout(
		method,
		path,
		body,
		timeout,
	)
	defer cancel()
	res, err := client.Do(req)
	if err != nil {
		log.Error("%s: Could nor perform request", context)
		return nil, err
	}
	if res.StatusCode != 200 {
		log.Error("%s: Request failed with status code %d", context, res.StatusCode)
		return nil, err
	}
	return res, nil
}
