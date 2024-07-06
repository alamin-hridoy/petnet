package user

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	mpb "brank.as/rbac/gunk/v1/mfa"
	upb "brank.as/rbac/gunk/v1/user"
)

func (s *Handler) ChangePassword(ctx context.Context, req *upb.ChangePasswordRequest) (*upb.ChangePasswordResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.user.changepassword")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.OldPassword, validation.Required, validation.Length(1, 64)),
		validation.Field(&req.NewPassword, validation.Required, validation.Length(1, 64)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if hydra.ClientID(ctx) != req.GetUserID() {
		return nil, status.Errorf(codes.InvalidArgument, "user ID mismatch")
	}

	mfa, err := s.usr.ChangePass(ctx, req.GetUserID(), req.GetOldPassword(), req.GetNewPassword())
	if err != nil {
		logging.WithError(err, log).Error("change pass")
		return nil, err
	}
	if mfa != nil {
		return &upb.ChangePasswordResponse{
			MFAEventID: mfa.EventID,
			MFAType:    mpb.MFA(mpb.MFA_value[mfa.Type]),
		}, nil
	}

	return &upb.ChangePasswordResponse{
		Updated: tspb.Now(),
	}, nil
}
