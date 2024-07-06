package cicopartnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Svc struct {
	spb.UnimplementedCICOPartnerListServiceServer
	core CICOPartnerListCore
}

type CICOPartnerListCore interface {
	CreateCICOPartnerList(context.Context, *spb.CreateCICOPartnerListRequest) (*spb.CreateCICOPartnerListResponse, error)
	UpdateCICOPartnerList(context.Context, *spb.UpdateCICOPartnerListRequest) (*spb.UpdateCICOPartnerListResponse, error)
	GetCICOPartnerList(context.Context, *spb.GetCICOPartnerListRequest) (*spb.GetCICOPartnerListResponse, error)
	DeleteCICOPartnerList(context.Context, *spb.DeleteCICOPartnerListRequest) (*spb.DeleteCICOPartnerListResponse, error)
	EnableCICOPartnerList(context.Context, *spb.EnableCICOPartnerListRequest) (*spb.EnableCICOPartnerListResponse, error)
	DisableCICOPartnerList(context.Context, *spb.DisableCICOPartnerListRequest) (*spb.DisableCICOPartnerListResponse, error)
}

func New(core CICOPartnerListCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterPartner with grpc server.
func (s *Svc) Register(srv *grpc.Server) { spb.RegisterCICOPartnerListServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return spb.RegisterCICOPartnerListServiceHandlerFromEndpoint(ctx, mux, address, options)
}
