package org

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/core"
)

// ActivaceOrg initialize a new organization structure.
func (s *Svc) ActivateOrg(ctx context.Context, user string, org string) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "org.activateorg")

	if err := s.org.ActivateOrg(ctx, org); err != nil {
		logging.WithError(err, log).Error("activating org")
		return "", status.Error(codes.Internal, "processing failed")
	}

	p, err := s.pm.ListRole(ctx, core.ListRoleFilter{
		ID:    []string{},
		OrgID: org,
	})
	if err == nil && len(p) != 0 {
		log.WithField("org", org).Debug("previously activated")
		return org, nil
	}

	// Bootstrap default permissions to org
	ps, err := s.org.ListServicePublic(ctx)
	if err != nil {
		logging.WithError(err, log).Error("listing default services")
		return "", status.Error(codes.Internal, "processing failed")
	}

	for _, svc := range ps {
		if _, err := s.pm.AssignService(ctx, user, core.Grant{
			RoleID:      org,
			GrantID:     svc.ServiceID,
			Environment: svc.Environment,
			Default:     true,
		}); err != nil {
			logging.WithError(err, log).Error("granting", svc.ServiceID)
		}
	}

	// Bootstrap default permissions to org
	pm, err := s.org.ListOrgPermission(ctx, org)
	if err != nil {
		logging.WithError(err, log).Error("listing default services")
		return "", status.Error(codes.Internal, "processing failed")
	}

	// Create Owner Role with owner user
	r, err := s.pm.CreateRole(ctx, core.Role{
		OrgID:     org,
		Name:      "Owner",
		Desc:      "Orgaization Owner",
		CreateUID: hydra.ClientID(ctx),
		Members:   []string{user},
	})
	if err != nil {
		logging.WithError(err, log).Error("creating role")
		return "", status.Error(codes.Internal, "processing failed")
	}

	// Assign Permissions to Role
	for _, p := range pm {
		if _, err := s.pm.RoleGrant(ctx, core.Grant{
			RoleID:      r,
			GrantID:     p.ID,
			Environment: p.Environment,
			Default:     true,
		}); err != nil {
			logging.WithError(err, log).Error("granting permission")
			return "", status.Error(codes.Internal, "processing failed")
		}
	}
	return org, nil
}
