package main

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/storage/postgres"
)

type OrgLoader struct {
	st *postgres.Storage
}

func NewOrg(st *postgres.Storage) *OrgLoader {
	return &OrgLoader{st: st}
}

// Metadata for service account validation
func (s *OrgLoader) Metadata(ctx context.Context) (context.Context, error) {
	log := logging.FromContext(ctx).WithField("method", "orgloader.metadata")
	// log.WithField("md", metautils.ExtractIncoming(ctx)).Info("metadata")

	if mw.GetOrg(ctx) != "" {
		return ctx, nil
	}
	v, err := s.st.GetUserByID(ctx, hydra.ClientID(ctx))
	if err != nil {
		logging.WithError(err, log).Trace("get user by id")
		return ctx, nil
	}

	return metautils.ExtractIncoming(ctx).
		Add(mw.OrgIDKey, v.OrgID).
		Add(mw.IDNameKey, v.Username).ToIncoming(ctx), nil
}
