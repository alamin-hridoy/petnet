package mfa

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) EnableMFA(ctx context.Context, req *mpb.EnableMFARequest) (*mpb.EnableMFAResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.mfa.enablemfa")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.Type, validation.By(func(interface{}) error {
			if mpb.MFA_name[int32(req.GetType())] == "" {
				return validation.NewError("invalid_type", "mfa type invalid or not supported")
			}
			return nil
		})),
		validation.Field(&req.Source, validation.By(func(interface{}) error {
			s := req.GetSource()
			switch req.GetType() {
			case mpb.MFA_EMAIL:
				return validation.Validate(s, is.EmailFormat)
			case mpb.MFA_SMS:
				return validation.Validate(s, is.Int.Error("use only numbers in phone"))
			case mpb.MFA_CODE:
				return validation.Validate(s, is.Int.Error("PIN must be numerical"))
			case mpb.MFA_RECOVERY, mpb.MFA_TOTP:
			}
			return nil
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := s.a.RegisterMFA(ctx, core.MFA{
		UserID: req.UserID,
		Type:   req.Type.String(),
		Source: req.Source,
	})
	if err != nil {
		logging.WithError(err, log).Error("storage mfa registration")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "mfa activation failed")
	}

	return &mpb.EnableMFAResponse{
		ID:             m.MFAID,
		InitializeCode: m.Source,
		EventID:        m.ConfirmID,
	}, nil
}
