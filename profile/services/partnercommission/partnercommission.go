package partnercommission

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	rc "brank.as/petnet/gunk/dsa/v2/partnercommission"
)

type Svc struct {
	rc.UnsafePartnerCommissionServiceServer
	core PartnerCommissionCore
}

type PartnerCommissionCore interface {
	CreatePartnerCommission(ctx context.Context, req *rc.CreatePartnerCommissionRequest) (*rc.CreatePartnerCommissionResponse, error)
	UpdatePartnerCommission(ctx context.Context, req *rc.UpdatePartnerCommissionRequest) (*rc.UpdatePartnerCommissionResponse, error)
	CreatePartnerCommissionTier(ctx context.Context, req *rc.CreatePartnerCommissionTierRequest) (*rc.CreatePartnerCommissionTierResponse, error)
	UpdatePartnerCommissionTier(ctx context.Context, req *rc.UpdatePartnerCommissionTierRequest) (*rc.UpdatePartnerCommissionTierResponse, error)
	GetPartnerCommissionsList(ctx context.Context, req *rc.GetPartnerCommissionsListRequest) (*rc.GetPartnerCommissionsListResponse, error)
	GetPartnerCommissionsTierList(ctx context.Context, req *rc.GetPartnerCommissionsTierListRequest) (*rc.GetPartnerCommissionsTierListResponse, error)
	DeletePartnerCommission(ctx context.Context, req *rc.DeletePartnerCommissionRequest) error
	DeletePartnerCommissionTier(ctx context.Context, req *rc.DeletePartnerCommissionTierRequest) error
	DeletePartnerCommissionTierById(ctx context.Context, req *rc.DeletePartnerCommissionTierByIdRequest) error
}

func New(core PartnerCommissionCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { rc.RegisterPartnerCommissionServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return rc.RegisterPartnerCommissionServiceHandlerFromEndpoint(ctx, mux, address, options)
}
