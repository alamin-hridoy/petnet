package md

import (
	"context"
	"fmt"

	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	mta "github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

const (
	OrgIDKey    = "org-id"
	UserIDKey   = "user-id"
	UsernameKey = "user-name"
)

type Id struct {
	UserID   string
	Username string
	OrgID    string
}

type IdentityGetter interface {
	GetIdentity(ctx context.Context, client string) (Id, error)
}

type Identity struct {
	id IdentityGetter
}

// NewIdentity creates an identity metadata middleware.
func NewIdentity(getter IdentityGetter) *Identity {
	return &Identity{id: getter}
}

// Metadata loads the platform identity into context.
func (ity Identity) Metadata(ctx context.Context) (context.Context, error) {
	log := logging.FromContext(ctx)
	cl := hydra.ClientID(ctx)
	if cl == "" {
		return nil, fmt.Errorf("identity: missing client ID in metadata")
	}
	idn, err := ity.id.GetIdentity(ctx, cl)
	if err != nil {
		logging.WithError(err, log).Error("fetch id failed")
		return nil, fmt.Errorf("identity: client ID does not match a platform identity")
	}
	md := mta.ExtractIncoming(ctx).
		Set(UserIDKey, idn.UserID).
		Set(UsernameKey, idn.Username).
		Set(OrgIDKey, idn.OrgID)
	return md.ToIncoming(ctx), nil
}

// GetUser returns the User ID
func GetUser(ctx context.Context) string { return mta.ExtractIncoming(ctx).Get(UserIDKey) }

// GetUserName returns the User name
func GetUserName(ctx context.Context) string { return mta.ExtractIncoming(ctx).Get(UsernameKey) }

// GetOrg returns the org ID
func GetOrgID(ctx context.Context) string { return mta.ExtractIncoming(ctx).Get(OrgIDKey) }
