package v1

import (
	"context"
)

// CreateCharge handles POST /payments/charges (Domain Owner: Long)
// Atomically executes credit charge authorizations with idempotency key enforcement.
func (s *Server) CreateCharge(ctx context.Context, request CreateChargeRequestObject) (CreateChargeResponseObject, error) {
	detail := "Payments & Core Engine charge processing is currently under active development (Owner: Long)."
	instance := "/payments/charges"
	return CreateCharge501ApplicationProblemPlusJSONResponse(newNotImplementedProblem(instance, detail)), nil
}
