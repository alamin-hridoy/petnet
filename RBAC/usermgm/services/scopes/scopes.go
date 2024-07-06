package scopes

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"

	cpb "brank.as/rbac/gunk/v1/consent"
)

type Svc struct {
	cpb.UnsafeScopeServiceServer
	cpb.UnsafeGrantServiceServer
	sc ScopeStore
	gr GrantStore
}

func New(sc ScopeStore, gr GrantStore) *Svc {
	return &Svc{sc: sc, gr: gr}
}

type ScopeStore interface {
	UpsertScope(context.Context, core.Scope) (*core.Scope, error)
	UpdateGroup(context.Context, core.ScopeGroup) (*core.ScopeGroup, error)
	GetScopes(context.Context, []string) (map[string]core.ScopeGroup, error)
}

type GrantStore interface {
	GetGrant(context.Context, core.ConsentGrant) (*core.OfferGrant, error)
	RecordGrant(context.Context, core.ConsentGrant) (*core.ConsentGrant, error)
}

func (s *Svc) RegisterSvc(svr *grpc.Server) error {
	cpb.RegisterScopeServiceServer(svr, s)
	cpb.RegisterGrantServiceServer(svr, s)
	return nil
}

func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, addr string, opt []grpc.DialOption) error {
	if err := cpb.RegisterGrantServiceHandlerFromEndpoint(ctx, mux, addr, opt); err != nil {
		return err
	}
	return cpb.RegisterScopeServiceHandlerFromEndpoint(ctx, mux, addr, opt)
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"UpsertScope": {res: "RBAC:scope", act: "create"},
		"UpdateGroup": {res: "RBAC:scope", act: "create"},
		"GetScope":    {res: "RBAC:scope", act: "view"},
		"ServeGrant":  {res: "RBAC:consent", act: "view", pub: true},
		"Grant":       {res: "RBAC:consent", act: "create", pub: true},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
