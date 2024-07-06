package revenuesharingreport

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	rsp "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
)

type Svc struct {
	rsp.UnsafeRevenueSharingReportServiceServer
	core RevenueSharingReportCore
}

type RevenueSharingReportCore interface {
	CreateRevenueSharingReport(ctx context.Context, req *rsp.CreateRevenueSharingReportRequest) (*rsp.CreateRevenueSharingReportResponse, error)
	GetRevenueSharingReportList(ctx context.Context, req *rsp.GetRevenueSharingReportListRequest) (*rsp.GetRevenueSharingReportListResponse, error)
	UpdateRevenueSharingReport(ctx context.Context, req *rsp.UpdateRevenueSharingReportRequest) (*rsp.UpdateRevenueSharingReportResponse, error)
}

func New(core RevenueSharingReportCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { rsp.RegisterRevenueSharingReportServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return rsp.RegisterRevenueSharingReportServiceHandlerFromEndpoint(ctx, mux, address, options)
}
