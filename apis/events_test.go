package apis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/volvlabs/towncryer-sdk/apis"
	"github.com/volvlabs/towncryer-sdk/helpers"
	"github.com/volvlabs/towncryer-sdk/models"
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
