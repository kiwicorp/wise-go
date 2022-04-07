package wisesca

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"net/http"

	"github.com/kiwicorp/wise-go/pkg/wisehttp"
)

var (
	_ SCAClient = (*PersonalToken)(nil)
)

// An SCA client that uses a personal token.
//
// See https://api-docs.transferwise.com/#strong-customer-authentication-personal-token.
type PersonalToken struct {
	client wisehttp.HTTPClient
	key    *rsa.PrivateKey
}

// Create a SCA client with a personal token.
func NewWithPersonalToken(client wisehttp.HTTPClient, key *rsa.PrivateKey) *PersonalToken {
	return &PersonalToken{
		client: client,
		key:    key,
	}
}

func (c *PersonalToken) Do(req *http.Request) (*http.Response, error) {
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusForbidden {
		return resp, nil
	}
	ott := resp.Header.Get("x-2fa-approval")
	if ott == "" {
		return resp, nil
	}
	ottHash := sha256.Sum256([]byte(ott))
	ottSig, err := rsa.SignPKCS1v15(rand.Reader, c.key, crypto.SHA256, ottHash[:])
	if err != nil {
		return nil, err
	}
	sig := base64.StdEncoding.EncodeToString(ottSig)
	req.Header.Set("x-2fa-approval", ott)
	req.Header.Set("X-Signature", sig)
	return c.client.Do(req)
}
