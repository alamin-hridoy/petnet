package mw

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
)

const (
	orgID      = "org-id"
	orgName    = "org-name"
	orgCountry = "org-country"
)

type Identity struct {
	OrgID      string
	OrgName    string
	OrgCountry string
}

type IdentityGetter interface {
	GetIdentity(ctx context.Context, id string) (Identity, error)
}

type Identify struct {
	id IdentityGetter
}

// NewIdentity creates an identity metadata middleware.
func NewIdentity(getter IdentityGetter) *Identify {
	return &Identify{id: getter}
}

// Metadata loads the org identity into context.
func (ity Identify) Metadata(ctx context.Context) (context.Context, error) {
	log := logging.FromContext(ctx).WithField("method", "identity.metadata")
	cl := hydra.ClientID(ctx)
	if cl == "" {
		return nil, fmt.Errorf("identity: missing client ID in metadata")
	}
	idn, err := ity.id.GetIdentity(ctx, cl)
	if err != nil {
		logging.WithError(err, log).Error("fetch id failed")
		return nil, fmt.Errorf("identity: client ID does not match a org identity")
	}
	md := metautils.ExtractIncoming(ctx).
		Set(orgID, idn.OrgID).
		Set(orgName, idn.OrgName).
		Set(orgCountry, idn.OrgCountry)
	return md.ToIncoming(ctx), nil
}

// GetOrg returns the org ID
func GetOrg(ctx context.Context) string { return metautils.ExtractIncoming(ctx).Get(orgID) }

// GetOrgName returns the org name
func GetOrgName(ctx context.Context) string { return metautils.ExtractIncoming(ctx).Get(orgName) }

// GetCountry returns the country ID
func GetCountry(ctx context.Context) string { return metautils.ExtractIncoming(ctx).Get(orgCountry) }
func GetEnv(ctx context.Context) string     { return metautils.ExtractIncoming(ctx).Get(EnvKey) }
