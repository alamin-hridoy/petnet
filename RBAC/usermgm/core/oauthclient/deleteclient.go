package oauthclient

import (
	"context"
	"database/sql"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

// DeleteClient ...
func (s *Svc) DeleteClient(ctx context.Context, cl core.AuthCodeClient) (*core.AuthCodeClient, error) {
	log := logging.FromContext(ctx).WithField("method", "core.oauthclient.deleteclient")

	c, err := s.st.GetOauthClientByID(ctx, cl.ClientID)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "client not found")
		}
		return nil, err
	}

	if cl.OrgID != c.OrgID {
		log.WithField("request_org", cl.OrgID).
			WithField("target_org", c.OrgID).
			Error("org mismatch")
		return nil, status.Error(codes.NotFound, "client not found")
	}

	if err := s.hy.DeleteClient(ctx, cl.ClientID); err != nil {
		logging.WithError(err, log).Error("hydra delete")
		return nil, err
	}

	c.DeleteUserID = cl.DeletedBy
	st, err := s.st.DeleteOauthClient(ctx, *c)
	if err != nil {
		logging.WithError(err, log).Error("storage delete")
		return nil, err
	}
	c.Deleted = sql.NullTime{
		Time:  *st,
		Valid: true,
	}
	return &cl, nil
}
