package userprofile

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/petnet/profile/storage"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

type Svc struct {
	ppb.UnimplementedUserProfileServiceServer
	core ProfileStore
}

type ProfileStore interface {
	CreateUserProfile(ctx context.Context, p storage.UserProfile) (string, error)
	UpdateUserProfile(ctx context.Context, p storage.UserProfile) (string, error)
	UpdateUserProfileByOrgID(ctx context.Context, p storage.UpdateOrgProfileOrgIDUserID) (string, error)
	GetUserProfile(ctx context.Context, uid string) (*storage.UserProfile, error)
	GetUserProfiles(ctx context.Context, oid string) ([]storage.UserProfile, error)
	DeleteUserProfile(ctx context.Context, req *ppb.DeleteUserProfileRequest) (ppb.DeleteUserProfileResponse, error)
	EnableUserProfile(ctx context.Context, req *ppb.EnableUserProfileRequest) (*ppb.EnableUserProfileResponse, error)
	GetUserProfileByEmail(ctx context.Context, uid string) (*storage.UserProfile, error)
}

func New(core ProfileStore) *Svc {
	h := &Svc{
		core: core,
	}
	return h
}

// RegisterService with grpc server.
func (h *Svc) Register(srv *grpc.Server) { ppb.RegisterUserProfileServiceServer(srv, h) }

// RegisterGateway grpcgw
func (h *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterUserProfileServiceHandlerFromEndpoint(ctx, mux, address, options)
}
