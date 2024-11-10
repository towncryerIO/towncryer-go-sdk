package apis

import (
	"fmt"
	"net/http"
	"runtime"

	"github.com/carlmjohnson/requests"
	"github.com/volvlabs/towncryer-sdk/config"
)

// BaseAPI represents the base configuration for interacting with the API,
// including essential configurations and request setup for consistent behavior
// across API calls.
//
// Fields:
//   - config (*config.Configuration): Holds configuration settings required for the API connection.
//   - requestsBuilder (*requests.Builder): A pre-configured request builder that
//     manages headers, endpoint, content types, and transport settings.
type BaseAPI struct {
	config          *config.Configuration
	requestsBuilder *requests.Builder
}

// NewBaseAPI creates and returns a new BaseAPI instance with the specified
// configuration and endpoint. It applies necessary headers and configures
// retry and rate limit transports to handle retries and payload size limits.
//
// Parameters:
//   - config (*config.Configuration): The configuration struct containing the base
//     URL and authentication keys for the API.
//   - endpoint (string): The specific API endpoint to connect to, which will
//     be appended to the base URL.
//
// Returns:
//   - *BaseAPI: A pointer to the initialized BaseAPI instance with configured
//     headers, content types, retry, and rate limit transports.
func NewBaseAPI(config *config.Configuration, endpoint string) *BaseAPI {
	userAgent := fmt.Sprintf("TowncryerSDK/0.0.1 (Go/%s; %s/%s)", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	requestsBuilder := requests.URL(fmt.Sprintf("%s/%s", config.BaseUrl, endpoint)).
		Header(API_HEADER_KEY, config.ApiKey).
		Header(API_HEADER_SECRET_KEY, config.ApiSecret).
		Header(http.CanonicalHeaderKey("User-Agent"), userAgent).
		ContentType("application/json").
		Accept("application/json").
		Transport(RetryTransport()).
		Transport(RateLimitTransport(nil))
	return &BaseAPI{config: config, requestsBuilder: requestsBuilder}
}
