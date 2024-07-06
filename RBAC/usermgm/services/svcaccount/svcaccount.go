package svcaccount

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/storage"

	ppb "brank.as/rbac/gunk/v1/permissions"
	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

const defaultEnv = "development"

type Svc struct {
	sapb.UnimplementedSvcAccountServiceServer
	sapb.UnimplementedValidationServiceServer
	env     []interface{}
	plt     OrgLookup
	store   AcctStore
	get     AcctGetter
	svcName string
	val     ppb.ValidationServiceClient
}

type OrgLookup interface {
	UserOrg(ctx context.Context) (*storage.Organization, error)
}

type AcctStore interface {
	CreateSvcAccount(ctx context.Context, sa storage.SvcAccount) (id, secret string, err error)
	DisableSvcAccount(ctx context.Context, sa storage.SvcAccount) (*time.Time, error)
	ValidateSvcAccount(ctx context.Context, key string) (*storage.SvcAccount, error)
}

type AcctGetter interface {
	GetSvcAccountByOrgID(context.Context, string) ([]storage.SvcAccount, error)
	GetSvcAccountByID(context.Context, string) (*storage.SvcAccount, error)
	GetOrgByID(context.Context, string) (*storage.Organization, error)
}

func New(store AcctStore, get AcctGetter, plt OrgLookup, envs []string, val ppb.ValidationServiceClient) *Svc {
	e := make([]interface{}, len(envs))
	for i, ev := range envs {
		e[i] = ev
	}
	return &Svc{
		env:   e,
		plt:   plt,
		get:   get,
		store: store,
		val:   val,
	}
}

// RegisterService with grpc server.
func (h *Svc) RegisterSvc(srv *grpc.Server) error {
	sapb.RegisterSvcAccountServiceServer(srv, h)
	return nil
}

// RegisterGateway grpcgw
func (h *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return sapb.RegisterSvcAccountServiceHandlerFromEndpoint(ctx, mux, address, options)
}

func (h *Svc) RegisterInternal(srv *grpc.Server) error {
	sapb.RegisterValidationServiceServer(srv, h)
	return nil
}

// RegisterGateway grpcgw
func (h *Svc) RegisterGatewayInternal(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return sapb.RegisterValidationServiceHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (h *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"CreateAccount":  {res: "ACCOUNT:service", act: "create"},
		"ListAccounts":   {res: "ACCOUNT:service", act: "view"},
		"DisableAccount": {res: "ACCOUNT:service", act: "create"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
