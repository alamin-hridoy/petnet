package userprofile

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func (h *Svc) DeleteUserProfile(ctx context.Context, req *ppb.DeleteUserProfileRequest) (*ppb.DeleteUserProfileResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err := h.core.DeleteUserProfile(ctx, &ppb.DeleteUserProfileRequest{
		UserID: req.UserID,
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to delete profile")
	}
	return &ppb.DeleteUserProfileResponse{}, nil
}
