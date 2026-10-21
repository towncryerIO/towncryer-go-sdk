package apis

import (
	"context"
	"errors"

	"github.com/towncryerIO/towncryer-go-sdk/config"
	"github.com/towncryerIO/towncryer-go-sdk/models"
)

var ErrNoMessagesSetInPayload = errors.New("message payload is totally empty")

// MessagesAPI is the interface satisfied by Messages, the Messages API client.
// Consumers should depend on this interface rather than the concrete type so that
// a generated mock (see the mocks package) can be substituted in their own tests.
type MessagesAPI interface {
	SendMessage(ctx context.Context, payload *models.SendMessagesPayload) (*models.ApiResponse[models.SendMessageResponsePayload], error)
	SendBulkMessages(ctx context.Context, payload *models.SendMessagesPayload) (*models.ScheduleInfo, error)
	SendEmail(ctx context.Context, from models.EmailFrom, to models.EmailRecipient, title string, content string) (*models.ApiResponse[models.SendMessageResponsePayload], error)
	SendSMS(ctx context.Context, from string, to models.SMSRecipient, content string) (*models.ApiResponse[models.SendMessageResponsePayload], error)
	SendPushNotification(ctx context.Context, to models.PushNotificationRecipient, title string, content string) (*models.ApiResponse[models.SendMessageResponsePayload], error)
}

type Messages struct {
	*BaseAPI
}

var _ MessagesAPI = (*Messages)(nil)

func NewMessages(config *config.Configuration) *Messages {
	m := &Messages{
		BaseAPI: NewBaseAPI(config, "messages", config.EnableRetriesForNonIdempotentRequests),
	}
	return m
}

func (m *Messages) SendMessage(
	ctx context.Context,
	payload *models.SendMessagesPayload,
) (*models.ApiResponse[models.SendMessageResponsePayload], error) {
	if err := checkMessagesPayload(payload); err != nil {
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

// SendBulkMessages sends a batch of emails, SMS, and/or push notifications in a single call
// and returns the resulting job's schedule info.
func (m *Messages) SendBulkMessages(
	ctx context.Context,
	payload *models.SendMessagesPayload,
) (*models.ScheduleInfo, error) {
	if err := checkMessagesPayload(payload); err != nil {
		return nil, err
	}

	response := &models.ScheduleInfo{}
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

func checkMessagesPayload(payload *models.SendMessagesPayload) error {
	if payload == nil || len(payload.Emails) == 0 && len(payload.PushNotifications) == 0 && len(payload.SMSes) == 0 {
		return ErrNoMessagesSetInPayload
	}
	return validatePayload(payload)
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
