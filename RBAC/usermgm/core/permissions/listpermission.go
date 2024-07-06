package permissions

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) ListPermission(ctx context.Context, f core.ListPermissionFilter) ([]core.OrgPermission, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.listpermission")

	sp, err := s.store.ListOrgPermission(ctx, f.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("read storage")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	grants, err := s.store.ListServiceAssignOrg(ctx, f.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("read grant storage")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	gr := make(map[string]storage.ServiceAssignment)
	for _, g := range grants {
		gr[g.GrantID] = g
	}
	snm, err := s.store.ListService(ctx)
	if err != nil {
		logging.WithError(err, log).Error("read service storage")
		return nil, status.Error(codes.Internal, "processing failed")
	}
	sn := make(map[string]storage.Service)
	for _, s := range snm {
		sn[s.ID] = s
	}

	var corePerms []core.OrgPermission
	for _, p := range sp {
		// Filter IDs
		if len(f.ID) > 0 {
			found := false
			for _, id := range f.ID {
				if id == p.ID {
					found = true
					break
				}
			}

			if !found {
				continue
			}
		}
		if f.Environment != "" && p.Environment != f.Environment {
			continue
		}

		pm, err := s.keto.GetPermission(ctx, p.ID)
		if err != nil {
			logging.WithError(err, log).Error("read keto")
			return nil, err // status.Error(codes.Internal, "processing failed")
		}
		g := gr[p.GrantID]
		s := sn[p.ServiceID]
		for _, act := range pm.Actions {
			for _, res := range pm.Resources {
				r := strings.SplitN(res, ":", 3)
				if len(r) < 3 {
					continue
				}
				res = r[2]
				corePerms = append(corePerms, core.OrgPermission{
					ID:          p.ID,
					OrgID:       g.OrgID,
					SvcID:       g.ServiceID,
					SvcName:     s.Name,
					SvcPermID:   p.PermissionID,
					GrantID:     p.GrantID,
					CreateUID:   g.AssignUserID,
					DeleteUID:   g.RevokeUserID.String,
					Name:        p.Name,
					Description: p.Description,
					Environment: pm.Environment,
					Allow:       pm.Allow,
					Action:      act,
					Resource:    res,
					Groups:      pm.Groups,
				})
			}
		}
	}

	return corePerms, nil
}
