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

func (h *Svc) UpdateUserProfileByOrgID(ctx context.Context, req *ppb.UpdateUserProfileByOrgIDRequest) (*ppb.UpdateUserProfileByOrgIDResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OldOrgID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pid, err := h.core.UpdateUserProfileByOrgID(ctx, storage.UpdateOrgProfileOrgIDUserID{
		OldOrgID: req.GetOldOrgID(),
		NewOrgID: req.GetNewOrgID(),
		UserID:   req.GetUserID(),
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update profile")
	}
	return &ppb.UpdateUserProfileByOrgIDResponse{ID: pid}, nil
}
