package mfa

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	mpb "brank.as/petnet/gunk/v1/mfa"
	core "brank.as/petnet/profile/core/mfa"
)

type Svc struct {
	mpb.UnimplementedMFAServiceServer
	core MFACore
}

type MFACore interface {
	EnableMFA(context.Context, core.EnableMFAReq) (*core.EnableMFAResp, error)
}

func New(core MFACore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { mpb.RegisterMFAServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return mpb.RegisterMFAServiceHandlerFromEndpoint(ctx, mux, address, options)
}
