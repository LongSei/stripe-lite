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

// ListWebhookEndpoints handles GET /webhook-endpoints
func (s *Server) ListWebhookEndpoints(ctx context.Context, request ListWebhookEndpointsRequestObject) (ListWebhookEndpointsResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)"
	instance := "/webhook-endpoints"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}

// CreateWebhookEndpoint handles POST /webhook-endpoints
// Configures HTTPS delivery target and event subscriptions with signing secrets.
func (s *Server) CreateWebhookEndpoint(ctx context.Context, request CreateWebhookEndpointRequestObject) (CreateWebhookEndpointResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)"
	instance := "/webhook-endpoints"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}

// Response visitors for the remaining Webhook operations.
// All return HTTP 501 with an RFC 7807 Problem Document.

func (response webhookNotImplementedResponse) VisitGetWebhookEndpointResponse(w http.ResponseWriter) error {
	return response.VisitListWebhookEndpointsResponse(w)
}

func (response webhookNotImplementedResponse) VisitUpdateWebhookEndpointResponse(w http.ResponseWriter) error {
	return response.VisitListWebhookEndpointsResponse(w)
}

func (response webhookNotImplementedResponse) VisitDeleteWebhookEndpointResponse(w http.ResponseWriter) error {
	return response.VisitListWebhookEndpointsResponse(w)
}

func (response webhookNotImplementedResponse) VisitRotateWebhookSecretResponse(w http.ResponseWriter) error {
	return response.VisitListWebhookEndpointsResponse(w)
}

// GetWebhookEndpoint handles GET /webhook-endpoints/{id}
func (s *Server) GetWebhookEndpoint(ctx context.Context, request GetWebhookEndpointRequestObject) (GetWebhookEndpointResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)"
	instance := "/webhook-endpoints/{id}"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}

// UpdateWebhookEndpoint handles PATCH /webhook-endpoints/{id}
func (s *Server) UpdateWebhookEndpoint(ctx context.Context, request UpdateWebhookEndpointRequestObject) (UpdateWebhookEndpointResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)"
	instance := "/webhook-endpoints/{id}"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}

// DeleteWebhookEndpoint handles DELETE /webhook-endpoints/{id}
func (s *Server) DeleteWebhookEndpoint(ctx context.Context, request DeleteWebhookEndpointRequestObject) (DeleteWebhookEndpointResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)"
	instance := "/webhook-endpoints/{id}"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}

// RotateWebhookSecret handles POST /webhook-endpoints/{id}/rotate-secret
func (s *Server) RotateWebhookSecret(ctx context.Context, request RotateWebhookSecretRequestObject) (RotateWebhookSecretResponseObject, error) {
	detail := "Webhooks management and subscription dispatcher is currently under active development (Owner: Sơn)"
	instance := "/webhook-endpoints/{id}/rotate-secret"
	return webhookNotImplementedResponse(newNotImplementedProblem(instance, detail)), nil
}
