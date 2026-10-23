package v1

import (
	"context"
	"time"
)

// GetHealthz handles GET /healthz (Observability & liveness probe)
func (s *Server) GetHealthz(ctx context.Context, request GetHealthzRequestObject) (GetHealthzResponseObject, error) {
	return GetHealthz200JSONResponse{
		Status:    Ok,
		Timestamp: time.Now().UTC(),
		Version:   s.version,
	}, nil
}
