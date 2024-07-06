package permissions

import (
	"context"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) PublicService(ctx context.Context, g core.Grant) error {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.publicservice")

	log.WithField("service", g.RoleID).Info("publishing")

	svc, err := s.store.GetService(ctx, g.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("getting service")
		return status.Error(codes.NotFound, "service not found")
	}
	log = log.WithField("service", svc.ID)

	pub, err := s.store.GetPublicService(ctx, svc.ID)
	if err != nil && err != storage.NotFound {
		logging.WithError(err, log).Error("getting public service")
		return status.Error(codes.NotFound, "service not found")
	}
	for _, p := range pub {
		if p.Environment == g.Environment {
			log.WithFields(logrus.Fields{
				"service": g.GrantID,
				"env":     g.Environment,
			}).Debug("already public")
			return nil
		}
	}

	if _, err := s.store.UpsertPublicService(ctx, storage.DefaultService{
		ServiceID:       g.GrantID,
		Environment:     g.Environment,
		PublishedUserID: g.RoleID,
	}); err != nil {
		logging.WithError(err, log).Error("store publish record")
		return status.Error(codes.Internal, "failed to publish")
	}

	orgs, err := s.store.GetOrgs(ctx)
	if err != nil {
		logging.WithError(err, log).Error("store publish record")
		return status.Error(codes.Internal, "failed to publish")
	}

	asn, err := s.store.ListServiceAssignID(ctx, g.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("store fetch service assign record")
		return status.Error(codes.Internal, "failed to publish")
	}
	m := make(map[string]bool)
	for _, a := range asn {
		m[a.OrgID] = true
	}

	gr := g
	uid := g.RoleID
	gr.Default = true
	for _, o := range orgs {
		if m[o.ID] {
			continue
		}
		gr.RoleID = o.ID
		_, err := s.AssignService(ctx, uid, gr)
		if err != nil {
			logging.WithError(err, log).WithField("org", o.ID).Error("granting service")
		}
	}
	return nil
}
