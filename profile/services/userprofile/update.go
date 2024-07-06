package userprofile

import (
	"context"
	"database/sql"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
	"brank.as/petnet/profile/storage"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (h *Svc) UpdateUserProfile(ctx context.Context, req *ppb.UpdateUserProfileRequest) (*ppb.UpdateUserProfileResponse, error) {
	pf := req.GetProfile()
	if err := validation.ValidateStruct(pf,
		validation.Field(&pf.ID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	pid, err := h.core.UpdateUserProfile(ctx, storage.UserProfile{
		ID:             pf.GetID(),
		ProfilePicture: pf.GetProfilePicture(),
		Deleted:        sql.NullTime{Time: pf.GetDeleted().AsTime(), Valid: pf.GetDeleted().IsValid()},
	})
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to update profile")
	}
	return &ppb.UpdateUserProfileResponse{ID: pid}, nil
}
