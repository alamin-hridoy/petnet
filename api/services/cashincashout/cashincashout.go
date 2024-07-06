package cashincashout

import (
	"context"

	"google.golang.org/grpc"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	cico "brank.as/petnet/gunk/drp/v1/cashincashout"
)

type CicoCore interface {
	CiCoInquire(ctx context.Context, req *cico.CiCoInquireRequest) (*cico.CiCoInquireResponse, error)
	CiCoExecute(ctx context.Context, req *cico.CiCoExecuteRequest) (*cico.CiCoExecuteResponse, error)
	CiCoRetry(ctx context.Context, req *cico.CiCoRetryRequest) (*cico.CiCoRetryResponse, error)
	CiCoOTPConfirm(ctx context.Context, req *cico.CiCoOTPConfirmRequest) (*cico.CiCoOTPConfirmResponse, error)
	CiCoValidate(ctx context.Context, req *cico.CiCoValidateRequest) (*cico.CiCoValidateResponse, error)
	CICOTransactList(ctx context.Context, req *cico.CICOTransactListRequest) (*cico.CICOTransactListResponse, error)
}

type Svc struct {
	cico.UnimplementedCashInCashOutServiceServer
	core CicoCore
}

func New(cio CicoCore) *Svc {
	s := &Svc{
		core: cio,
	}
	return s
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	cico.RegisterCashInCashOutServiceServer(srv, s)
	return nil
}

// RegisterGateway parter endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return cico.RegisterCashInCashOutServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
