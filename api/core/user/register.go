package user

import (
	"context"

	"brank.as/petnet/api/core"
	phmw "brank.as/petnet/api/perahub-middleware"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) RegisterUser(ctx context.Context, req core.RegisterUserReq) (*core.RegisterUserResp, error) {
	um, ok := s.usermanagers[phmw.GetPartner(ctx)]
	if !ok {
		return nil, status.Error(codes.NotFound, "missing user management for partner")
	}
	return um.RegisterUser(ctx, req)
}
