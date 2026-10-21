package models

// CreateCustomerPayload identifies a customer and the contact details the API
// can reach them on. Only CustomerId is required: a caller may know a customer
// by an email address or phone number alone, with no name to give.
type CreateCustomerPayload struct {
	CustomerId            string `json:"customerId" validate:"required"`
	FirstName             string `json:"firstName"`
	LastName              string `json:"lastName"`
	Email                 string `json:"email"`
	PhoneNumber           string `json:"phoneNumber"`
	PushNotificationToken string `json:"pushNotificationToken"`
}
