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

func TestSendMessages(t *testing.T) {
	type args struct {
		payload *models.SendMessagesPayload
	}
	scenarios := []struct {
		name                string
		args                args
		mockResponse        helpers.MockResponse
		expectedApiResponse *models.ApiResponse[models.SendMessageResponsePayload]
		errorIsNil          bool
	}{
		{
			name: "should return error because no messages were set in the payload at all",
			args: args{
				payload: &models.SendMessagesPayload{},
			},
			mockResponse:        helpers.MockResponse{},
			expectedApiResponse: nil,
			errorIsNil:          false,
		},
		{
			name: "should return successfully send message",
			args: args{
				payload: &models.SendMessagesPayload{
					Emails: []models.SendEmailPayload{
						{
							Title:   "New Transaction",
							Content: "hello {{firstName}}, we received your transaction",
							From: models.EmailFrom{
								Name:  "John & Sons",
								Email: "no-reply@company.com",
							},
							Recipients: []models.EmailRecipient{
								{
									Name:    "John Doe",
									Email:   "john.doe@gmail.com",
									Context: []byte(`{"firstName": "John"}`),
								},
							},
						},
					},
					PushNotifications: []models.SendPushNotificationPayload{
						{
							Title:   "New Transaction",
							Content: "{{firstName}}, we recieved your new transaction",
							Recipients: []models.PushNotificationRecipient{
								{
									Token:   "push-notification-token",
									Context: []byte(`{"firstName": "John"}`),
								},
							},
						},
					},
					SMSes: []models.SendSMSPayload{
						{
							Content: "Your OTP is {{otp}}",
							From:    "2348091607289",
							Recipients: []models.SMSRecipient{
								{
									PhoneNumber: "+2348091607399",
									Context:     []byte(`{"otp": 123456}`),
								},
								{
									PhoneNumber: "+2348091607392",
									Context:     []byte(`{"otp": 123455}`),
								},
							},
						},
					},
				},
			},
			mockResponse: helpers.MockResponse{
				StatusCode: http.StatusOK,
				Payload: []byte(`{
					"code": "00",
					"message": "success",
					"data": {
						"id": "550e8400-e29b-41d4-a716-446655440000",
						"status": "pending"
					}
				}`),
			},
			expectedApiResponse: &models.ApiResponse[models.SendMessageResponsePayload]{
				Code:    "00",
				Message: "success",
				Data: models.SendMessageResponsePayload{
					Id:     "550e8400-e29b-41d4-a716-446655440000",
					Status: "pending",
				},
			},
			errorIsNil: true,
		},
		{
			name: "should return error because payload is invalid",
			args: args{
				payload: &models.SendMessagesPayload{
					Emails: []models.SendEmailPayload{
						{
							Title:   "New Transaction",
							Content: "hello {{firstName}}, we received your transaction",
							From: models.EmailFrom{
								Name:  "John & Sons",
								Email: "no-reply@company.com",
							},
							Recipients: []models.EmailRecipient{
								{
									Name:    "John Doe",
									Email:   "john.doe@gmail",
									Context: []byte(`{"firstName": "John"}`),
								},
							},
						},
					},
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
				"/messages",
				scenario.mockResponse.StatusCode,
				scenario.mockResponse.Payload,
			)
			defer server.Close()

			messagesApi := apis.NewMessages(config)
			ctx := context.Background()

			// Act:
			apiResponse, err := messagesApi.SendMessage(ctx, scenario.args.payload)

			// Assert:
			assert.Equal(t, scenario.errorIsNil, err == nil, err)
			assert.Equal(t, scenario.expectedApiResponse, apiResponse)
		})
	}
}
