package org

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) UpdateOrg(ctx context.Context, org storage.Organization) (*storage.Organization, error) {
	log := logging.FromContext(ctx).WithField("method", "core.org.updateorg")

	o, err := s.org.UpdateOrgByID(ctx, org)
	if err != nil {
		logging.WithError(err, log).Error("update db")
		if err == storage.NotFound {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}

	return o, nil
}
