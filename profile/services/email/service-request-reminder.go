package fees

import (
	"context"
	"strings"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	epb "brank.as/petnet/gunk/dsa/v1/email"
	eml "brank.as/petnet/profile/integrations/email"
)

func (s *Svc) SendDsaServiceRequestNotification(ctx context.Context, req *epb.SendDsaServiceRequestNotificationRequest) (*epb.SendDsaServiceRequestNotificationResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
		validation.Field(&req.Status, validation.Required),
		validation.Field(&req.Remark, validation.Required),
		validation.Field(&req.ServiceName, validation.Required),
		validation.Field(&req.PartnerNames, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ptnr := strings.Join(req.GetPartnerNames(), ", ")
	if err := s.es.SendDsaServiceRequestNotification(ctx, eml.DsaServiceRequestNotificationForm{
		Email:        req.GetEmail(),
		Status:       req.GetStatus().String(),
		Remark:       req.GetRemark(),
		PartnerNames: ptnr,
	}); err != nil {
		return nil, status.Error(codes.Internal, "failed to send service request reminder")
	}

	return &epb.SendDsaServiceRequestNotificationResponse{}, nil
}
