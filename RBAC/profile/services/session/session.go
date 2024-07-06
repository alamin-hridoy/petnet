package session

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/profile/core"

	authpb "brank.as/rbac/gunk/v1/authenticate"
)

type Svc struct {
	authpb.UnimplementedSessionServiceServer
	be Auth
}

type Auth interface {
	UserLogin(ctx context.Context, username, pass string) (*core.UserSession, error)
}

func New(be Auth) *Svc {
	return &Svc{be: be}
}

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	authpb.RegisterSessionServiceServer(srv, s)
	return nil
}

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return authpb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, address, options)
}
