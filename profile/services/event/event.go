package event

import (
	"context"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/profile/storage"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Svc struct {
	tpb.UnimplementedEventServiceServer
	core EventCore
}

type EventCore interface {
	CreateEventData(context.Context, *storage.EventData) error
	GetEventData(context.Context, string) (*storage.EventData, error)
}

func New(core EventCore) *Svc {
	return &Svc{
		core: core,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { tpb.RegisterEventServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return tpb.RegisterEventServiceHandlerFromEndpoint(ctx, mux, address, options)
}
