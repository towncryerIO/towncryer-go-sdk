package apis_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/towncryerIO/towncryer-go-sdk/apis"
	"github.com/towncryerIO/towncryer-go-sdk/config"
	"github.com/towncryerIO/towncryer-go-sdk/internal/helpers"
	"github.com/towncryerIO/towncryer-go-sdk/models"
)

func TestPublishEvent(t *testing.T) {
	validRequestPayload := &models.PublishEventPayload{
		Name: "UserRegisteredEvent",
		Customer: models.CreateCustomerPayload{
			CustomerId: "CUST-1234",
			FirstName:  "John",
			LastName:   "Doe",
		},
		Data: map[string]any{
			"testKey": "testValue",
		},
	}
	type args struct {
		payload *models.PublishEventPayload
	}
	scenarios := []struct {
		name                string
		args                args
		mockResponse        helpers.MockResponse
		expectedApiResponse *models.ApiResponse[map[string]any]
		errorIsNil          bool
	}{
		{
			name: "should publish event successfully",
			args: args{
				payload: validRequestPayload,
			},
			mockResponse: helpers.MockResponse{
				StatusCode: http.StatusOK,
				Payload: []byte(`
					{
						"code": "00",
						"message": "event published successfully"
					}
				`),
			},
			expectedApiResponse: &models.ApiResponse[map[string]any]{
				Code:    "00",
				Message: "event published successfully",
			},
			errorIsNil: true,
		},
		{
			name: "should return error because request payload is invalid",
			args: args{
				payload: &models.PublishEventPayload{},
			},
			mockResponse:        helpers.MockResponse{},
			expectedApiResponse: nil,
			errorIsNil:          false,
		},
		{
			name: "should return error because server returned a non-200 status code",
			args: args{
				payload: validRequestPayload,
			},
			mockResponse: helpers.MockResponse{
				StatusCode: http.StatusInternalServerError,
				Payload:    []byte(`{"code": "500", "message": "internal server error"}`),
			},
			expectedApiResponse: nil,
			errorIsNil:          false,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Arrange:
			config, server := helpers.CreateHttpServer(
				t,
				"/events",
				scenario.mockResponse.StatusCode,
				scenario.mockResponse.Payload,
			)
			defer server.Close()

			eventsApi := apis.NewEvents(config)
			ctx := context.Background()

			// Act:
			apiResponse, err := eventsApi.PublishEvent(ctx, scenario.args.payload)

			// Assert:
			assert.Equal(t, scenario.errorIsNil, err == nil, err)
			assert.Equal(t, scenario.expectedApiResponse, apiResponse)
		})
	}
}

func TestPublishEventRetryOptIn(t *testing.T) {
	payload := &models.PublishEventPayload{
		Name: "UserRegisteredEvent",
		Customer: models.CreateCustomerPayload{
			CustomerId: "CUST-1234",
			FirstName:  "John",
			LastName:   "Doe",
		},
	}

	scenarios := []struct {
		name            string
		optIn           bool
		expectedAttempt int
	}{
		{name: "should not retry by default", optIn: false, expectedAttempt: 1},
		{name: "should retry when explicitly opted in", optIn: true, expectedAttempt: apis.DefaultMaxRetries},
	}

	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Arrange:
			requestCount := 0
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				requestCount++
				w.WriteHeader(http.StatusInternalServerError)
			}))
			defer server.Close()

			cfg := &config.Configuration{
				BaseUrl:                               server.URL,
				ApiKey:                                "dummy_api_key",
				ApiSecret:                             "dummy_api_secret",
				EnableRetriesForNonIdempotentRequests: scenario.optIn,
			}

			// Act:
			_, err := apis.NewEvents(cfg).PublishEvent(context.Background(), payload)

			// Assert:
			assert.Error(t, err)
			assert.Equal(t, scenario.expectedAttempt, requestCount)
		})
	}
}

// TestPublishEvent_WithoutCustomerNames covers a customer the caller knows only
// by email address: CustomerId is the one required field, not the names.
func TestPublishEvent_WithoutCustomerNames(t *testing.T) {
	config, server := helpers.CreateHttpServer(t, "/events", http.StatusOK, []byte(`
		{"code": "00", "message": "event published successfully"}
	`))
	defer server.Close()

	apiResponse, err := apis.NewEvents(config).PublishEvent(context.Background(), &models.PublishEventPayload{
		Name: "UserRegisteredEvent",
		Customer: models.CreateCustomerPayload{
			CustomerId: "CUST-1234",
			Email:      "nameless@example.com",
		},
	})

	assert.NoError(t, err)
	assert.Equal(t, "00", apiResponse.Code)
}
