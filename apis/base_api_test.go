package apis_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/towncryerIO/towncryer-go-sdk/apis"
	"github.com/towncryerIO/towncryer-go-sdk/config"
	"github.com/towncryerIO/towncryer-go-sdk/models"
)

// TestDefaultTimeoutAppliedWithoutContextDeadline verifies that a request made with
// context.Background() (no deadline) still fails within the configured timeout instead
// of hanging forever when the server never responds.
func TestDefaultTimeoutAppliedWithoutContextDeadline(t *testing.T) {
	blockedTimeout := 200 * time.Millisecond

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
	}))
	defer server.Close()

	cfg := &config.Configuration{
		BaseUrl:    server.URL,
		ApiKey:     "dummy_api_key",
		ApiSecret:  "dummy_api_secret",
		Timeout:    blockedTimeout,
		MaxRetries: 1,
	}

	customers := apis.NewCustomers(cfg)

	start := time.Now()
	_, err := customers.CreateCustomer(context.Background(), &models.CreateCustomerPayload{
		CustomerId: "customer-1",
		FirstName:  "Ada",
		LastName:   "Lovelace",
	})
	elapsed := time.Since(start)

	assert.Error(t, err)
	assert.Less(t, elapsed, 1*time.Second, "request should fail within the configured timeout instead of hanging")
}

// mockRoundTripper lets a test stub out HTTP responses without making a real
// network connection.
type mockRoundTripper struct {
	response *http.Response
	requests []*http.Request
}

func (m *mockRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	m.requests = append(m.requests, req)
	return m.response, nil
}

func TestNewBaseAPI_UsesInjectedHTTPClient(t *testing.T) {
	// Arrange: a client whose transport never touches the network.
	mockTransport := &mockRoundTripper{
		response: &http.Response{
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewBufferString(`{"code":"00","message":"ok","data":{}}`)),
			Header:     http.Header{"Content-Type": []string{"application/json"}},
		},
	}
	cfg := &config.Configuration{
		BaseUrl:    "https://unreachable.invalid",
		ApiKey:     "dummy_api_key",
		ApiSecret:  "dummy_api_secret",
		HTTPClient: &http.Client{Transport: mockTransport},
	}

	eventsApi := apis.NewEvents(cfg)

	// Act:
	_, err := eventsApi.PublishEvent(context.Background(), &models.PublishEventPayload{
		Name: "test-event",
		Customer: models.CreateCustomerPayload{
			CustomerId: "CUST-12345",
			FirstName:  "John",
			LastName:   "Doe",
		},
	})

	// Assert: request was routed through the injected transport, not the network.
	assert.NoError(t, err)
	assert.Len(t, mockTransport.requests, 1)
	assert.Equal(t, "unreachable.invalid", mockTransport.requests[0].URL.Host)
}
