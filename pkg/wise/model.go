package wise

import "time"

// Account statement type.
type StatementType string

const (
	// COMPACT for a single statement line per transaction.
	StatementTypeCompact = StatementType("COMPACT")
	// FLAT for accounting statements where transaction fees are on a separate
	// line.
	StatementTypeFlat = StatementType("FLAT")
)

// Profile type.
type ProfileType string

const (
	// Profile type for a personal account.
	ProfileTypePersonal = StatementType("PERSONAL")
	// Profile type for a business account.
	ProfileTypeBusiness = StatementType("BUSINESS")
)

type Balance struct {
	ID               int           `json:"id"`
	Currency         string        `json:"currency"`
	Type             string        `json:"type"`
	Name             string        `json:"name"`
	Icon             string        `json:"icon"`
	Amount           BalanceAmount `json:"amount"`
	CreationTime     time.Time     `json:"creationTime"`
	ModificationTime time.Time     `json:"modificationTime"`
	Visible          bool          `json:"visible"`
}

type BalanceAmount struct {
	Value    float32 `json:"value"`
	Currency string  `json:"currency"`
}

type Profile struct {
	ID           int            `json:"id"`
	Type         string         `json:"type"`
	UserID       int            `json:"userId"`
	Address      ProfileAddress `json:"address"`
	Email        string         `json:"email"`
	CreatedAd    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	Obfuscated   bool           `json:"obfuscated"`
	CurrentState string         `json:"currentState"`
	FullName     string         `json:"fullName"`
	// incomplete
}

type ProfileAddress struct {
	AddressFirstLine string `json:"addressFirstLine"`
	City             string `json:"city"`
	CountryISO2Code  string `json:"countryIso2Code"`
	CountryISO3Code  string `json:"countryIso3Code"`
	PostCode         string `json:"postCode"`
	StateCode        string `json:"stateCode"`
}
