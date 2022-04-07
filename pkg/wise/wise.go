package wise

import (
	"context"
	"time"
)

// Wise API client.
type WiseClient interface {
	// Retrieve the user's multi-currency account balances. It returns all
	// currency balances the profile has.
	GetBalances(context.Context, *GetBalancesRequest) (*GetBalancesResponse, error)

	// Get a multi-currency account statement for one balance or Jar and for the
	// specified time range. The period between intervalStart and intervalEnd
	// cannot exceed 455 days (around 1 year and 3 month)s. PDF format.
	GetStatementPDF(context.Context, *GetStatementPDFRequest) (*GetStatementPDFResponse, error)

	// Get a list of all profiles belonging to user.
	ListProfiles(context.Context) (*ListProfilesResponse, error)
}

type GetStatementPDFRequest struct {
	ProfileID     int
	BalanceID     int
	IntervalStart time.Time
	IntervalEnd   time.Time
	Type          StatementType
}

type GetStatementPDFResponse struct {
	Data []byte
}

type GetBalancesRequest struct {
	ProfileID int
}

type GetBalancesResponse struct {
	Balances []Balance
}

type ListProfilesResponse struct {
	Profiles []Profile
}
