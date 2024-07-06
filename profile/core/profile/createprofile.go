package profile

import (
	"context"
	"database/sql"
	"time"

	"brank.as/petnet/profile/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CreateOrgProfile ...
func (s *Svc) CreateOrgProfile(ctx context.Context, p storage.OrgProfile) (id string, err error) {
	if !p.DateApplied.Valid {
		p.DateApplied = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	}
	id, err = s.st.CreateOrgProfile(ctx, &p)
	if err != nil {
		return "", err
	}
	return id, nil
}

// UpdateOrgProfile ...
func (s *Svc) UpdateOrgProfile(ctx context.Context, p storage.OrgProfile) (id string, err error) {
	id, err = s.st.UpdateOrgProfile(ctx, &p)
	if err != nil {
		return "", err
	}
	return id, nil
}

// CreateUserProfile ...
func (s *Svc) CreateUserProfile(ctx context.Context, p storage.UserProfile) (id string, err error) {
	id, err = s.st.CreateUserProfile(ctx, &p)
	if err != nil {
		if err == storage.Conflict {
			return "", status.Error(codes.AlreadyExists, "account already exists")
		}
		return "", err
	}
	return id, nil
}

// UpdateUserProfile ...
func (s *Svc) UpdateUserProfile(ctx context.Context, p storage.UserProfile) (id string, err error) {
	id, err = s.st.UpdateUserProfile(ctx, &p)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Svc) UpdateOrgProfileUserID(ctx context.Context, p storage.UpdateOrgProfileOrgIDUserID) (id string, err error) {
	id, err = s.st.UpdateOrgProfileUserID(ctx, &p)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *Svc) UpdateUserProfileByOrgID(ctx context.Context, p storage.UpdateOrgProfileOrgIDUserID) (id string, err error) {
	id, err = s.st.UpdateUserProfileByOrgID(ctx, &p)
	if err != nil {
		return "", err
	}
	return id, nil
}
