package v1

import (
	"context"
	"encoding/json"
	"net/http"
)

type webhookNotImplementedResponse ProblemDocument

func (response webhookNotImplementedResponse) VisitListWebhookEndpointsResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusNotImplemented)

	return json.NewEncoder(w).Encode(response)
}

func (response webhookNotImplementedResponse) VisitCreateWebhookEndpointResponse(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(http.StatusNotImplemented)

	return json.NewEncoder(w).Encode(response)
}

// ListWebhookEndpoints handles GET /webhook-endpoints (Domain Owner: Sơn)
// Lists webhook delivery targets without implementing persistence yet.
func (s *Server) ListWebhookEndpoints(ctx context.Context, request ListWebhookEndpointsRequestObject) (ListWebhookEndpointsResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)."
	instance := "/webhook-endpoints"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}

// CreateWebhookEndpoint handles POST /webhook-endpoints (Domain Owner: Sơn)
// Configures HTTPS delivery target and event subscriptions with signing secrets.
func (s *Server) CreateWebhookEndpoint(ctx context.Context, request CreateWebhookEndpointRequestObject) (CreateWebhookEndpointResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)."
	instance := "/webhook-endpoints"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}
