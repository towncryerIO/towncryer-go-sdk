package apis

import (
	"encoding/json"
	"fmt"
	"io"
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
		Transport(RateLimitTransport(nil)).
		AddValidator(apiErrorValidator)
	return &BaseAPI{config: config, requestsBuilder: requestsBuilder}
}

// apiErrorBody mirrors the error-carrying fields of models.ApiResponse, without
// requiring a type parameter for the (irrelevant, on an error path) data payload.
type apiErrorBody struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
}

// apiErrorValidator replaces the default 2xx status check so that non-2xx
// responses are decoded into an *APIError, preserving the status code and any
// error detail the API reported in the response body.
func apiErrorValidator(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return nil
	}

	body, _ := io.ReadAll(res.Body)

	var decoded apiErrorBody
	_ = json.Unmarshal(body, &decoded)

	return &APIError{
		StatusCode: res.StatusCode,
		Code:       decoded.Code,
		Message:    decoded.Message,
		Errors:     decoded.Errors,
	}
}
