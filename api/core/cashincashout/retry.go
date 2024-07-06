package cashincashout

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	cio "brank.as/petnet/gunk/drp/v1/cashincashout"
	"google.golang.org/grpc/codes"
)

func (s *Svc) CiCoRetry(ctx context.Context, req *cio.CiCoRetryRequest) (res *cio.CiCoRetryResponse, er error) {
	ci, err := s.ph.CicoRetry(ctx, perahub.CicoRetryRequest{
		PartnerCode:      req.GetPartnerCode(),
		PetnetTrackingno: req.GetPetnetTrackingno(),
		TrxDate:          req.GetTrxDate(),
	})
	if err != nil {
		return nil, handleCiCoError(err)
	}

	if ci == nil || ci.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, "not found")
	}

	return &cio.CiCoRetryResponse{
		Code:    int32(ci.Code),
		Message: ci.Message,
		Result: &cio.CicoRetryResult{
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
			OTPPayload: &cio.OTPPayload{
				CommandID: int32(ci.Result.OTPPayload.CommandID),
				Payload:   ci.Result.OTPPayload.Payload,
			},
		},
	}, nil
}
