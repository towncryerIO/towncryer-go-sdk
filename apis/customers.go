package apis

import (
	"context"

	"github.com/volvlabs/towncryer-sdk/config"
	"github.com/volvlabs/towncryer-sdk/models"
)

// CustomersAPI is the interface satisfied by Customers, the Customers API client.
// Consumers should depend on this interface rather than the concrete type so that
// a generated mock (see the mocks package) can be substituted in their own tests.
type CustomersAPI interface {
	CreateCustomer(ctx context.Context, payload *models.CreateCustomerPayload) (*models.ApiResponse[models.Customer], error)
}

type Customers struct {
	*BaseAPI
}

var _ CustomersAPI = (*Customers)(nil)

// CreateCustomer is idempotent (customers are keyed by CustomerId), so retries are
// safe and enabled by default.
func NewCustomers(config *config.Configuration) *Customers {
	return &Customers{BaseAPI: NewBaseAPI(config, "customers", true)}
}

func (c *Customers) CreateCustomer(ctx context.Context, payload *models.CreateCustomerPayload) (*models.ApiResponse[models.Customer], error) {
	if err := validatePayload(payload); err != nil {
		return nil, err
	}

	response := &models.ApiResponse[models.Customer]{}
	err := c.requestsBuilder.
		BodyJSON(payload).
		ToJSON(response).
		Post().
		Fetch(ctx)
	if err != nil {
		return nil, err
	}

	return response, nil
}
