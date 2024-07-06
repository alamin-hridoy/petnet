package profile

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/petnet/profile/storage"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/profile/core/profile"
	uspb "brank.as/rbac/gunk/v1/user"
)

type Svc struct {
	ppb.UnimplementedOrgProfileServiceServer
	ps   ProfileStore
	cl   uspb.UserServiceClient
	core *profile.Svc
}

type ProfileStore interface {
	CreateOrgProfile(ctx context.Context, p *storage.OrgProfile) (string, error)
	UpdateOrgProfile(ctx context.Context, p *storage.OrgProfile) (string, error)
	GetOrgProfile(ctx context.Context, id string) (*storage.OrgProfile, error)
	GetOrgProfiles(ctx context.Context, f storage.FilterList) ([]storage.OrgProfile, error)
	GetProfileByDsaCode(ctx context.Context, dsaCode string) (*storage.OrgProfile, error)
}

func New(ps ProfileStore, cl uspb.UserServiceClient, core *profile.Svc) *Svc {
	h := &Svc{
		ps:   ps,
		cl:   cl,
		core: core,
	}
	return h
}

// RegisterService with grpc server.
func (h *Svc) Register(srv *grpc.Server) { ppb.RegisterOrgProfileServiceServer(srv, h) }

// RegisterGateway grpcgw
func (h *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterOrgProfileServiceHandlerFromEndpoint(ctx, mux, address, options)
}
