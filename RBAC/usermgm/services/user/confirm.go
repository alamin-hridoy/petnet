package user

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	mpb "brank.as/rbac/gunk/v1/mfa"
	upb "brank.as/rbac/gunk/v1/user"
)

func (s *Handler) ConfirmUpdate(ctx context.Context, req *upb.ConfirmUpdateRequest) (*upb.ConfirmUpdateResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.user.confirmupdate")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.MFAEventID, validation.Required, is.UUIDv4),
		validation.Field(&req.MFAType, validation.By(func(interface{}) error {
			return validation.Validate(mpb.MFA_name[mpb.MFA_value[req.MFAType.String()]],
				validation.Required)
		})),
		validation.Field(&req.MFAToken, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	mfa := core.MFAChallenge{
		EventID: req.GetMFAEventID(),
		UserID:  req.GetUserID(),
		Type:    mpb.MFA_name[mpb.MFA_value[req.GetMFAType().String()]],
		Token:   req.GetMFAToken(),
	}
	if err := s.usr.ConfirmPass(ctx, mfa); err != nil {
		logging.WithError(err, log).Error("core confirm pass")
		return nil, err
	}

	return &upb.ConfirmUpdateResponse{Updated: tspb.Now()}, nil
}
