package hue_control_http

import (
	"context"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func HasError(res *http.Response, body *[]byte) bool {
	defer res.Body.Close()
	if res.StatusCode >= 300 || (body != nil && strings.Contains(string(*body), "error")) {
		log.Errorf("Error response: %v", string(*body))
		return true
	} else {
		return false
	}
}

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
