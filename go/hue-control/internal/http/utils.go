package hue_control_http

import (
	"net/http"
	"strings"
)

func HasError(res *http.Response, body *[]byte) bool {
	defer res.Body.Close()
	if res.StatusCode >= 300 || (body != nil && strings.Contains(string(*body), "error")) {
		return true
	} else {
		return false
	}
}
