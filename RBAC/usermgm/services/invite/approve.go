package invite

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "brank.as/rbac/gunk/v1/invite"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/storage"
)

// Approve an invite
func (s *Svc) Approve(ctx context.Context, req *pb.ApproveRequest) (*pb.ApproveResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.invite.Approve")
	log.Trace("request received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, is.UUIDv4),
	); err != nil {
		logging.WithError(err, log).Info("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	usr, err := s.store.GetUserByID(ctx, req.GetID())
	if err != nil {
		logging.WithError(err, log).Error("failed to get user created")
		return nil, status.Error(codes.NotFound, "user not found")
	}
	org, err := s.plt.GetOrgByID(ctx, usr.OrgID)
	if err != nil {
		logging.WithError(err, log).Error("failed to get organization created")
		return nil, status.Error(codes.NotFound, "org not found")
	}
	if err := s.plt.ActivateOrg(ctx, org.ID); err != nil {
		logging.WithError(err, log).Error("org activation")
		return nil, status.Error(codes.Aborted, "failed to activate record")
	}

	org.Active = true
	org, err = s.store.UpdateOrgByID(ctx, *org)
	if err != nil {
		logging.WithError(err, log).Error("failed to update organization")
		return nil, status.Error(codes.Internal, "activation failed")
	}

	usr.EmailVerified = true
	usr.InviteStatus = storage.Approved
	if _, err := s.store.UpdateUserByID(ctx, *usr); err != nil {
		logging.WithError(err, log).Error("failed to update user")
		return nil, err
	}

	ap := email.Approved{
		CompanyName: org.OrgName,
		OrgID:       org.ID,
	}
	log.Trace("sending approved email")
	err = s.mailer.Approved(usr.Email, ap)
	if err != nil {
		const errMsg = "failed to send approved email"
		logging.WithError(err, log).Error(errMsg)
		if s.env != "development" {
			return nil, status.Error(codes.Internal, errMsg)
		}
	}
	return &pb.ApproveResponse{}, nil
}
