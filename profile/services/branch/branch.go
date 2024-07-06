package branch

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/petnet/profile/storage"

	bpb "brank.as/petnet/gunk/dsa/v2/branch"
)

type Svc struct {
	bpb.UnimplementedBranchServiceServer
	core BranchCore
}

type BranchCore interface {
	UpsertBranch(context.Context, storage.Branch) (*storage.Branch, error)
	ListBranches(ctx context.Context, org string, lim, off int, title string) ([]storage.Branch, error)
}

func New(core BranchCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { bpb.RegisterBranchServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return bpb.RegisterBranchServiceHandlerFromEndpoint(ctx, mux, address, options)
}
