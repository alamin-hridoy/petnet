package invite

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"

	ipb "brank.as/rbac/gunk/v1/invite"
)

const defaultEnv = "development"

type Mailer interface {
	InviteUser(email, invCode string, inv email.Invite) error
	Approved(email string, ap email.Approved) error
	// InviteExpired(email string) error
}

// FakeMailer for local dev and testing
type FakeMailer struct {
	WantErr bool
	Code    string
}

func (f FakeMailer) Invite(email, code string) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	if f.Code != code {
		return fmt.Errorf("different code")
	}
	return nil
}

type Svc struct {
	ipb.UnimplementedInviteServiceServer
	env         string
	servicename string
	inv         UserInviter
	store       InviteStore
	plt         OrgStore
	mailer      Mailer
}

type UserInviter interface {
	InviteUser(context.Context, core.Invite) (*core.Invite, error)
}

type InviteStore interface {
	CreateUser(context.Context, storage.User, storage.Credential) (*storage.User, error)
	GetUsersByOrg(context.Context, string) ([]storage.User, error)
	GetUserByID(context.Context, string) (*storage.User, error)
	GetUserByInvite(ctx context.Context, code string) (*storage.User, error)
	UpdateOrgByID(context.Context, storage.Organization) (*storage.Organization, error)
	UpdateUserByID(context.Context, storage.User) (*storage.User, error)
	ReviveOrgByID(context.Context, string) error
	ReviveUserByID(context.Context, string) error
}

type OrgStore interface {
	ActivateOrg(context.Context, string) error
	GetOrgByID(context.Context, string) (*storage.Organization, error)
	CreateOrg(context.Context, storage.Organization) (string, error)
}

// Option for the signup Handler constructor
type Option func(h *Svc)

// WithEnv used by the logger
func WithEnv(env string) Option {
	return func(h *Svc) { h.env = env }
}

// TODO: refactor core logic out of service endpoints.
func New(store InviteStore, plt OrgStore, inv UserInviter, mailer Mailer, opts ...Option) *Svc {
	h := &Svc{
		env:    defaultEnv,
		store:  store,
		plt:    plt,
		mailer: mailer,
		inv:    inv,
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// Register service.
func (h *Svc) RegisterSvc(srv *grpc.Server) error {
	ipb.RegisterInviteServiceServer(srv, h)
	return nil
}

// RegisterGateway grpcgw
func (h *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ipb.RegisterInviteServiceHandlerFromEndpoint(ctx, mux, address, options)
}

func validateInviteStatus(statusMap map[string]struct{}, inviteStatus string) bool {
	_, ok := statusMap[inviteStatus]
	return ok
}

type resAct struct {
	res, act string
	pub      bool
}

func (h *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"InviteUser":     {res: "ACCOUNT:user", act: "invite"},
		"Resend":         {res: "ACCOUNT:user", act: "invite"},
		"Approve":        {res: "ACCOUNT:user", act: "create"},
		"ListInvite":     {res: "ACCOUNT:user", act: "view"},
		"RetrieveInvite": {res: "ACCOUNT:user", act: "view"},
		"Revoke":         {res: "ACCOUNT:user", act: "delete"},
		"CancelInvite":   {res: "ACCOUNT:user", act: "delete"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
