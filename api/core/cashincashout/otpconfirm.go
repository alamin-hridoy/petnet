package cashincashout

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	cio "brank.as/petnet/gunk/drp/v1/cashincashout"
	"google.golang.org/grpc/codes"
)

func (s *Svc) CiCoOTPConfirm(ctx context.Context, req *cio.CiCoOTPConfirmRequest) (res *cio.CiCoOTPConfirmResponse, er error) {
	ci, err := s.ph.CicoOTPConfirm(ctx, perahub.CicoOTPConfirmRequest{
		PartnerCode:      req.GetPartnerCode(),
		PetnetTrackingno: req.GetPetnetTrackingno(),
		TrxDate:          req.GetTrxDate(),
		OTP:              req.OTP,
		OTPPayload: perahub.OTPpayload{
			CommandID: int(req.GetOTPPayload().GetCommandID()),
			Payload:   req.GetOTPPayload().GetPayload(),
		},
	})
	if err != nil {
		return nil, handleCiCoError(err)
	}

	if ci == nil || ci.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, "not found")
	}

	return &cio.CiCoOTPConfirmResponse{
		Code:    int32(ci.Code),
		Message: ci.Message,
		Result: &cio.CicoOTPConfirmResult{
			PartnerCode:        ci.Result.PartnerCode,
			Provider:           ci.Result.Provider,
			PetnetTrackingno:   ci.Result.PetnetTrackingno,
			TrxDate:            ci.Result.TrxDate,
			TrxType:            ci.Result.TrxType,
			ProviderTrackingno: ci.Result.ProviderTrackingno,
			ReferenceNumber:    ci.Result.ReferenceNumber,
			PrincipalAmount:    int32(ci.Result.PrincipalAmount),
			Charges:            int32(ci.Result.Charges),
			TotalAmount:        int32(ci.Result.TotalAmount),
		},
	}, nil
}
