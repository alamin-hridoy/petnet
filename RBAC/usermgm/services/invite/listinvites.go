package invite

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	ipb "brank.as/rbac/gunk/v1/invite"
)

func (h *Svc) ListInvite(ctx context.Context, req *ipb.ListInviteRequest) (*ipb.ListInviteResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.ListInvite")
	log.Trace("request received")

	// TODO: Validate fields
	err := validation.ValidateStruct(req)
	if err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}
	orgID := hydra.OrgID(ctx)

	usrs, err := h.store.GetUsersByOrg(ctx, orgID)
	if err != nil {
		logging.WithError(err, log).Error("failed to create organization storage entry")
		return nil, status.Error(codes.Internal, err.Error())
	}
	resp := &ipb.ListInviteResponse{
		Invites: make([]*ipb.Invite, len(usrs)),
	}
	for i, usr := range usrs {
		resp.Invites[i] = &ipb.Invite{
			ID:           usr.ID,
			OrgID:        usr.OrgID,
			ContactEmail: usr.Email,
			Status:       usr.InviteStatus,
			Invited:      tspb.New(usr.Updated),
		}
	}
	return resp, nil
}
