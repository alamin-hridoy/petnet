package userauth

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	apb "brank.as/rbac/gunk/v1/authenticate"
	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) RetryMFA(ctx context.Context, req *apb.RetryMFARequest) (*apb.Session, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.userauth.retrymfa")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.MFAEventID, validation.Required, is.UUIDv4),
		validation.Field(&req.Subject, validation.Required, is.UUIDv4),
		validation.Field(&req.ClientID, validation.Required, is.ASCII),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	usr, err := s.auth.NewMFA(ctx, core.AuthCredential{
		AuthClientID: req.GetClientID(),
		MFA: &core.MFAChallenge{
			EventID: req.GetMFAEventID(),
			UserID:  req.GetSubject(),
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("resend mfa")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to send mfa token")
	}

	return &apb.Session{
		UserID:     usr.ID,
		OrgID:      usr.OrgID,
		MFAEventID: usr.EventID,
		MFAType:    mpb.MFA(mpb.MFA_value[usr.MFA]),
		Attempt:    int32(usr.MFATrial),
	}, err
}
