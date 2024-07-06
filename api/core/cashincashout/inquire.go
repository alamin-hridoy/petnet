package cashincashout

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	cio "brank.as/petnet/gunk/drp/v1/cashincashout"
	"google.golang.org/grpc/codes"
)

func (s *Svc) CiCoInquire(ctx context.Context, req *cio.CiCoInquireRequest) (res *cio.CiCoInquireResponse, er error) {
	ci, err := s.ph.CicoInquire(ctx, perahub.CicoInquireRequest{
		PartnerCode:        req.GetPartnerCode(),
		Provider:           req.GetProvider(),
		TrxType:            req.GetTrxType(),
		ReferenceNumber:    req.GetReferenceNumber(),
		PetnetTrackingno:   req.GetPetnetTrackingno(),
		ProviderTrackingno: req.GetProviderTrackingno(),
		Message:            req.GetMessage(),
	})
	if err != nil {
		return nil, handleCiCoError(err)
	}

	if ci == nil || ci.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, "not found")
	}

	return &cio.CiCoInquireResponse{
		Code:    int32(ci.Code),
		Message: ci.Message,
		Result: &cio.CicoInquireResult{
			StatusMessage:      ci.Result.StatusMessage,
			PetnetTrackingno:   ci.Result.PetnetTrackingno,
			TrxType:            ci.Result.TrxType,
			ReferenceNumber:    ci.Result.ReferenceNumber,
			Amount:             ci.Result.Amount,
			ProviderTrackingno: ci.Result.ProviderTrackingno,
			Expiry:             ci.Result.Expiry,
			CustomerName:       ci.Result.CustomerName,
			CustomerFirstname:  ci.Result.CustomerFirstname,
			CustomerLastname:   ci.Result.CustomerLastname,
			MerchantID:         ci.Result.MerchantID,
			PartnerCode:        ci.Result.PartnerCode,
			AccountNumber:      ci.Result.AccountNumber,
			ServiceCharge:      int32(ci.Result.ServiceCharge),
			CreatedAt:          ci.Result.CreatedAt,
		},
	}, nil
}
