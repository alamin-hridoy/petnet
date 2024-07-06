package email

import (
	"context"

	eml "brank.as/petnet/profile/integrations/email"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Mailer interface {
	OnboardingReminder(email string, orgID string, userID string) error
	DsaServiceRequestNotification(req eml.DsaServiceRequestNotificationForm) error
}

type Svc struct {
	mailer Mailer
}

func New(mailer Mailer) *Svc {
	return &Svc{mailer: mailer}
}

func (s *Svc) SendOnboardingReminder(ctx context.Context, email string, orgID string, userID string) error {
	log := logging.FromContext(ctx)

	if err := s.mailer.OnboardingReminder(email, orgID, userID); err != nil {
		logging.WithError(err, log).Error("send onboarding reminder")
		return status.Error(codes.Internal, "failed to send onboarding email reminder")
	}
	return nil
}

func (s *Svc) SendDsaServiceRequestNotification(ctx context.Context, req eml.DsaServiceRequestNotificationForm) error {
	log := logging.FromContext(ctx)

	if err := s.mailer.DsaServiceRequestNotification(req); err != nil {
		logging.WithError(err, log).Error("service request status")
		return status.Error(codes.Internal, "failed to send service request status email reminder")
	}
	return nil
}
