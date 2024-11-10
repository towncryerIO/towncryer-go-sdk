package models

import "github.com/google/uuid"

type Customer struct {
	Id                     uuid.UUID `json:"id"`
	CustomerId             string    `json:"customerId"`
	FirstName              string    `json:"firstName"`
	LastName               string    `json:"lastName"`
	Emails                 []string  `json:"emails"`
	PhoneNumbers           []string  `json:"phoneNumbers"`
	PushNotificationTokens []string  `json:"pushNotificationTokens"`
}
