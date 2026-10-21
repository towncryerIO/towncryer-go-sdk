package apis

import (
	"context"

	"github.com/volvlabs/towncryer-sdk/config"
	"github.com/volvlabs/towncryer-sdk/models"
)

type Events struct {
	*BaseAPI
}

// PublishEvent is not naturally idempotent: a retried request can cause the server to
// publish the same event twice. Retries are therefore only applied for this endpoint if
// the caller has explicitly opted in via config.EnableRetriesForNonIdempotentRequests.
func NewEvents(config *config.Configuration) *Events {
	return &Events{BaseAPI: NewBaseAPI(config, "events", config.EnableRetriesForNonIdempotentRequests)}
}

func (e *Events) PublishEvent(ctx context.Context, payload *models.PublishEventPayload) (*models.ApiResponse[map[string]any], error) {
	if err := validatePayload(payload); err != nil {
		return nil, err
	}

	response := &models.ApiResponse[map[string]any]{}
	err := e.requestsBuilder.
		BodyJSON(payload).
		ToJSON(response).
		Post().
		Fetch(ctx)
	if err != nil {
		return nil, err
	}
	return response, nil
}
