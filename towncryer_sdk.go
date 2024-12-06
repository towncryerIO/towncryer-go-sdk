package towncryer

import (
	"github.com/volvlabs/towncryer-sdk/apis"
	"github.com/volvlabs/towncryer-sdk/config"
)

// TowncryerSdk is the main SDK struct for interacting with the Towncryer APIs. It provides access to various
// API services (e.g., Customers and Events) using lazy initialization, meaning that each service API is only
// created when first accessed. This SDK is configured using the provided `config.Configuration`.

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
		cfg: cfg,
	}
}

// CustomersApi returns an instance of the Customers API client. If it has not been initialized yet,
// it will create a new instance using the provided configuration. The Customers API client is used
// for all customer-related interactions with the Towncryer API, such as managing customer data and retrieving information.
//
// Returns:
//   - A pointer to an `apis.Customers` instance.
func (t *TowncryerSdk) CustomersApi() *apis.Customers {
	if t.customersApi == nil {
		t.customersApi = apis.NewCustomers(t.cfg)
	}
	return t.customersApi
}

// EventsApi returns an instance of the Events API client. If it has not been initialized yet,
// it will create a new instance using the provided configuration. The Events API client is used
// for managing and querying event-related data in the Towncryer API.
//
// Returns:
//   - A pointer to an `apis.Events` instance.
func (t *TowncryerSdk) EventsApi() *apis.Events {
	if t.eventsApi == nil {
		t.eventsApi = apis.NewEvents(t.cfg)
	}

	return t.eventsApi
}

func (t *TowncryerSdk) MessagesApi() *apis.Messages {
	if t.messagesApi == nil {
		t.messagesApi = apis.NewMessages(t.cfg)
	}

	return t.messagesApi
}

func (t *TowncryerSdk) Run() {
	
}
