package userauth

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"

	apb "brank.as/rbac/gunk/v1/authenticate"
	ppb "brank.as/rbac/gunk/v1/user"
)

type UserStore interface {
	GetUser(ctx context.Context, username, password string) (*storage.User, error)
	GetUserByID(context.Context, string) (*storage.User, error)
}

type AuthStore interface {
	AuthUser(context.Context, core.AuthCredential) (*core.Identity, error)
	UserSession(context.Context, string) (*core.Identity, error)
	NewMFA(context.Context, core.AuthCredential) (*core.Identity, error)
}

type Svc struct {
	ppb.UnsafeUserAuthServiceServer
	apb.UnsafeSessionServiceServer
	perm UserStore
	auth AuthStore
}

func New(prm UserStore, a AuthStore) *Svc {
	return &Svc{
		perm: prm,
		auth: a,
	}
}

func (s *Svc) RegisterSvc(svr *grpc.Server) error {
	ppb.RegisterUserAuthServiceServer(svr, s)
	apb.RegisterSessionServiceServer(svr, s)
	return nil
}

func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, addr string, opt []grpc.DialOption) error {
	if err := apb.RegisterSessionServiceHandlerFromEndpoint(ctx, mux, addr, opt); err != nil {
		return err
	}
	return ppb.RegisterUserAuthServiceHandlerFromEndpoint(ctx, mux, addr, opt)
}
