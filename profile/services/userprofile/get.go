package userprofile

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/profile/storage"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func (h *Svc) GetUserProfile(ctx context.Context, req *ppb.GetUserProfileRequest) (*ppb.GetUserProfileResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pf, err := h.core.GetUserProfile(ctx, req.UserID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get profile")
	}
	return &ppb.GetUserProfileResponse{Profile: storageToProto(pf)}, nil
}

func storageToProto(spf *storage.UserProfile) *ppb.Profile {
	ppf := &ppb.Profile{
		ID:             spf.ID,
		UserID:         spf.UserID,
		OrgID:          spf.OrgID,
		Email:          spf.Email,
		ProfilePicture: spf.ProfilePicture,
		Created:        tspb.New(spf.Created),
		Updated:        tspb.New(spf.Updated),
		Deleted:        tspb.New(spf.Deleted.Time),
	}

	return ppf
}

func (h *Svc) GetUserProfileByEmail(ctx context.Context, req *ppb.GetUserProfileByEmailRequest) (*ppb.GetUserProfileByEmailResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pf, err := h.core.GetUserProfileByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get profile by email")
	}
	return &ppb.GetUserProfileByEmailResponse{Profile: storageToProto(pf)}, nil
}
