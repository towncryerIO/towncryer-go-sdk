package apis_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk/apis"
	"github.com/volvlabs/towncryer-sdk/config"
	"github.com/volvlabs/towncryer-sdk/models"
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
