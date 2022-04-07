package wise

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/kiwicorp/wise-go/pkg/wisehttp"
	"github.com/kiwicorp/wise-go/pkg/wisesca"
)

var (
	_ WiseClient = (*Client)(nil)
)

type Client struct {
	baseURL    *url.URL
	httpClient wisehttp.HTTPClient

	token string
}

// Create a new sandbox API client with the default settings.
func NewDefaultSandbox(token string) *Client {
	return NewSandbox(wisehttp.NewWithUA(&http.Client{}, defaultUA), token)
}

// Create a new production API client with the default settings.
func NewDefaultProduction(sca wisesca.SCAClient, token string) *Client {
	return NewProduction(wisehttp.NewWithUA(sca, defaultUA), token)
}

// Create a new sandbox API client.
func NewSandbox(httpClient wisehttp.HTTPClient, token string) *Client {
	return &Client{
		baseURL:    &*sandboxAPIBaseURL,
		httpClient: httpClient,
		token:      token,
	}
}

// Create a new production API client.
func NewProduction(httpClient wisehttp.HTTPClient, token string) *Client {
	return &Client{
		baseURL:    &*productionAPIBaseURL,
		httpClient: httpClient,
		token:      token,
	}
}

func (c *Client) GetBalances(ctx context.Context, req *GetBalancesRequest) (*GetBalancesResponse, error) {
	u := c.baseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("v3/profiles/%d/balances", req.ProfileID),
	})
	q := u.Query()
	q.Add("types", "STANDARD")
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	if err != nil {
		panic(err)
	}
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}
	var balances []Balance
	if err := json.Unmarshal(b, &balances); err != nil {
		return nil, err
	}
	return &GetBalancesResponse{
		Balances: balances,
	}, nil
}

func (c *Client) GetStatementPDF(ctx context.Context, req *GetStatementPDFRequest) (*GetStatementPDFResponse, error) {
	u := c.baseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("v1/profiles/%d/balance-statements/%d/statement.pdf", req.ProfileID, req.BalanceID),
	})
	q := u.Query()
	q.Add("intervalStart", req.IntervalStart.Format("2006-01-02T15:04:05.000Z"))
	q.Add("intervalEnd", req.IntervalEnd.Format("2006-01-02T15:04:05.000Z"))
	q.Add("type", string(req.Type))
	u.RawQuery = q.Encode()
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	if err != nil {
		panic(err)
	}
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New(string(b))
	}
	return &GetStatementPDFResponse{
		Data: b,
	}, nil
}

func (c *Client) ListProfiles(ctx context.Context) (*ListProfilesResponse, error) {
	u := c.baseURL.ResolveReference(&url.URL{
		Path: fmt.Sprintf("v2/profiles"),
	})
	r, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.token))
	if err != nil {
		panic(err)
	}
	resp, err := c.httpClient.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		var wiseErr Error
		if err := json.Unmarshal(b, &wiseErr); err != nil {
			return nil, err
		}
		return nil, wiseErr
	}
	var profiles []Profile
	if err := json.Unmarshal(b, &profiles); err != nil {
		return nil, err
	}
	return &ListProfilesResponse{
		Profiles: profiles,
	}, nil
}
