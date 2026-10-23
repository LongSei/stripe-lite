package v1

import (
	"net/http"
)

// Server implements StrictServerInterface for Stripe-lite v1 API gateway.
type Server struct {
	version string
}

// NewServer constructs a new baseline Server instance.
func NewServer(version string) *Server {
	return &Server{
		version: version,
	}
}

var _ StrictServerInterface = (*Server)(nil)

// newNotImplementedProblem builds a standardized RFC 7807 problem document for feature stubs.
func newNotImplementedProblem(instance, detail string) ProblemDocument {
	return ProblemDocument{
		Type:     "https://api.stripe-lite.internal/errors/not-implemented",
		Title:    "Feature Not Implemented",
		Status:   http.StatusNotImplemented,
		Detail:   &detail,
		Instance: &instance,
	}
}
