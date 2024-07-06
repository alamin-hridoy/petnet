package invite

import (
	"context"
	"fmt"
	"time"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	ipb "brank.as/rbac/gunk/v1/invite"
)

var sgTz = time.FixedZone("Asia/Manila", 8*3600)

// Resend resends the invitation for user ready to be activated
func (h *Svc) Resend(ctx context.Context, req *ipb.ResendRequest) (*ipb.ResendResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.Resend")
	log.Debug("request received")

	// TODO: Validate fields
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, is.UUIDv4),
	); err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	usr, err := h.store.GetUserByID(ctx, req.ID)
	if err != nil {
		logging.WithError(err, log).Error("fetch invited user")
		return nil, status.Error(codes.PermissionDenied, "invalid user account")
	}

	if err := h.userResend(ctx, *usr, req.CustomEmailData); err != nil {
		return nil, err
	}
	return &ipb.ResendResponse{}, nil
}

func (h *Svc) orgResend(ctx context.Context, usr storage.User) error {
	log := logging.FromContext(ctx).WithField("method", "service.invite.orgResend")
	// id := hydra.ClientID(ctx)

	// usr, err := h.store.GetUserByID(id)
	// if err != nil {
	// 	logging.WithError(err, log).Error("fetch inviting admin")
	// 	return status.Error(codes.PermissionDenied, "invalid user account")
	// }

	// check if the application is ready to be resent
	statuses := map[string]struct{}{
		storage.Invited: {}, storage.Expired: {}, storage.InProgress: {}, storage.Revoked: {},
	}
	if !validateInviteStatus(statuses, usr.InviteStatus) {
		log.WithField("status", usr.InviteStatus).Info("invalid request")
		return status.Errorf(codes.Internal, "cannot resend because the application status is %s", usr.InviteStatus)
	}

	usr.InviteStatus = storage.Invited
	usr.InviteExpiry = time.Now().AddDate(0, 0, 3) // invitation expires after 3 days
	for {
		usr.InviteCode = random.InvitationCode(16)
		if _, err := h.store.UpdateUserByID(ctx, usr); err != nil {
			if err == storage.InvCodeExists {
				continue
			}
			logging.WithError(err, log).Error("failed to update organization")
			return err
		}
		break
	}

	if err := h.store.ReviveOrgByID(ctx, usr.ID); err != nil {
		logging.WithError(err, log).Error("failed to revive organization")
		return err
	}

	// todo(robin): enable ones new implementation is in place
	exp := time.Now().In(sgTz).AddDate(0, 0, 3)
	inv := email.Invite{
		Username:   usr.FirstName + " " + usr.LastName,
		UserEmail:  usr.Email,
		Duration:   "3 days",
		ExpiryDate: fmt.Sprintf("%s %d", exp.Month().String(), exp.Day()),
	}
	log.Trace("resending invite email")
	if err := h.mailer.InviteUser(inv.UserEmail, usr.InviteCode, inv); err != nil {
		const errMsg = "failed to resend invite email"
		logging.WithError(err, log).Error(errMsg)
		return status.Error(codes.Internal, errMsg)
	}
	return nil
}

func (h *Svc) userResend(ctx context.Context, usr storage.User, customEmailData map[string]string) error {
	log := logging.FromContext(ctx).WithField("method", "service.invite.userResend")

	// check if the application is ready to be resent
	statuses := map[string]struct{}{storage.InviteSent: {}, storage.Expired: {}, storage.Revoked: {}}
	if !validateInviteStatus(statuses, usr.InviteStatus) {
		log.Error("invalid request")
		return status.Errorf(codes.Internal, "cannot resend because the application status is %s", usr.InviteStatus)
	}

	usr.InviteStatus = storage.InviteSent
	usr.InviteExpiry = time.Now().AddDate(0, 0, 3) // invitation expires after 3 days
regenerateCode:
	usr.InviteCode = random.InvitationCode(16)
	if _, err := h.store.UpdateUserByID(ctx, usr); err != nil {
		if err == storage.InvCodeExists {
			goto regenerateCode
		}
		logging.WithError(err, log).Error("failed to update user")
		return err
	}
	if err := h.store.ReviveUserByID(ctx, usr.ID); err != nil {
		logging.WithError(err, log).Error("failed to revive user")
		return err
	}

	// todo(robin): enable ones new implementation is in place
	invitingAdminID := hydra.ClientID(ctx)
	invitingAdmin, err := h.store.GetUserByID(ctx, invitingAdminID)
	if err != nil {
		logging.WithError(err, log).Error("fetch inviting admin")
		return status.Error(codes.PermissionDenied, "invalid user account")
	}

	exp := time.Now().In(sgTz).AddDate(0, 0, 3)
	inv := email.Invite{
		Username:        invitingAdmin.FirstName + " " + invitingAdmin.LastName,
		UserEmail:       usr.Email,
		Duration:        "3 days",
		ExpiryDate:      fmt.Sprintf("%s %d", exp.Month().String(), exp.Day()),
		CustomEmailData: customEmailData,
	}

	log.Trace("resending invite email")
	if err := h.mailer.InviteUser(inv.UserEmail, usr.InviteCode, inv); err != nil {
		const errMsg = "failed to resend invite email"
		logging.WithError(err, log).Error(errMsg)
		return status.Error(codes.Internal, errMsg)
	}
	return nil
}
