package revenuesharing

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
)

type Svc struct {
	rc.UnimplementedRevenueSharingServiceServer
	core RevenueSharingCore
}

type RevenueSharingCore interface {
	CreateRevenueSharing(ctx context.Context, req *rc.CreateRevenueSharingRequest) (*rc.CreateRevenueSharingResponse, error)
	UpdateRevenueSharing(ctx context.Context, req *rc.UpdateRevenueSharingRequest) (*rc.UpdateRevenueSharingResponse, error)
	CreateRevenueSharingTier(ctx context.Context, req *rc.CreateRevenueSharingTierRequest) (*rc.CreateRevenueSharingTierResponse, error)
	UpdateRevenueSharingTier(ctx context.Context, req *rc.UpdateRevenueSharingTierRequest) (*rc.UpdateRevenueSharingTierResponse, error)
	GetRevenueSharingList(ctx context.Context, req *rc.GetRevenueSharingListRequest) (*rc.GetRevenueSharingListResponse, error)
	GetRevenueSharingTierList(ctx context.Context, req *rc.GetRevenueSharingTierListRequest) (*rc.GetRevenueSharingTierListResponse, error)
	DeleteRevenueSharing(ctx context.Context, req *rc.DeleteRevenueSharingRequest) error
	DeleteRevenueSharingTier(ctx context.Context, req *rc.DeleteRevenueSharingTierRequest) error
	DeleteRevenueSharingTierById(ctx context.Context, req *rc.DeleteRevenueSharingTierByIdRequest) error
	GetPartnerTransactionType(ctx context.Context, req *rc.GetPartnerTransactionTypeRequest) (*rc.GetPartnerTransactionTypeResponse, error)
}

func New(core RevenueSharingCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { rc.RegisterRevenueSharingServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return rc.RegisterRevenueSharingServiceHandlerFromEndpoint(ctx, mux, address, options)
}
