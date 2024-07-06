package mfa

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/logging"

	mpb "brank.as/rbac/gunk/v1/mfa"
	"brank.as/rbac/usermgm/core"
)

func (s *Svc) ExternalMFA(ctx context.Context, req *mpb.ExternalMFARequest) (*mpb.ExternalMFAResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.mfa.external")
	log.Trace("received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.EventID, validation.Required, is.UUIDv4),
		validation.Field(&req.Value, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := s.a.ExternalMFA(ctx, core.MFAChallenge{
		EventID:    req.EventID,
		ExternalID: req.SourceID,
		Token:      req.Value,
	})
	if err != nil {
		logging.WithError(err, log).Error("record external")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update event")
	}

	return &mpb.ExternalMFAResponse{EventID: m.EventID, Updated: tspb.Now()}, nil
}
