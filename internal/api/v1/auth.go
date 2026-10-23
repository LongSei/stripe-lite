package v1

import (
	"context"
)

// CreateAuthToken handles POST /auth/tokens (Domain Owner: Khánh)
// Exchanges merchant credentials for scoped JWT session tokens.
func (s *Server) CreateAuthToken(ctx context.Context, request CreateAuthTokenRequestObject) (CreateAuthTokenResponseObject, error) {
	detail := "Authentication & Identity domain handler is currently under active development (Owner: Khánh)."
	instance := "/auth/tokens"
	return CreateAuthToken501ApplicationProblemPlusJSONResponse(newNotImplementedProblem(instance, detail)), nil
}
