package userprofile

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v1/user"
)

func (h *Svc) ListUserProfiles(ctx context.Context, req *ppb.ListUserProfilesRequest) (*ppb.ListUserProfilesResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pfs, err := h.core.GetUserProfiles(ctx, req.OrgID)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to get profiles")
	}
	ppfs := []*ppb.Profile{}
	for _, pf := range pfs {
		ppfs = append(ppfs, storageToProto(&pf))
	}
	return &ppb.ListUserProfilesResponse{Profiles: ppfs}, nil
}
