package validation

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

type Validator interface {
	Validate(context.Context, core.Validation) (*core.Identity, error)
}

type Svc struct {
	v Validator
	ppb.UnimplementedValidationServiceServer
}

func New(v Validator) *Svc {
	return &Svc{
		v: v,
	}
}

func (s *Svc) RegisterSvc(svr *grpc.Server) error {
	ppb.RegisterValidationServiceServer(svr, s)
	return nil
}

func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterValidationServiceHandlerFromEndpoint(ctx, mux, address, options)
}
