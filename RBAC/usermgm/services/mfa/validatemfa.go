package mfa

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
)

func (s *Svc) ValidateMFA(ctx context.Context, req *mpb.ValidateMFARequest) (*mpb.ValidateMFAResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.mfa.validatemfa")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.EventID, is.UUIDv4),
		validation.Field(&req.Type, validation.By(func(interface{}) error {
			if mpb.MFA_name[int32(req.GetType())] == "" {
				return validation.NewError("invalid_type", "mfa type invalid or not supported")
			}
			return nil
		})),
		validation.Field(&req.Token, validation.Required, validation.Length(1, 0)),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := s.a.MFAuth(ctx, core.MFAChallenge{
		EventID: req.GetEventID(),
		UserID:  req.GetUserID(),
		Type:    req.GetType().String(),
		Token:   req.GetToken(),
	})
	if err != nil {
		logging.WithError(err, log).Error("mfa validation")
		return nil, err
	}

	return &mpb.ValidateMFAResponse{
		EventID:   m.EventID,
		Validated: tspb.Now(),
	}, nil
}
