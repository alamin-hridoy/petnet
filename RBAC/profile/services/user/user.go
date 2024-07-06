package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/profile/storage"

	ipupb "brank.as/rbac/gunk/v1/user"
	upb "brank.as/rbac/profile/gunk/v1/useraccount"
)

const defaultEnv = "development"

type Handler struct {
	upb.UnimplementedUserServiceServer
	acct AcctGetter
	cl   ipupb.UserServiceClient
}

type AcctGetter interface {
	GetUserByID(context.Context, string) (*storage.User, error)
	GetUsersByOrg(context.Context, string) ([]storage.User, error)
}

func New(get AcctGetter, cl ipupb.UserServiceClient) *Handler {
	return &Handler{acct: get, cl: cl}
}

// RegisterService with grpc server.
func (h *Handler) RegisterSvc(srv *grpc.Server) error {
	upb.RegisterUserServiceServer(srv, h)
	return nil
}

// RegisterGateway grpcgw
func (h *Handler) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return upb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, address, options)
}
