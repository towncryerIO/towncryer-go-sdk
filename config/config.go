package config

import "time"

// Configuration holds the settings needed to connect to the Towncryer API.
//
// Timeout, MaxRetries, and RetryDelay are optional. When left at their zero
// value, the SDK applies its own defaults (see apis.DEFAULT_TIMEOUT,
// apis.DEFAULT_MAX_RETRIES, apis.DEFAULT_RETRY_DELAY) so requests never hang
// indefinitely and retry behavior is predictable out of the box.
type Configuration struct {
	ApiKey    string
	ApiSecret string
	BaseUrl   string

	// Timeout bounds how long a single request (including retries) may take.
	Timeout time.Duration

	// MaxRetries is the number of times a request is retried after a
	// retryable failure (network error or retryable status code).
	MaxRetries int

	// RetryDelay is the fixed delay between retry attempts.
	RetryDelay time.Duration
}
