package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	uapb "brank.as/rbac/gunk/v1/user"
)

func (h *Handler) GetUser(ctx context.Context, req *uapb.GetUserRequest) (*uapb.GetUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.user.GetUser")
	log.Trace("request received")

	if req.GetID() == "" {
		return nil, status.Error(codes.InvalidArgument, "missing id")
	}

	usr, err := h.acct.GetUserByID(ctx, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("no user found")
		return nil, status.Error(codes.NotFound, "user doesn't exist")
	}
	org, err := h.acct.GetOrgByID(ctx, usr.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("no org found")
		return nil, status.Error(codes.NotFound, "user doesn't exist")
	}

	resp := &uapb.GetUserResponse{
		User: &uapb.User{
			ID:           usr.ID,
			OrgID:        org.ID,
			OrgName:      org.OrgName,
			FirstName:    usr.FirstName,
			LastName:     usr.LastName,
			Email:        usr.Email,
			InviteStatus: usr.InviteStatus,
			Created:      tspb.New(usr.Created),
			Updated:      tspb.New(usr.Updated),
		},
	}
	return resp, nil
}
