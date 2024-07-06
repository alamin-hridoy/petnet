package invite

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"

	pb "brank.as/rbac/gunk/v1/invite"
)

// CancelInvite revokes an invite that was sent
func (h *Svc) CancelInvite(ctx context.Context, req *pb.CancelInviteRequest) (*pb.CancelInviteResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.CancelInvite")
	log.Trace("request received")

	// TODO: Validate fields
	err := validation.ValidateStruct(req)
	if err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	// check if the user was created
	usr, err := h.store.GetUserByID(ctx, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("failed to get user created")
		return nil, err
	}

	// check if the invite was already revoked
	if usr.InviteStatus == storage.Revoked {
		log.WithField("status", usr.InviteStatus).Info("invite already revoked")
		return &pb.CancelInviteResponse{}, nil
	}
	usr.InviteStatus = storage.Revoked
	usr, err = h.store.UpdateUserByID(ctx, *usr)
	if err != nil {
		logging.WithError(err, log).Error("failed to update user")
		return nil, err
	}

	return &pb.CancelInviteResponse{}, nil
}
