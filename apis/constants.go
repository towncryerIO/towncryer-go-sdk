package apis

import "time"

const (
	APIHeaderKey       = "X-API-KEY"
	APIHeaderSecretKey = "X-API-SECRET"

	MaxRequestPayloadSize = 32768 // 32kb

	// DefaultMaxRetries and DefaultRetryDelay are used when a
	// config.Configuration does not specify MaxRetries / RetryDelay.
	DefaultMaxRetries = 3
	DefaultRetryDelay = 150 * time.Millisecond

	// DefaultTimeout is applied to the underlying http.Client when a
	// config.Configuration does not specify a Timeout. It guarantees requests
	// fail instead of hanging forever when the caller passes a context with
	// no deadline (e.g. context.Background()).
	DefaultTimeout = 10 * time.Second
)
