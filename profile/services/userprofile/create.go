package userprofile

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/profile/storage"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func (h *Svc) CreateUserProfile(ctx context.Context, req *ppb.CreateUserProfileRequest) (*ppb.CreateUserProfileResponse, error) {
	pf := req.GetProfile()
	if err := validation.ValidateStruct(pf,
		validation.Field(&pf.UserID, validation.Required, is.UUID),
		validation.Field(&pf.OrgID, validation.Required, is.UUID),
		validation.Field(&pf.Email, validation.Required, is.Email),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pid, err := h.core.CreateUserProfile(ctx, storage.UserProfile{
		OrgID:  pf.GetOrgID(),
		UserID: pf.GetUserID(),
		Email:  pf.GetEmail(),
	})
	if err != nil {
		if status.Code(err) == codes.AlreadyExists {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to create profile")
	}
	return &ppb.CreateUserProfileResponse{ID: pid}, nil
}
