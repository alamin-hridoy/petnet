package partner

import (
	"context"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
	"brank.as/petnet/profile/storage"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Svc struct {
	ppb.UnimplementedPartnerServiceServer
	core PartnerCore
}

type PartnerCore interface {
	CreatePartners(context.Context, *ppb.Partners) error
	UpdatePartners(context.Context, *ppb.Partners) error
	GetPartners(context.Context, string) (*ppb.Partners, error)
	DeletePartner(context.Context, string) error
	ValidatePartnerAccess(ctx context.Context, oid string, pnr string) error
	EnablePartner(ctx context.Context, oid string, pnr string) error
	DisablePartner(ctx context.Context, oid string, pnr string) error
	GetPartner(ctx context.Context, oid string, tp string) (*storage.Partner, error)
}

func New(core PartnerCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterPartner with grpc server.
func (s *Svc) Register(srv *grpc.Server) { ppb.RegisterPartnerServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterPartnerServiceHandlerFromEndpoint(ctx, mux, address, options)
}
