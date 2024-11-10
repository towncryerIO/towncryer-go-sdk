package towncryer

import (
	"github.com/volvlabs/towncryer-sdk/apis"
	"github.com/volvlabs/towncryer-sdk/config"
)

type TowncryerSdk struct {
	cfg *config.Configuration

	customersApi *apis.Customers
	eventsApi    *apis.Events
}

func New(cfg *config.Configuration) *TowncryerSdk {
	return &TowncryerSdk{
		cfg: cfg,
	}
}

func (t *TowncryerSdk) CustomersApi() *apis.Customers {
	if t.customersApi == nil {
		t.customersApi = apis.NewCustomers(t.cfg)
	}
	return t.customersApi
}

func (t *TowncryerSdk) EventsApi() *apis.Events {
	if t.eventsApi == nil {
		t.eventsApi = apis.NewEvents(t.cfg)
	}

	return t.eventsApi
}
