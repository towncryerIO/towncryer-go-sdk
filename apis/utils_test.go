package apis_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/carlmjohnson/requests"
	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk/apis"
)

type MockTransport struct {
	RoundTripFn func(*http.Request) (*http.Response, error)
}

func (m *MockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFn(req)
}

func TestRateLimitTransport(t *testing.T) {
	scenarios := []struct {
		name        string
		payload     string
		wantErr     error
		contentType string
	}{
		{
			name:        "should pass - payload within limit",
			payload:     "small payload",
			wantErr:     nil,
			contentType: "application/json",
		},
		{
			name:        "should fail - payload exceeds limit",
			payload:     strings.Repeat("a", apis.MAX_REQUEST_PAYLOAD_SIZE+1),
			wantErr:     apis.ErrRequestPayloadLimitReached,
			contentType: "application/json",
		},
		{
			name:        "should pass - empty payload",
			payload:     "",
			wantErr:     nil,
			contentType: "application/json",
		},
		{
			name:        "should pass - payload at exact limit",
			payload:     strings.Repeat("a", apis.MAX_REQUEST_PAYLOAD_SIZE),
			wantErr:     nil,
			contentType: "application/json",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			mockTransport := &MockTransport{
				RoundTripFn: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(bytes.NewBufferString("mock response")),
					}, nil
				},
			}

			ctx := context.Background()

			req, err := http.NewRequestWithContext(ctx, http.MethodPost, "http://example.com", bytes.NewBufferString(scenario.payload))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", scenario.contentType)
			req.ContentLength = int64(len(scenario.payload))

			transport := apis.RateLimitTransport(mockTransport)
			resp, err := transport.RoundTrip(req)

			// Assert:
			if scenario.wantErr != nil {
				assert.Equal(t, scenario.wantErr, err)
			} else {
				assert.NoError(t, err)
			}

			if resp != nil {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				defer resp.Body.Close()

				assert.Equal(t, "mock response", string(body))
			}
		})
	}
}

func TestRateLimitTransportNilTransport(t *testing.T) {
	transport := apis.RateLimitTransport(nil)

	assert.NotNil(t, transport)
}

func TestRetryTransport(t *testing.T) {
	tests := []struct {
		name            string
		statusCode      int
		expectedRetries int
	}{
		{
			name:            "Retry on TooManyRequests",
			statusCode:      http.StatusTooManyRequests,
			expectedRetries: apis.DEFAULT_MAX_RETRIES,
		},
		{
			name:            "Retry on RetryableStatusCode",
			statusCode:      http.StatusInternalServerError,
			expectedRetries: apis.DEFAULT_MAX_RETRIES,
		},
		{
			name:            "No retry on 200 OK",
			statusCode:      http.StatusOK,
			expectedRetries: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange:
			retryCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if retryCount < tt.expectedRetries {
					retryCount++
					w.WriteHeader(tt.statusCode)
				} else {
					w.WriteHeader(http.StatusOK)
				}
			}))
			defer server.Close()

			ctx := context.Background()
			var response string

			// Act:
			requests.URL(server.URL).
				Transport(apis.RetryTransport(apis.DEFAULT_MAX_RETRIES, apis.DEFAULT_RETRY_DELAY)).
				ToString(&response).
				//nolint:errcheck
				Fetch(ctx)

			// Assert:
			assert.Equal(t, tt.expectedRetries, retryCount)
		})
	}
}
