package org

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/storage"
)

// CreateOrg initialize a new organization structure.
func (s *Svc) CreateOrg(ctx context.Context, org storage.Organization) (string, error) {
	log := logging.FromContext(ctx).WithField("method", "org.createorg")
	id, err := s.org.CreateOrg(ctx, org)
	if err != nil {
		logging.WithError(err, log).Error("creating org")
		return "", status.Error(codes.Internal, "processing failed")
	}

	return id, nil
}
