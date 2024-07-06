package fees

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	epb "brank.as/petnet/gunk/dsa/v1/email"
)

func (s *Svc) SendOnboardingReminder(ctx context.Context, req *epb.SendOnboardingReminderRequest) (*epb.SendOnboardingReminderResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.OrgID, validation.Required, is.UUID),
		validation.Field(&req.UserID, validation.Required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.es.SendOnboardingReminder(ctx, req.Email, req.OrgID, req.UserID); err != nil {
		return nil, status.Error(codes.Internal, "failed to send onboarding reminder")
	}

	return &epb.SendOnboardingReminderResponse{}, nil
}
