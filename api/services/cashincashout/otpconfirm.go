package cashincashout

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	cico "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) CiCoOTPConfirm(ctx context.Context, req *cico.CiCoOTPConfirmRequest) (*cico.CiCoOTPConfirmResponse, error) {
	if err := s.CiCoOTPConfirmValidate(ctx, req); err != nil {
		return nil, util.HandleServiceErr(err)
	}
	res, err := s.core.CiCoOTPConfirm(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}
	return res, err
}

func (s *Svc) CiCoOTPConfirmValidate(ctx context.Context, req *cico.CiCoOTPConfirmRequest) error {
	log := logging.FromContext(ctx)
	if req == nil {
		return status.Error(codes.InvalidArgument, "please provide required fields")
	}
	if req.OTPPayload == nil {
		return status.Error(codes.InvalidArgument, "otp payload is required")
	}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.PartnerCode, validation.Required),
		validation.Field(&req.PetnetTrackingno, validation.Required),
		validation.Field(&req.TrxDate, validation.Required),
		validation.Field(&req.OTP, validation.Required),
		validation.Field(&req.Provider, validation.Required),
		validation.Field(&req.OTPPayload, validation.Required, validation.By(func(interface{}) error {
			r := req.OTPPayload
			return validation.ValidateStruct(r,
				validation.Field(&r.CommandID, validation.Required),
				validation.Field(&r.Payload, validation.Required),
			)
		})),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
