package v1

import (
	"context"
)

// CreateWebhookEndpoint handles POST /webhook-endpoints (Domain Owner: Sơn)
// Configures HTTPS delivery target and event subscriptions with signing secrets.
func (s *Server) CreateWebhookEndpoint(ctx context.Context, request CreateWebhookEndpointRequestObject) (CreateWebhookEndpointResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)."
	instance := "/webhook-endpoints"
	return CreateWebhookEndpoint501ApplicationProblemPlusJSONResponse(newNotImplementedProblem(instance, detail)), nil
}
