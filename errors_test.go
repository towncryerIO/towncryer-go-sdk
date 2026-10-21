package towncryer_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk"
	"github.com/volvlabs/towncryer-sdk/internal/helpers"
	"github.com/volvlabs/towncryer-sdk/models"
)

func TestAPIError_ErrorsAs(t *testing.T) {
	// Arrange:
	cfg, server := helpers.CreateHttpServer(
		t,
		"/events",
		http.StatusBadRequest,
		[]byte(`{"code": "40", "message": "invalid event", "errors": ["name is required"]}`),
	)
	defer server.Close()

	sdk := towncryer.New(cfg)
	payload := &models.PublishEventPayload{
		Name: "UserRegisteredEvent",
		Customer: models.CreateCustomerPayload{
			CustomerId: "CUST-1234",
			FirstName:  "John",
			LastName:   "Doe",
		},
	}

	// Act:
	_, err := sdk.EventsApi().PublishEvent(context.Background(), payload)

	// Assert:
	var apiErr *towncryer.APIError
	if assert.True(t, errors.As(err, &apiErr), "expected err to be a *towncryer.APIError, got %T: %v", err, err) {
		assert.Equal(t, http.StatusBadRequest, apiErr.StatusCode)
		assert.Equal(t, "40", apiErr.Code)
		assert.Equal(t, "invalid event", apiErr.Message)
		assert.Equal(t, []string{"name is required"}, apiErr.Errors)
	}
}
