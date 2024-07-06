package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Svc struct {
	spb.UnimplementedPartnerListServiceServer
	core PartnerListCore
}

type PartnerListCore interface {
	CreatePartnerList(context.Context, *spb.CreatePartnerListRequest) (*spb.CreatePartnerListResponse, error)
	UpdatePartnerList(context.Context, *spb.UpdatePartnerListRequest) (*spb.UpdatePartnerListResponse, error)
	GetPartnerList(context.Context, *spb.GetPartnerListRequest) (*spb.GetPartnerListResponse, error)
	DeletePartnerList(context.Context, *spb.DeletePartnerListRequest) (*spb.DeletePartnerListResponse, error)
	EnablePartnerList(context.Context, *spb.EnablePartnerListRequest) (*spb.EnablePartnerListResponse, error)
	DisablePartnerList(context.Context, *spb.DisablePartnerListRequest) (*spb.DisablePartnerListResponse, error)
	EnableMultiplePartnerList(context.Context, *spb.EnableMultiplePartnerListRequest) (*spb.EnableMultiplePartnerListResponse, error)
	DisableMultiplePartnerList(context.Context, *spb.DisableMultiplePartnerListRequest) (*spb.DisableMultiplePartnerListResponse, error)
	GetDSAPartnerList(context.Context, *spb.DSAPartnerListRequest) (*spb.GetDSAPartnerListResponse, error)
	GetPartnerByStype(context.Context, *spb.GetPartnerByStypeRequest) (*spb.GetPartnerByStypeResponse, error)
}

func New(core PartnerListCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterPartner with grpc server.
func (s *Svc) Register(srv *grpc.Server) { spb.RegisterPartnerListServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return spb.RegisterPartnerListServiceHandlerFromEndpoint(ctx, mux, address, options)
}
