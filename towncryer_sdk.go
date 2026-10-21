package towncryer

import (
	"github.com/towncryerIO/towncryer-go-sdk/apis"
	"github.com/towncryerIO/towncryer-go-sdk/config"
)

// TowncryerSdk is the main SDK struct for interacting with the Towncryer APIs. It exposes each API
// service (e.g., Customers, Events and Messages) as a public field, constructed eagerly in New() so that
// they are safe to access concurrently. This SDK is configured using the provided `config.Configuration`.

type TowncryerSdk struct {
	cfg *config.Configuration

	Customers apis.CustomersAPI
	Events    apis.EventsAPI
	Messages  apis.MessagesAPI
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
	resolved := cfg.WithEnvDefaults()
	return &TowncryerSdk{
		cfg:       &resolved,
		Customers: apis.NewCustomers(&resolved),
		Events:    apis.NewEvents(&resolved),
		Messages:  apis.NewMessages(&resolved),
	}
}
