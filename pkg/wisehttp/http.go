package wisehttp

import "net/http"

// An http client interface used for sending requests.
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

var (
	_ HTTPClient = (*http.Client)(nil)
)
