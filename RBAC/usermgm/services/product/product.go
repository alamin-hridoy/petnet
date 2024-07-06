package product

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core/permissions"
	"brank.as/rbac/usermgm/storage/postgres"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

type Svc struct {
	ppb.UnimplementedProductServiceServer
	pm  *permissions.Svc
	st  *postgres.Storage
	env []interface{}
}

func New(pm *permissions.Svc, st *postgres.Storage, allowedEnv []string) *Svc {
	e := make([]interface{}, len(allowedEnv)+1)
	for i, env := range allowedEnv {
		e[i] = env
	}
	return &Svc{pm: pm, st: st, env: e}
}

func (s *Svc) RegisterSvc(svr *grpc.Server) error {
	ppb.RegisterProductServiceServer(svr, s)
	return nil
}

func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"GrantService":           {res: "RBAC:service", act: "create"},
		"ListServiceAssignments": {res: "RBAC:service", act: "view"},
		"RevokeService":          {res: "RBAC:service", act: "delete"},
		"PublicService":          {res: "RBAC:service", act: "publish"},
		"ListServices":           {res: "RBAC:service", act: "view"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
