package v1

import (
	"context"
)

// CreatePaymentIntent handles POST /payment-intents (Domain Owner: Long)
// Initiates a new PaymentIntent with idempotency enforcement.
func (s *Server) CreatePaymentIntent(ctx context.Context, request CreatePaymentIntentRequestObject) (CreatePaymentIntentResponseObject, error) {
	detail := "Payments & Core Engine PaymentIntent processing is currently under active development (Owner: Long)."
	instance := "/payment-intents"
	return CreatePaymentIntent501ApplicationProblemPlusJSONResponse(newNotImplementedProblem(instance, detail)), nil
}
