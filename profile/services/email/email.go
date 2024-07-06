package fees

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	epb "brank.as/petnet/gunk/dsa/v1/email"
	eml "brank.as/petnet/profile/integrations/email"
)

type Svc struct {
	epb.UnimplementedEmailServiceServer
	es EmailStore
}

type EmailStore interface {
	SendOnboardingReminder(ctx context.Context, email string, orgID string, userID string) error
	SendDsaServiceRequestNotification(ctx context.Context, req eml.DsaServiceRequestNotificationForm) error
}

func New(es EmailStore) *Svc {
	return &Svc{
		es: es,
	}
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { epb.RegisterEmailServiceServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return epb.RegisterEmailServiceHandlerFromEndpoint(ctx, mux, address, options)
}
