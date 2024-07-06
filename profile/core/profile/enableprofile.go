package profile

import (
	"context"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func (s *Svc) EnableUserProfile(ctx context.Context, req *ppb.EnableUserProfileRequest) (*ppb.EnableUserProfileResponse, error) {
	err := s.st.EnableUserProfile(ctx, req.UserID)
	if err != nil {
		return &ppb.EnableUserProfileResponse{}, err
	}
	return &ppb.EnableUserProfileResponse{}, nil
}
