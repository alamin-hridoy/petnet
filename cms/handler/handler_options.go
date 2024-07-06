package handler

import "context"

type iRemcoCommissionSvc interface {
	SyncRemcoCommissionConfigForRemittance(ctx context.Context)
	SyncDSACommissionConfigForRemittance(ctx context.Context, orgID string)
}

// ServerOptions is type for creating Server with options
type ServerOptions func(s *Server)

// WithRemcoCommissionSvc ...
func WithRemcoCommissionSvc(rcs iRemcoCommissionSvc) ServerOptions {
	return func(s *Server) {
		s.remcoCommSvc = rcs
	}
}
