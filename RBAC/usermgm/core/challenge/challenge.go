package challenge

import (
	"context"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

type Validator interface {
	ValidateRequest(ctx context.Context, v core.Validation) (bool, error)
}

type IDLookup interface {
	GetUserByID(ctx context.Context, id string) (*storage.User, error)
	GetSvcAccountByID(ctx context.Context, clientID string) (*storage.SvcAccount, error)
	GetOrgByID(context.Context, string) (*storage.Organization, error)
}

type Svc struct {
	id  IDLookup
	val Validator
}

// New challenger
func New(st IDLookup, v Validator) *Svc {
	return &Svc{id: st, val: v}
}
