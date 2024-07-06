package auth

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/petnet/api/core"
	"brank.as/petnet/svcutil/mw/meta"

	authpb "brank.as/rbac/gunk/v1/authenticate"
	osapb "brank.as/rbac/gunk/v1/oauth2"
)

type Source interface {
	UserLogin(ctx context.Context, usr, pass string) (*core.User, error)
	GetUser(ctx context.Context, id string) (*core.User, error)
}

type Svc struct {
	authpb.UnimplementedSessionServiceServer
	src Source
	meta.PublicService
	cl osapb.AuthClientServiceClient
}

// New ...
func New(s Source, cl osapb.AuthClientServiceClient) *Svc {
	return &Svc{src: s, cl: cl}
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	authpb.RegisterSessionServiceServer(srv, s)
	return nil
}

// RegisterGateway ...
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) (err error) {
	return authpb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, addr, opts)
}
