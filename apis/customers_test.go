package apis_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/towncryerIO/towncryer-go-sdk/apis"
	"github.com/towncryerIO/towncryer-go-sdk/internal/helpers"
	"github.com/towncryerIO/towncryer-go-sdk/models"
)

func TestCreateCustomer(t *testing.T) {
	// Arrange:
	sampleRequestPayload := &models.CreateCustomerPayload{
		CustomerId:  "CUST-12345",
		FirstName:   "John",
		LastName:    "Doe",
		Email:       "john.doe@example.com",
		PhoneNumber: "123-456-7890",
	}

	sampleCustomerResponsePayload := []byte(`
		{
			"code": "00",
			"message": "customer created successfully",
			"data": {
				"id": "1aadd344-613d-448c-89b2-3bb3b7f458b2", 
				"customerId": "CUST-12345", 
				"firstName": "John", 
				"lastName": "Doe",
				"emails": [
					"john.doe@example.com"
				], 
				"phoneNumbers": [
					"123-456-7890"
				], 
				"pushNotificationTokens": ["token123"]
			}
		}
	`)

	type args struct {
		payload *models.CreateCustomerPayload
	}
	scenarios := []struct {
		name                string
		args                args
		mockResponse        helpers.MockResponse
		expectedApiResponse *models.ApiResponse[models.Customer]
		errorIsNil          bool
	}{
		{
			name: "should send create customers successfully",
			args: args{
				payload: sampleRequestPayload,
			},
			mockResponse: helpers.MockResponse{
				StatusCode: http.StatusOK,
				Payload:    sampleCustomerResponsePayload,
			},
			expectedApiResponse: &models.ApiResponse[models.Customer]{
				Code:    "00",
				Message: "customer created successfully",
				Data: models.Customer{
					Id:         uuid.MustParse("1aadd344-613d-448c-89b2-3bb3b7f458b2"),
					CustomerId: "CUST-12345",
					FirstName:  "John",
					LastName:   "Doe",
					Emails: []string{
						"john.doe@example.com",
					},
					PushNotificationTokens: []string{
						"token123",
					},
					PhoneNumbers: []string{
						"123-456-7890",
					},
				},
			},
			errorIsNil: true,
		},
		{
			name: "should return error because payload passed is nil",
			args: args{
				payload: nil,
			},
			mockResponse:        helpers.MockResponse{},
			expectedApiResponse: nil,
			errorIsNil:          false,
		},
		{
			name: "should return error because customerId is not present in request payload",
			args: args{
				payload: &models.CreateCustomerPayload{
					Email:       "john.doe@example.com",
					PhoneNumber: "123-456-7890",
				},
			},
			mockResponse:        helpers.MockResponse{},
			expectedApiResponse: nil,
			errorIsNil:          false,
		},
	}
	for _, scenario := range scenarios {
		t.Run(scenario.name, func(t *testing.T) {
			// Arrange:
			config, server := helpers.CreateHttpServer(
				t,
				"/customers",
				scenario.mockResponse.StatusCode,
				scenario.mockResponse.Payload,
			)
			defer server.Close()

			customersApi := apis.NewCustomers(config)
			ctx := context.Background()

			// Act:
			apiResponse, err := customersApi.CreateCustomer(ctx, scenario.args.payload)

			// Assert:
			assert.Equal(t, scenario.errorIsNil, err == nil, err)
			assert.Equal(t, scenario.expectedApiResponse, apiResponse)
		})
	}
}

// TestCreateCustomer_WithoutNames covers a customer the caller knows only by
// email address: CustomerId is the one required field, not the names.
func TestCreateCustomer_WithoutNames(t *testing.T) {
	config, server := helpers.CreateHttpServer(t, "/customers", http.StatusOK, []byte(`
		{"code": "00", "message": "customer created successfully"}
	`))
	defer server.Close()

	apiResponse, err := apis.NewCustomers(config).CreateCustomer(context.Background(), &models.CreateCustomerPayload{
		CustomerId: "CUST-12345",
		Email:      "nameless@example.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, "00", apiResponse.Code)
}
