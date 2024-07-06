package mfa

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) DisableMFA(ctx context.Context, req *mpb.DisableMFARequest) (*mpb.DisableMFAResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.mfa.enablemfa")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.MFAID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := s.a.DisableMFA(ctx, core.MFA{UserID: req.UserID, MFAID: req.MFAID})
	if err != nil {
		logging.WithError(err, log).Error("storage mfa disable")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "mfa update failed")
	}

	return &mpb.DisableMFAResponse{Disabled: tspb.New(m.Revoked)}, nil
}
