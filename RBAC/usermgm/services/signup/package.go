package signup

import (
	"context"
	"path"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/serviceutil/auth/token"

	"brank.as/rbac/usermgm/storage"

	upb "brank.as/rbac/gunk/v1/user"
)

const defaultEnv = "development"

type UserStore interface {
	ConfirmEmail(context.Context, string) (*storage.User, error)
	ResetPassword(ctx context.Context, code, pw string) error
	ResetPasswordInit(context.Context, string) error
	CreateUser(context.Context, storage.User, storage.Credential) (*storage.User, string, error)
	GetConfirmationCode(ctx context.Context, code string) (*storage.User, string, error)
}

// Handler implementing the gunk/v1/signup APIs
type Handler struct {
	upb.UnimplementedSignupServer
	token.NoopAuthOverride
	env    string
	store  UserStore
	mailer Mailer
}

// Mailer for connecting to SMTP server
type Mailer interface {
	ConfirmEmail(email, code string) error
}

func New(store UserStore, mailer Mailer) *Handler {
	h := &Handler{
		env:    defaultEnv,
		store:  store,
		mailer: mailer,
	}
	return h
}

func (h *Handler) RegisterSvc(srv *grpc.Server) error { upb.RegisterSignupServer(srv, h); return nil }

// RegisterGateway grpcgw
func (h *Handler) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return upb.RegisterSignupHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (h *Handler) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"Signup":             {res: "ACCOUNT:user", act: "create", pub: true},
		"ForgotPassword":     {res: "ACCOUNT:user", act: "forgot", pub: true},
		"ResetPassword":      {res: "ACCOUNT:user", act: "reset", pub: true},
		"EmailConfirmation":  {res: "ACCOUNT:user", act: "confirm", pub: true},
		"ResendConfirmEmail": {res: "ACCOUNT:user", act: "resend", pub: true},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}

func (h *Handler) PublicEndpoint(method string) bool {
	if _, _, pub := h.Permission(context.Background(), path.Base(method)); pub {
		return true
	}
	return false
}
