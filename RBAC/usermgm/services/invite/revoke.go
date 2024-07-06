package invite

import (
	"context"
	"database/sql"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"

	ipb "brank.as/rbac/gunk/v1/invite"
)

// Revoke a platform application
func (h *Svc) Revoke(ctx context.Context, req *ipb.RevokeRequest) (*ipb.RevokeResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.Revoke")
	log.Trace("request received")

	err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, validation.Length(1, 0)))
	if err != nil {
		log.WithError(err).Error("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	usr, err := h.store.GetUserByID(ctx, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("fetch invited user")
		return nil, status.Error(codes.PermissionDenied, "invalid user account")
	}
	if err := h.userRevoke(ctx, log, *usr); err != nil {
		return nil, err
	}

	return &ipb.RevokeResponse{}, nil
}

func (h *Svc) userRevoke(ctx context.Context, log *logrus.Entry, usr storage.User) error {
	// check if the application was already been revoked
	if usr.InviteStatus == storage.Revoked {
		log.Error("invalid request")
		return status.Errorf(codes.Internal, "the user was already revoked")
	}

	usr.Deleted = sql.NullTime{Time: time.Now(), Valid: true}
	usr.InviteStatus = storage.Revoked
	if _, err := h.store.UpdateUserByID(ctx, usr); err != nil {
		logging.WithError(err, log).Error("failed to update user")
		return err
	}

	return nil
}
