package wise

import "time"

var (
	_ error = Error{}
)

// A Wise API error.
type Error struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status"`
	Err       string    `json:"error"`
	Message   string    `json:"message"`
	Path      string    `json:"path"`
}

// Error implements error
func (Error) Error() string {
	panic("unimplemented")
}
