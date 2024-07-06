package rbsession

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/petnet/profile/storage"

	authpb "brank.as/rbac/gunk/v1/authenticate"
	rbupb "brank.as/rbac/gunk/v1/user"
)

type Svc struct {
	authpb.UnimplementedSessionServiceServer
	acl rbupb.UserAuthServiceClient
	ucl rbupb.UserServiceClient
	scl authpb.SessionServiceClient

	core ProfileStore
}

type ProfileStore interface {
	GetOrgProfile(ctx context.Context, id string) (*storage.OrgProfile, error)
	GetOrgProfiles(ctx context.Context) ([]storage.OrgProfile, error)
	GetUserProfileByEmail(ctx context.Context, email string) (*storage.UserProfile, error)
	SessionExists(ctx context.Context, uid string) bool
}

func New(core ProfileStore, acl rbupb.UserAuthServiceClient, ucl rbupb.UserServiceClient, scl authpb.SessionServiceClient) *Svc {
	return &Svc{
		acl:  acl,
		ucl:  ucl,
		scl:  scl,
		core: core,
	}
}

// RegisterGateway grpcgw
func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return authpb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, address, options)
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { authpb.RegisterSessionServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return authpb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, address, options)
}
