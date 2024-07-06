package permissions

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) GrantService(ctx context.Context, g core.Grant) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.grantservice")

	org, err := s.store.GetOrgByID(ctx, g.RoleID)
	if err != nil {
		logging.WithError(err, log).Error("getting org")
		return "", status.Error(codes.NotFound, "org not found")
	}
	if !org.Active {
		logging.WithError(err, log).Error("org inactive")
		return "", status.Error(codes.NotFound, "org inactive")
	}

	// converts the svcaccount clientid into uuid of the creator of the svc account
	id := hydra.ClientID(ctx)
	if _, err := uuid.Parse(id); err != nil {
		sa, err := s.store.GetSvcAccountByID(ctx, id)
		if err != nil {
			logging.WithError(err, log).Error("getting service account")
			return "", status.Error(codes.NotFound, "svc account not found")
		}
		id = sa.CreateUserID
	}
	return s.AssignService(ctx, id, g)
}

// AssignService ...
func (s *Svc) AssignService(ctx context.Context, uid string, g core.Grant) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.assignservice")
	svc, err := s.store.GetAssignedService(ctx, storage.ServiceAssignment{
		ServiceID:   g.GrantID,
		Environment: g.Environment,
		OrgID:       g.RoleID,
	})
	if err == nil {
		log.WithField("service_id", g.GrantID).Debug("service assigned")
		return svc.GrantID, nil
	}

	sp, err := s.store.ListServicePermission(ctx, g.GrantID)
	if err != nil {
		logging.WithError(err, log).Error("getting service permissions")
		return "", status.Error(codes.NotFound, "service not found")
	}

	asn := storage.ServiceAssignment{
		ServiceID:    g.GrantID,
		Environment:  g.Environment,
		OrgID:        g.RoleID,
		AssignUserID: uid,
		Default:      g.Default,
	}

	grant, err := s.store.AssignService(ctx, asn)
	if err != nil {
		logging.WithError(err, log).Error("store service assignment")
		return "", status.Error(codes.Internal, "processing failed")
	}

	for _, p := range sp {
		if _, err := s.CreateOrgPermission(ctx, core.OrgPermission{
			OrgID:       g.RoleID,
			SvcID:       p.ServiceID,
			SvcPermID:   p.ID,
			GrantID:     grant,
			CreateUID:   uid,
			Name:        p.Name,
			Description: p.Description,
			Environment: g.Environment,
			Allow:       true,
			Action:      p.Action,
			Resource:    p.Resource,
			Groups:      []string{},
		}); err != nil {
			logging.WithError(err, log).Error("create org permission")
		}
	}

	return grant, nil
}
