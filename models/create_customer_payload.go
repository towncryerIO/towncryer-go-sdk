package models

type CreateCustomerPayload struct {
	CustomerId            string `json:"customerId" validate:"required"`
	FirstName             string `json:"firstName" validate:"required"`
	LastName              string `json:"lastName" validate:"required"`
	Email                 string `json:"email"`
	PhoneNumber           string `json:"phoneNumber"`
	PushNotificationToken string `json:"pushNotificationToken"`
}
