package apis_test

import (
	"context"
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk/apis"
	"github.com/volvlabs/towncryer-sdk/helpers"
	"github.com/volvlabs/towncryer-sdk/models"
)

func TestAPIError_Error(t *testing.T) {
	scenarios := []struct {
		name     string
		err      *apis.APIError
		expected string
	}{
		{
			name:     "with a message",
			err:      &apis.APIError{StatusCode: http.StatusBadRequest, Code: "40", Message: "invalid payload"},
			expected: `towncryer: api error (status 400, code "40"): invalid payload`,
		},
		{
			name:     "without a message",
			err:      &apis.APIError{StatusCode: http.StatusInternalServerError},
			expected: "towncryer: api error (status 500)",
		},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			assert.Equal(t, scenario.expected, scenario.err.Error())
		})
	}
}

func TestNonSuccessResponses_ReturnTypedAPIError(t *testing.T) {
	// Arrange:
	config, server := helpers.CreateHttpServer(
		t,
		"/customers",
		http.StatusUnauthorized,
		[]byte(`{"code": "40", "message": "invalid credentials", "errors": ["api key is invalid"]}`),
	)
	defer server.Close()

	customersApi := apis.NewCustomers(config)
	payload := &models.CreateCustomerPayload{
		CustomerId: "CUST-12345",
		FirstName:  "John",
		LastName:   "Doe",
	}

	// Act:
	_, err := customersApi.CreateCustomer(context.Background(), payload)

	// Assert:
	var apiErr *apis.APIError
	if assert.True(t, errors.As(err, &apiErr), "expected err to be an *apis.APIError, got %T: %v", err, err) {
		assert.Equal(t, http.StatusUnauthorized, apiErr.StatusCode)
		assert.Equal(t, "40", apiErr.Code)
		assert.Equal(t, "invalid credentials", apiErr.Message)
		assert.Equal(t, []string{"api key is invalid"}, apiErr.Errors)
	}
}
