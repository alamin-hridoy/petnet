package mfa

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/logging"

	mpb "brank.as/rbac/gunk/v1/mfa"
)

func (s *Svc) GetRegisteredMFA(ctx context.Context, req *mpb.GetRegisteredMFARequest) (*mpb.GetRegisteredMFAResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.mfa.enablemfa")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	m, err := s.a.ListMFA(ctx, req.UserID)
	if err != nil {
		logging.WithError(err, log).Error("storage mfa registration")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "mfa activation failed")
	}

	ts := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}
	mfa := make([]*mpb.MFAEntry, len(m))
	for i, mf := range m {
		mfa[i] = &mpb.MFAEntry{
			ID:       mf.MFAID,
			Type:     mpb.MFA(mpb.MFA_value[mf.Type]),
			Source:   mf.Source,
			Enabled:  ts(mf.Confirmed),
			Disabled: ts(mf.Revoked),
			Updated:  ts(mf.Updated),
		}
	}

	return &mpb.GetRegisteredMFAResponse{
		UserID: req.UserID,
		MFA:    mfa,
	}, nil
}
