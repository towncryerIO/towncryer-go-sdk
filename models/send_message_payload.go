package models

import "encoding/json"

type EmailRecipient struct {
	Name    string          `json:"name" validate:"required"`
	Email   string          `json:"email" validate:"required,email"`
	Context json.RawMessage `json:"context"`
}

type EmailFrom struct {
	Name  string `json:"name" validate:"required"`
	Email string `json:"email" validate:"required,email"`
}

type SendEmailPayload struct {
	Title      string           `json:"title" validate:"required_without=TemplateId"`
	TemplateId string           `json:"templateId" validate:"required_without=Content"`
	Content    string           `json:"content" validate:"required_without=TemplateId"`
	From       EmailFrom        `json:"from" validate:"required_without=TemplateId"`
	Recipients []EmailRecipient `json:"recipients" validate:"required,dive"`
}

type PushNotificationRecipient struct {
	Token   string          `json:"token" validate:"required"`
	Context json.RawMessage `json:"context"`
}

type SendPushNotificationPayload struct {
	Title      string                      `json:"title" validate:"required_without=TemplateId"`
	TemplateId string                      `json:"templateId" validate:"required_without=Content"`
	Content    string                      `json:"content" validate:"required_without=TemplateId"`
	Recipients []PushNotificationRecipient `json:"recipients" validate:"required,dive"`
	Data       json.RawMessage             `json:"data"`
}

type SMSRecipient struct {
	PhoneNumber string          `json:"phoneNumber" validate:"required"`
	Context     json.RawMessage `json:"context"`
}

type SendSMSPayload struct {
	TemplateId string         `json:"templateId" validate:"required_without=Content"`
	Content    string         `json:"content" validate:"required_without=TemplateId"`
	From       string         `json:"from" validate:"required"`
	Recipients []SMSRecipient `json:"to" validate:"required,dive"`
}

type SendMessagesPayload struct {
	Emails            []SendEmailPayload            `json:"emails" validate:"omitempty,dive,required_without_all=PushNotifications SMSes"`
	PushNotifications []SendPushNotificationPayload `json:"pushNotifications" validate:"omitempty,dive,required_without_all=Emails SMSes"`
	SMSes             []SendSMSPayload              `json:"smses" validate:"omitempty,dive,required_without_all=Emails PushNotifications"`
}
