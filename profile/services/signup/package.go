package signup

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	upb "brank.as/petnet/gunk/dsa/v1/user"
	core "brank.as/petnet/profile/core/signup"
)

type Svc struct {
	upb.UnimplementedSignupServiceServer
	core SignupCore
}

type SignupCore interface {
	Signup(context.Context, core.SignupReq) (*core.SignupResp, error)
	RetrieveInvite(ctx context.Context, code string) (*upb.RetrieveInviteResponse, error)
}

func New(core SignupCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { upb.RegisterSignupServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return upb.RegisterSignupServiceHandlerFromEndpoint(ctx, mux, address, options)
}
