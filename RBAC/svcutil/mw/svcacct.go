package mw

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

type SvcAcct struct {
	cl sapb.ValidationServiceClient
}

func NewServiceAccount(cl sapb.ValidationServiceClient) *SvcAcct {
	return &SvcAcct{cl: cl}
}

// Metadata for service account validation
func (s *SvcAcct) Metadata(ctx context.Context) (context.Context, error) {
	log := logging.FromContext(ctx).WithField("method", "saloader.metadata")
	if GetOrg(ctx) != "" {
		return ctx, nil
	}
	v, err := s.cl.ValidateAccount(ctx, &sapb.ValidateAccountRequest{
		ClientID: hydra.ClientID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Trace("validate SA")
		return ctx, nil
	}

	return metautils.ExtractIncoming(ctx).
		Add(OrgIDKey, v.GetOrgID()).
		Add(EnvKey, v.GetEnvironment()).
		Add(IDNameKey, v.GetClientName()).ToIncoming(ctx), nil
}
