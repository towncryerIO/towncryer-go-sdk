package apis

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"time"

	"github.com/carlmjohnson/requests"
	"github.com/go-playground/validator/v10"
	"github.com/philippta/trip"
)

// validatePayload checks if the given payload is valid according to predefined
// struct validation rules.
//
// Parameters:
//   - payload (any): The payload to validate. It should be a struct that follows
//     the necessary validation tags required for the application.
//
// Returns:
//   - error: An error if the payload is nil or if validation fails; otherwise, nil.
//
// Usage:
//
//	If the payload is nil, validatePayload will return ErrPayloadCannotBeNil.
//	Otherwise, it will validate the payload based on struct field tags, returning
//	any error if the payload does not meet required validation constraints.
func validatePayload(payload any) error {
	if reflect.ValueOf(payload).IsNil() {
		return ErrPayloadCannotBeNil
	}

	validate := validator.New(validator.WithRequiredStructEnabled())
	return validate.Struct(payload)
}

// PayloadSizeLimitTransport wraps an existing http.RoundTripper to enforce request payload size limits.
// It limits the size of request bodies to a maximum defined by MaxRequestPayloadSize.
//
// Parameters:
//   - rt (http.RoundTripper): The underlying RoundTripper to wrap. If nil, it defaults
//     to http.DefaultTransport.
//
// Returns:
//   - http.RoundTripper: A RoundTripper that enforces a payload size limit.
//
// Behavior:
//
//	If the request payload exceeds MaxRequestPayloadSize, PayloadSizeLimitTransport will return
//	an error (ErrRequestPayloadLimitReached) without sending the request. If the payload size
//	is within the limit, it passes the request through to the underlying RoundTripper.
func PayloadSizeLimitTransport(rt http.RoundTripper) http.RoundTripper {
	if rt == nil {
		rt = http.DefaultTransport
	}

	return requests.RoundTripFunc(func(req *http.Request) (res *http.Response, err error) {
		if req.Body != nil {
			var bodyBuffer bytes.Buffer
			limitReader := io.LimitReader(req.Body, int64(MaxRequestPayloadSize)+1)

			written, err := io.Copy(&bodyBuffer, limitReader)
			if err != nil {
				return nil, fmt.Errorf("failed to read request body: %w", err)
			}

			if written > int64(MaxRequestPayloadSize) {
				return nil, ErrRequestPayloadLimitReached
			}

			req.Body = io.NopCloser(&bodyBuffer)
		}

		return rt.RoundTrip(req)
	})
}

// RetryTransport configures an HTTP RoundTripper with retry logic. It automatically retries
// requests that fail with a network error or receive a retryable HTTP status code (e.g. 429
// Too Many Requests, or any other code in trip.RetryableStatusCodes), up to maxRetries times,
// waiting retryDelay between attempts.
//
// Parameters:
//   - rt (http.RoundTripper): The underlying RoundTripper to wrap. If nil, it defaults
//     to http.DefaultTransport.
//   - maxRetries (int): The maximum number of attempts made for a single request.
//   - retryDelay (time.Duration): The fixed delay between retry attempts.
//
// Returns:
//   - http.RoundTripper: A RoundTripper configured to retry requests based on a single,
//     unambiguous retry rule.
//
// Only apply this to idempotent requests (safe to send more than once). Retrying a
// non-idempotent send/publish call can cause the server to process it more than once.
func RetryTransport(rt http.RoundTripper, maxRetries int, retryDelay time.Duration) http.RoundTripper {
	return trip.New(rt,
		trip.Retry(maxRetries, retryDelay, trip.RetryableStatusCodes...),
	)
}
