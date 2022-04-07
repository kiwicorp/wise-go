package wisesca

import (
	"net/http"

	"github.com/kiwicorp/wise-go/pkg/wisehttp"
)

// Strong customer authentication (SCA) client for the Wise API.
//
// See https://api-docs.transferwise.com/#strong-customer-authentication.
type SCAClient interface {
	// Do a request for an SCA-protected endpoint.
	//
	// Will return an error if signing the one-time token fails.
	Do(*http.Request) (*http.Response, error)
}

var (
	// ensure that an SCA client conforms to the http client interface
	_ wisehttp.HTTPClient = (SCAClient)(nil)
)
