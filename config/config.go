package config

import (
	"net/http"
	"os"
	"time"
)

const (
	envApiKey    = "TOWNCRYER_API_KEY"
	envApiSecret = "TOWNCRYER_API_SECRET"
	envBaseUrl   = "TOWNCRYER_BASE_URL"
)

// Configuration holds the settings needed to connect to the Towncryer API.
//
// ApiKey, ApiSecret, and BaseUrl fall back to the TOWNCRYER_API_KEY,
// TOWNCRYER_API_SECRET, and TOWNCRYER_BASE_URL environment variables
// respectively whenever left blank, so a zero-value Configuration can be
// used in env-var-driven (12-factor-style) deployments.
//
// Timeout, MaxRetries, and RetryDelay are optional. When left at their zero
// value, the SDK applies its own defaults (see apis.DefaultTimeout,
// apis.DefaultMaxRetries, apis.DefaultRetryDelay) so requests never hang
// indefinitely and retry behavior is predictable out of the box.
//
// HTTPClient, if set, is used for all outgoing requests instead of
// http.DefaultClient. This allows consumers to inject a custom transport,
// e.g. for proxies in production or mock transports in tests.
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

	// HTTPClient optionally supplies the client (and its Transport) used for requests.
	HTTPClient *http.Client
}

// WithEnvDefaults returns a copy of cfg with any blank ApiKey, ApiSecret,
// or BaseUrl fields filled in from the TOWNCRYER_API_KEY,
// TOWNCRYER_API_SECRET, and TOWNCRYER_BASE_URL environment variables.
func (cfg Configuration) WithEnvDefaults() Configuration {
	if cfg.ApiKey == "" {
		cfg.ApiKey = os.Getenv(envApiKey)
	}
	if cfg.ApiSecret == "" {
		cfg.ApiSecret = os.Getenv(envApiSecret)
	}
	if cfg.BaseUrl == "" {
		cfg.BaseUrl = os.Getenv(envBaseUrl)
	}
	return cfg
}
