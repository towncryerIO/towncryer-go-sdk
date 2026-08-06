package apis

import (
	"context"

	"github.com/volvlabs/towncryer-sdk/config"
	"github.com/volvlabs/towncryer-sdk/models"
)

type Customers struct {
	*BaseAPI
}

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
