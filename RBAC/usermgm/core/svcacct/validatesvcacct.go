package svcacct

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/storage"
)

// ValidateSvcAccount as an api key token.
func (s *Svc) ValidateSvcAccount(ctx context.Context, key string) (*storage.SvcAccount, error) {
	spl := strings.Split(key, ".")
	if len(spl) == 2 && len(spl[0]) >= APIKeyPrefix {
		return s.ValidateKey(ctx, key)
	}
	return s.ValidateOAuth(ctx, key)
}

// ValidateOAuth as an oauth2 client.
func (s *Svc) ValidateOAuth(ctx context.Context, id string) (*storage.SvcAccount, error) {
	log := logging.FromContext(ctx).WithField("method", "core.svcacct.validateoauth")
	sa, err := s.store.GetSvcAccountByID(ctx, id)
	if err != nil {
		logging.WithError(err, log).Error("storage validate")
		return nil, status.Error(codes.NotFound, "validation failed")
	}
	org, err := s.store.GetOrgByID(ctx, sa.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("storage org")
		return nil, status.Error(codes.NotFound, "validation failed")
	}
	if !org.Active {
		log.WithField("org_id", org.ID).Error("inactive org")
		return nil, status.Error(codes.NotFound, "validation failed")
	}
	return sa, nil
}

// ValidateKey as an api key token.
func (s *Svc) ValidateKey(ctx context.Context, key string) (*storage.SvcAccount, error) {
	log := logging.FromContext(ctx).WithField("method", "core.svcacct.validatekey")
	spl := strings.Split(key, ".")
	if len(spl) != 2 {
		return nil, status.Error(codes.InvalidArgument, "invalid key")
	}
	id, k := spl[0], strings.Join(spl, "")
	sa, err := s.store.ValidateSvcAccount(ctx, id, k)
	if err != nil {
		logging.WithError(err, log).Error("storage validate")
		return nil, status.Error(codes.NotFound, "validation failed")
	}
	org, err := s.store.GetOrgByID(ctx, sa.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("storage org")
		return nil, status.Error(codes.NotFound, "validation failed")
	}
	if !org.Active {
		log.WithField("org_id", org.ID).Error("inactive org")
		return nil, status.Error(codes.NotFound, "validation failed")
	}
	return sa, nil
}
