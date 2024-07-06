package validation

import (
	"context"

	"google.golang.org/grpc"

	ppb "brank.as/rbac/gunk/v1/permissions"
	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

type Svc struct {
	v  ppb.ValidationServiceServer
	sa sapb.ValidationServiceServer
}

func NewLocal(v ppb.ValidationServiceServer, sa sapb.ValidationServiceServer) *Svc {
	return &Svc{v: v, sa: sa}
}

func (s *Svc) ValidatePermission(ctx context.Context, in *ppb.ValidatePermissionRequest, _ ...grpc.CallOption) (*ppb.ValidatePermissionResponse, error) {
	return s.v.ValidatePermission(ctx, in)
}

func (s *Svc) ValidateAccount(ctx context.Context, in *sapb.ValidateAccountRequest, _ ...grpc.CallOption) (*sapb.ValidateAccountResponse, error) {
	return s.sa.ValidateAccount(ctx, in)
}
