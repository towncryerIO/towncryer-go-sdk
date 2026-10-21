package towncryer

import (
	"github.com/volvlabs/towncryer-sdk/apis"
	"github.com/volvlabs/towncryer-sdk/config"
)

// TowncryerSdk is the main SDK struct for interacting with the Towncryer APIs. It provides access to various
// API services (e.g., Customers and Events), which are constructed eagerly in New() so that they are safe
// to access concurrently. This SDK is configured using the provided `config.Configuration`.

type TowncryerSdk struct {
	cfg *config.Configuration

	customersApi *apis.Customers
	eventsApi    *apis.Events
	messagesApi  *apis.Messages
}

// New initializes a new instance of the TowncryerSdk with the given configuration.
// This configuration contains essential information such as API credentials and the API base URL.
//
// Parameters:
//   - cfg: A pointer to a `config.Configuration` struct that holds the SDK's configuration details.
//
// Returns:
//   - A pointer to a newly created `TowncryerSdk` instance.
func New(cfg *config.Configuration) *TowncryerSdk {
	return &TowncryerSdk{
		cfg:          cfg,
		customersApi: apis.NewCustomers(cfg),
		eventsApi:    apis.NewEvents(cfg),
		messagesApi:  apis.NewMessages(cfg),
	}
}

// CustomersApi returns an instance of the Customers API client. The Customers API client is used
// for all customer-related interactions with the Towncryer API, such as managing customer data and retrieving information.
//
// Returns:
//   - A pointer to an `apis.Customers` instance.
func (t *TowncryerSdk) CustomersApi() *apis.Customers {
	return t.customersApi
}

// EventsApi returns an instance of the Events API client. The Events API client is used
// for managing and querying event-related data in the Towncryer API.
//
// Returns:
//   - A pointer to an `apis.Events` instance.
func (t *TowncryerSdk) EventsApi() *apis.Events {
	return t.eventsApi
}

// MessagesApi returns an instance of the Messages API client, used for sending emails, SMS,
// and push notifications through the Towncryer API.
//
// Returns:
//   - A pointer to an `apis.Messages` instance.
func (t *TowncryerSdk) MessagesApi() *apis.Messages {
	return t.messagesApi
}
