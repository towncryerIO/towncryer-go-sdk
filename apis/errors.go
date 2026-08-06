package apis

import (
	"errors"
	"fmt"
)

var (
	ErrRequestPayloadLimitReached = fmt.Errorf("request limit reached, max allowed is %d", MAX_REQUEST_PAYLOAD_SIZE)
	ErrPayloadCannotBeNil         = errors.New("payload cannot be nil")
)

// APIError represents a non-2xx response from the Towncryer API. It carries the
// HTTP status code alongside the error code, message, and field-level errors
// reported in the response body, so callers can distinguish API-level failures
// (e.g. bad credentials, validation errors) from network-level failures using
// errors.As, rather than string-matching on the error message.
type APIError struct {
	// StatusCode is the HTTP status code returned by the API.
	StatusCode int
	// Code is the API's own error code, if the response body included one.
	Code string
	// Message is the API's human-readable error message, if the response body included one.
	Message string
	// Errors holds field-level validation errors reported by the API, if any.
	Errors []string
}

func (e *APIError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("towncryer: api error (status %d, code %q): %s", e.StatusCode, e.Code, e.Message)
	}
	return fmt.Sprintf("towncryer: api error (status %d)", e.StatusCode)
}
