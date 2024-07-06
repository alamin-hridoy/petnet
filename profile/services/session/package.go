package session

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	spb "brank.as/petnet/gunk/v1/session"
	c "brank.as/petnet/profile/core/session"
)

type Svc struct {
	spb.UnimplementedSessionServiceServer
	core SessionStore
}

type SessionStore interface {
	UpsertSession(ctx context.Context, req *c.UpsertSessionReq) (string, error)
	GetSession(ctx context.Context, req *c.GetSessionReq) (*c.GetSessionRes, error)
}

func New(c SessionStore) *Svc {
	return &Svc{
		core: c,
	}
}

// RegisterGateway grpcgw
func RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return spb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, address, options)
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { spb.RegisterSessionServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return spb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, address, options)
}
