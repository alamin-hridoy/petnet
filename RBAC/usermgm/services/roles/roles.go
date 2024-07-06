package role

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

type RoleStore interface {
	CreateRole(context.Context, core.Role) (string, error)
	RoleGrant(context.Context, core.Grant) (*core.Role, error)
	RoleRevoke(context.Context, core.Grant) (*core.Role, error)
	AssignRole(context.Context, core.Grant) (*core.Role, error)
	UnassignRole(context.Context, core.Grant) (*core.Role, error)
	ListRole(context.Context, core.ListRoleFilter) ([]core.Role, error)
	ListUserRoles(context.Context, core.ListUserRolesRequest) (map[string]*ppb.UserRoles, error)
	UpdateRole(context.Context, core.Role) (core.Role, error)
	DeleteRole(ctx context.Context, id string) error
}

type Svc struct {
	ppb.UnimplementedRoleServiceServer
	perm RoleStore
}

func New(prm RoleStore) *Svc {
	return &Svc{
		perm: prm,
	}
}

func (s *Svc) RegisterSvc(svr *grpc.Server) error {
	ppb.RegisterRoleServiceServer(svr, s)
	return nil
}

func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, addr string, opt []grpc.DialOption) error {
	return ppb.RegisterRoleServiceHandlerFromEndpoint(ctx, mux, addr, opt)
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"CreateRole":           {res: "RBAC:role", act: "create"},
		"UpdateRole":           {res: "RBAC:role", act: "create"},
		"ListRole":             {res: "RBAC:role", act: "view"},
		"ListUserRoles":        {res: "RBAC:role", act: "view", pub: true},
		"DeleteRole":           {res: "RBAC:role", act: "delete"},
		"AddUser":              {res: "RBAC:role", act: "assign"},
		"RemoveUser":           {res: "RBAC:role", act: "assign"},
		"AssignRolePermission": {res: "RBAC:role", act: "assign"},
		"RevokeRolePermission": {res: "RBAC:role", act: "assign"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
