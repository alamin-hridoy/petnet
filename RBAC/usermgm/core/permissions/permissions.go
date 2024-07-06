package permissions

import (
	"context"

	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage/postgres"
)

type ketoClient interface {
	CreatePermission(ctx context.Context, p keto.Permission) (string, error)
	GetPermission(ctx context.Context, id string) (keto.Permission, error)
	UpdatePermission(ctx context.Context, p keto.Permission) error
	DeletePermission(ctx context.Context, id string) error

	CreateRole(ctx context.Context, ro keto.Role) (string, error)
	UpdateRole(ctx context.Context, ro keto.Role) (string, error)
	ListRoles(ctx context.Context, uid string) ([]string, error)
	GetRole(ctx context.Context, id string) (keto.Role, error)
	DeleteRole(ctx context.Context, id string) error

	GetRolePermissions(ctx context.Context, role string) ([]string, error)
}

type Svc struct {
	store *postgres.Storage
	keto  ketoClient
}

func New(store *postgres.Storage, k ketoClient) *Svc {
	return &Svc{
		store: store,
		keto:  k,
	}
}
