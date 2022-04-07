package wise

import (
	"fmt"
	"net/url"
)

const (
	// Library version.
	version = "0.0.0"
)

var (
	// Default user agent.
	defaultUA = fmt.Sprintf("wise-go / %s", version)
)

var (
	sandboxAPIBaseURL    *url.URL
	productionAPIBaseURL *url.URL
)

func init() {
	var err error
	sandboxAPIBaseURL, err = url.Parse("https://api.sandbox.transferwise.tech/")
	if err != nil {
		panic(fmt.Sprintf("wise-go: invalid sandbox base URL, please open an issue at https://github.com/kiwicorp/wise-go: %v", err))
	}
	productionAPIBaseURL, err = url.Parse("https://api.transferwise.com/")
	if err != nil {
		panic(fmt.Sprintf("wise-go: invalid production base URL, please open an issue at https://github.com/kiwicorp/wise-go: %v", err))
	}
}
