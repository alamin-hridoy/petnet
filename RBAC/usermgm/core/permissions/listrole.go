package permissions

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/rbac/gunk/v1/permissions"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"
	st "brank.as/rbac/usermgm/storage"
)

func (s *Svc) ListUserRoles(ctx context.Context, f core.ListUserRolesRequest) (map[string]*ppb.UserRoles, error) {
	listRole := make(map[string]*ppb.UserRoles)
	for _, uuid := range f.UserID {
		userRoles, err := s.keto.ListRoles(ctx, uuid)
		if err == nil {
			listRole[uuid] = &ppb.UserRoles{
				UserRoles: userRoles,
			}
		} else {
			listRole[uuid] = &ppb.UserRoles{}
		}
	}
	return listRole, nil
}

func (s *Svc) ListRole(ctx context.Context, f core.ListRoleFilter) ([]core.Role, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permissions.listrole")
	var storeRoles []st.Role
	var err error
	if f.UserID != "" {
		userRoles, err := s.keto.ListRoles(ctx, f.UserID)
		if err != nil {
			logging.WithError(err, log).Error("Unable to connect Keto api")
			return nil, status.Error(codes.Internal, "processing failed")
		}
		if len(userRoles) > 0 {
			storeRoles, err = s.store.ListRole(ctx, core.ListRoleFilter{
				ID: userRoles,
			})
			if err != nil {
				logging.WithError(err, log).Error("read storage")
				return nil, status.Error(codes.Internal, "processing failed")
			}
		}
	} else {
		storeRoles, err = s.store.ListRole(ctx, f)
		if err != nil {
			logging.WithError(err, log).Error("read storage")
			return nil, status.Error(codes.Internal, "processing failed")
		}
	}

	var coreRoles []core.Role
	SkipCount := 0
	for _, r := range storeRoles {
		// Filter IDs
		if len(f.ID) > 0 {
			found := false
			for _, id := range f.ID {
				if id == r.ID {
					found = true
					break
				}
			}
			if !found {
				SkipCount = SkipCount + 1
				continue
			}
		}
		if r.Delete.Valid {
			SkipCount = SkipCount + 1
			continue
		}

		ketoRole, err := s.keto.GetRole(ctx, r.ID)
		if err != nil {
			logging.WithError(err, log).Error("read keto role")
			return nil, status.Error(codes.Internal, "processing failed")
		}

		ketoRolePerms, err := s.keto.GetRolePermissions(ctx, r.ID)
		if err != nil {
			logging.WithError(err, log).Error("read keto role permission")
			return nil, status.Error(codes.Internal, "processing failed")
		}
		coreRoles = append(coreRoles, core.Role{
			ID:          r.ID,
			OrgID:       r.OrgID,
			Name:        r.Name,
			Desc:        r.Description,
			CreateUID:   r.CreateUID,
			DeleteUID:   r.DeleteUID.String,
			Permissions: ketoRolePerms,
			Members:     ketoRole.Members,
			Created:     r.Created,
			UpdatedUID:  r.UpdatedUID,
			Updated:     r.Updated,
			Count:       r.Count - SkipCount,
		})
	}

	return coreRoles, nil
}
