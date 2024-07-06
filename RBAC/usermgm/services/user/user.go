package user

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"

	upb "brank.as/rbac/gunk/v1/user"
)

const defaultEnv = "development"

type Handler struct {
	upb.UnimplementedUserServiceServer
	acct AcctGetter
	usr  UserStore
}

type AcctGetter interface {
	GetUserByEmail(context.Context, string) (*storage.User, error)
	GetUserByID(context.Context, string) (*storage.User, error)
	GetUsersByOrg(context.Context, string) ([]storage.User, error)
	GetUsers(context.Context, storage.FilterList) ([]storage.User, error)
	GetOrgByID(context.Context, string) (*storage.Organization, error)
}

type UserStore interface {
	ChangePass(ctx context.Context, userid, oldPass, newPass string) (*core.MFAChallenge, error)
	ConfirmPass(context.Context, core.MFAChallenge) error
	UpdateUser(context.Context, core.User) (*core.MFAChallenge, error)
	DisableUser(context.Context, core.UserActivation) error
	EnableUser(context.Context, core.UserActivation) error
}

func New(get AcctGetter, usr UserStore) *Handler {
	return &Handler{acct: get, usr: usr}
}

// RegisterService with grpc server.
func (h *Handler) RegisterSvc(srv *grpc.Server) error {
	upb.RegisterUserServiceServer(srv, h)
	return nil
}

// RegisterGateway grpcgw.
func (h *Handler) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return upb.RegisterUserServiceHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (h *Handler) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"GetUser":        {res: "ACCOUNT:user", act: "view"},
		"ListUsers":      {res: "ACCOUNT:user", act: "view"},
		"DisableUser":    {res: "ACCOUNT:user", act: "delete"},
		"EnableUser":     {res: "ACCOUNT:user", act: "update"},
		"ChangePassword": {res: "ACCOUNT:user", act: "edit", pub: true},
		"ConfirmUpdate":  {res: "ACCOUNT:user", act: "edit", pub: true},
		"UpdateUser":     {res: "ACCOUNT:user", act: "edit", pub: true},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
