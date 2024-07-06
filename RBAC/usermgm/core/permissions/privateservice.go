package permissions

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) PrivateService(ctx context.Context, g core.Grant) error {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.publicservice")

	svc, err := s.store.GetService(ctx, g.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("getting service")
		return status.Error(codes.NotFound, "service not found")
	}
	log = log.WithField("service", svc.ID)

	pub, err := s.store.GetPublicService(ctx, svc.ID)
	if err != nil && err != storage.NotFound {
		if err == storage.NotFound {
			log.WithFields(logrus.Fields{
				"service": g.GrantID,
				"env":     g.Environment,
			}).Debug("already private")
		}
		logging.WithError(err, log).Error("getting public service")
		return status.Error(codes.NotFound, "service not found")
	}
	pid := ""
	for _, p := range pub {
		if p.Environment == g.Environment {
			pid = p.GrantID
			break
		}
	}

	if err := s.store.RetractService(ctx, storage.DefaultService{
		GrantID:         pid,
		ServiceID:       g.GrantID,
		Environment:     g.Environment,
		RetractedUserID: hydra.ClientID(ctx),
	}); err != nil {
		logging.WithError(err, log).Error("store retract record")
		return status.Error(codes.Internal, "failed to retract")
	}

	asn, err := s.store.ListServiceAssignID(ctx, g.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("store publish record")
		return status.Error(codes.Internal, "failed to retract")
	}

	gr := g
	gr.Default = true
	for _, a := range asn {
		if !a.Default {
			continue
		}
		gr.RoleID = a.OrgID
		if err := s.RevokeService(ctx, gr); err != nil {
			logging.WithError(err, log).WithField("org", a.OrgID).Error("revoking service")
		}
	}
	return nil
}
