package permissions

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"

	"brank.as/rbac/usermgm/core"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

type PermissionStore interface {
	CreatePermission(ctx context.Context, p core.ServicePermission) (*core.ServicePermission, error)
	DeletePermission(context.Context, string) error
	ListPermission(context.Context, core.ListPermissionFilter) ([]core.OrgPermission, error)
}

type Validator interface {
	Validate(context.Context, core.Validation) (*core.Identity, error)
}

type Svc struct {
	ppb.UnimplementedPermissionServiceServer
	perm PermissionStore
	val  Validator
}

func New(prm PermissionStore, val Validator) *Svc {
	return &Svc{perm: prm, val: val}
}

func (s *Svc) RegisterSvc(svr *grpc.Server) error {
	ppb.RegisterPermissionServiceServer(svr, s)
	return nil
}

func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return ppb.RegisterPermissionServiceHandlerFromEndpoint(ctx, mux, address, options)
}

type resAct struct {
	res, act string
	pub      bool
}

func (s *Svc) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"CreatePermission": {res: "RBAC:permission", act: "create"},
		"ListPermission":   {res: "RBAC:permission", act: "view"},
		"DeletePermission": {res: "RBAC:permission", act: "delete"},
		"AssignPermission": {res: "RBAC:permission", act: "assign"},
		"RevokePermission": {res: "RBAC:permission", act: "assign"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
