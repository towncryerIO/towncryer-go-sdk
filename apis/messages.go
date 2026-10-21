package apis

import (
	"context"
	"errors"

	"github.com/volvlabs/towncryer-sdk/config"
	"github.com/volvlabs/towncryer-sdk/models"
)

var ErrNoMessagesSetInPayload = errors.New("message payload is totally empty")

type Messages struct {
	*BaseAPI
}

func NewMessages(config *config.Configuration) *Messages {
	m := &Messages{
		BaseAPI: NewBaseAPI(config, "messages"),
	}
	return m
}

func (m *Messages) SendMessage(
	ctx context.Context,
	payload *models.SendMessagesPayload,
) (*models.ApiResponse[models.SendMessageResponsePayload], error) {
	if len(payload.Emails) == 0 && len(payload.PushNotifications) == 0 ||
		len(payload.PushNotifications) == 0 {
		return nil, ErrNoMessagesSetInPayload
	}
	if err := validatePayload(payload); err != nil {
		return nil, err
	}

	response := &models.ApiResponse[models.SendMessageResponsePayload]{}
	err := m.requestsBuilder.
		BodyJSON(payload).
		ToJSON(response).
		Post().
		Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (m *Messages) SendEmail(
	ctx context.Context,
	from models.EmailFrom,
	to models.EmailRecipient,
	title string,
	content string,
) (*models.ApiResponse[models.SendMessageResponsePayload], error) {
	payload := &models.SendMessagesPayload{
		Emails: []models.SendEmailPayload{
			{
				From:       from,
				Recipients: []models.EmailRecipient{to},
				Title:      title,
				Content:    content,
			},
		},
	}
	return m.SendMessage(ctx, payload)
}

func (m *Messages) SendSMS(
	ctx context.Context,
	from string,
	to models.SMSRecipient,
	content string,
) (*models.ApiResponse[models.SendMessageResponsePayload], error) {
	payload := &models.SendMessagesPayload{
		SMSes: []models.SendSMSPayload{
			{
				From: from,
				Recipients: []models.SMSRecipient{
					to,
				},
				Content: content,
			},
		},
	}
	return m.SendMessage(ctx, payload)
}

func (m *Messages) SendPushNotification(
	ctx context.Context,
	to models.PushNotificationRecipient,
	title string,
	content string,
) (*models.ApiResponse[models.SendMessageResponsePayload], error) {
	payload := &models.SendMessagesPayload{
		PushNotifications: []models.SendPushNotificationPayload{
			{
				Recipients: []models.PushNotificationRecipient{to},
				Title:      title,
				Content:    content,
			},
		},
	}
	return m.SendMessage(ctx, payload)
}
