package models

type PublishEventPayload struct {
	Name     string                `json:"name" validate:"required"`
	Customer CreateCustomerPayload `json:"customer"`
	Data     map[string]any        `json:"data"`
}
