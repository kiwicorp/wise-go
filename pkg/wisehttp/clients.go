package wisehttp

import "net/http"

// An http client that automatically attaches an user agent.
type WithUA struct {
	c  HTTPClient
	ua string
}

// Create a new http client that always attaches an user agent.
func NewWithUA(c HTTPClient, ua string) *WithUA {
	return &WithUA{
		c:  c,
		ua: ua,
	}
}

var (
	_ HTTPClient = (*WithUA)(nil)
)

func (c *WithUA) Do(req *http.Request) (*http.Response, error) {
	req.Header.Set("user-agent", c.ua)
	return c.c.Do(req)
}
