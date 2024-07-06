package mfa

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) RetryMFA(ctx context.Context, req *mpb.RetryMFARequest) (*mpb.RetryMFAResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.mfa.retrymfa")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.EventID, validation.Required, is.UUIDv4),
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := s.a.RestartMFA(ctx, core.MFAChallenge{
		EventID: req.EventID,
		UserID:  req.UserID,
	})
	if err != nil {
		logging.WithError(err, log).Error("storage mfa retry")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "mfa retry failed")
	}

	ts := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}
	src := make([]*mpb.MFAEntry, len(m.Sources))
	for i, s := range m.Sources {
		src[i] = &mpb.MFAEntry{
			ID:       s.MFAID,
			Type:     mpb.MFA(mpb.MFA_value[m.Type]),
			Source:   s.Source,
			Enabled:  ts(s.Confirmed),
			Disabled: ts(s.Revoked),
			Updated:  ts(s.Updated),
		}
	}

	return &mpb.RetryMFAResponse{
		EventID: m.EventID,
		Sources: src,
		Value:   m.Token,
		Attempt: int32(m.Attempt),
	}, nil
}
