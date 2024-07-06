package fees

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/petnet/profile/storage"

	fpb "brank.as/petnet/gunk/dsa/v2/fees"
)

type Svc struct {
	fpb.UnimplementedOrgFeesServiceServer
	core FeeCore
}

type FeeCore interface {
	UpsertFee(ctx context.Context, f storage.FeeCommission) (*storage.FeeCommission, error)
	ListFees(ctx context.Context, oid string, f storage.LimitOffsetFilter) ([]storage.FeeCommission, error)
	UpsertRate(ctx context.Context, f storage.Rate) (*storage.Rate, error)
	ListRates(ctx context.Context, fcid string) ([]storage.Rate, error)
}

func New(core FeeCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { fpb.RegisterOrgFeesServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return fpb.RegisterOrgFeesServiceHandlerFromEndpoint(ctx, mux, address, options)
}
