package profile

import (
	"context"
	"time"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GetOrgProfile ..
func (s *Svc) GetOrgProfile(ctx context.Context, id string) (*storage.OrgProfile, error) {
	log := logging.FromContext(ctx)

	org, err := s.st.GetOrgProfile(ctx, id)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "org not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "profile failed")
	}
	return org, nil
}

// GetOrgProfiles ..
func (s *Svc) GetOrgProfiles(ctx context.Context) ([]storage.OrgProfile, error) {
	return s.st.GetOrgProfiles(ctx, storage.FilterList{})
}

// GetUserProfile ..
func (s *Svc) GetUserProfile(ctx context.Context, uid string) (*storage.UserProfile, error) {
	log := logging.FromContext(ctx)

	org, err := s.st.GetUserProfile(ctx, uid)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user profile not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "profile failed")
	}
	return org, nil
}

// GetUserProfiles ..
func (s *Svc) GetUserProfiles(ctx context.Context, oid string) ([]storage.UserProfile, error) {
	return s.st.GetUserProfiles(ctx, oid)
}

// GetUserProfileByEmail ..
func (s *Svc) GetUserProfileByEmail(ctx context.Context, email string) (*storage.UserProfile, error) {
	log := logging.FromContext(ctx)

	org, err := s.st.GetUserProfileByEmail(ctx, email)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "user profile not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "profile failed")
	}
	return org, nil
}

// SessionExists ...
func (s *Svc) SessionExists(ctx context.Context, uid string) bool {
	sess, err := s.st.GetSession(ctx, uid)
	if err != nil {
		return false
	}
	if time.Until(sess.Expiry.Time) > 0 {
		return true
	}
	return false
}

// GetOrgProfile by DSA Code ..
func (s *Svc) GetProfileByDsaCode(ctx context.Context, dsaCode string) (*storage.OrgProfile, error) {
	log := logging.FromContext(ctx)

	org, err := s.st.GetProfileByDsaCode(ctx, dsaCode)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "org not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "profile failed")
	}
	return org, nil
}
