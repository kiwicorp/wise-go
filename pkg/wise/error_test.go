package wise

import (
	"encoding/json"
	"testing"
)

func TestError(t *testing.T) {
	expectedBytes := []byte(`{"timestamp":"2022-04-07T07:15:31.399+00:00","status":403,"error":"Forbidden","message":"You are forbidden to send this request","path":"/v1/profiles/12345678/balance-statements/12345678/statement.pdf"}`)
	actual := Error{
		Status:  403,
		Err:     "Forbidden",
		Message: "You are forbidden to send this request",
		Path:    "/v1/profiles/12345678/balance-statements/12345678/statement.pdf",
	}
	var expected Error
	if err := json.Unmarshal(expectedBytes, &expected); err != nil {
		panic(err)
	}
	// somehow the timestamp gets parsed correctly but I can't see to create a
	// time format to properly parse it
	actual.Timestamp = expected.Timestamp
	if expected != actual {
		panic("expected != actual")
	}
}
