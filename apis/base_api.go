package apis

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"runtime/debug"

	"github.com/carlmjohnson/requests"
	"github.com/volvlabs/towncryer-sdk/config"
)

const sdkModulePath = "github.com/volvlabs/towncryer-sdk"

// sdkVersion returns the version of this module as recorded by the Go module
// system in the consuming binary's build info, falling back to "dev" when
// that information isn't available (e.g. when running the SDK's own tests).
func sdkVersion() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return "dev"
	}

	if info.Main.Path == sdkModulePath && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	for _, dep := range info.Deps {
		if dep.Path == sdkModulePath {
			return dep.Version
		}
	}

	return "dev"
}

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
	timeout := config.Timeout
	if timeout <= 0 {
		timeout = DefaultTimeout
	}

	maxRetries := config.MaxRetries
	if maxRetries <= 0 {
		maxRetries = DefaultMaxRetries
	}

	retryDelay := config.RetryDelay
	if retryDelay <= 0 {
		retryDelay = DefaultRetryDelay
	}

	userAgent := fmt.Sprintf("TowncryerSDK/%s (Go/%s; %s/%s)", sdkVersion(), runtime.Version(), runtime.GOOS, runtime.GOARCH)

	client := &http.Client{Timeout: timeout}
	var baseTransport http.RoundTripper
	if config.HTTPClient != nil {
		injected := *config.HTTPClient
		if injected.Timeout <= 0 {
			injected.Timeout = timeout
		}
		client = &injected
		baseTransport = injected.Transport
	}

	requestsBuilder := requests.URL(fmt.Sprintf("%s/%s", config.BaseUrl, endpoint)).
		Header(APIHeaderKey, config.ApiKey).
		Header(APIHeaderSecretKey, config.ApiSecret).
		Header(http.CanonicalHeaderKey("User-Agent"), userAgent).
		ContentType("application/json").
		Accept("application/json").
		Client(client).
		Transport(RetryTransport(maxRetries, retryDelay)).
		Transport(PayloadSizeLimitTransport(baseTransport)).
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
