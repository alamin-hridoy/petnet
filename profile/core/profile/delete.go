package profile

import (
	"context"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func (s *Svc) DeleteUserProfile(ctx context.Context, req *ppb.DeleteUserProfileRequest) (ppb.DeleteUserProfileResponse, error) {
	err := s.st.DeleteUserProfile(ctx, req.UserID)
	if err != nil {
		return ppb.DeleteUserProfileResponse{}, err
	}
	return ppb.DeleteUserProfileResponse{}, nil
}
