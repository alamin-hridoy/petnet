package cashincashout

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	cio "brank.as/petnet/gunk/drp/v1/cashincashout"
	"google.golang.org/grpc/codes"
)

func (s *Svc) CiCoValidate(ctx context.Context, req *cio.CiCoValidateRequest) (res *cio.CiCoValidateResponse, er error) {
	ci, err := s.ph.CicoValidate(ctx, perahub.CicoValidateRequest{
		PartnerCode: req.GetPartnerCode(),
		Trx: perahub.CicoValidateTrx{
			Provider:        req.Trx.GetProvider(),
			ReferenceNumber: req.Trx.GetReferenceNumber(),
			TrxType:         req.Trx.GetTrxType(),
			PrincipalAmount: int(req.Trx.GetPrincipalAmount()),
		},
		Customer: perahub.CicoValidateCustomer{
			CustomerID:        req.Customer.GetCustomerID(),
			CustomerFirstname: req.Customer.GetCustomerFirstname(),
			CustomerLastname:  req.Customer.GetCustomerLastname(),
			CurrAddress:       req.Customer.GetCurrAddress(),
			CurrBarangay:      req.Customer.GetCurrBarangay(),
			CurrCity:          req.Customer.GetCurrCity(),
			CurrProvince:      req.Customer.GetCurrProvince(),
			CurrCountry:       req.Customer.GetCurrCountry(),
			BirthDate:         req.Customer.GetBirthDate(),
			BirthPlace:        req.Customer.GetBirthPlace(),
			BirthCountry:      req.Customer.GetBirthCountry(),
			ContactNo:         req.Customer.GetContactNo(),
			IDType:            req.Customer.GetIDType(),
			IDNumber:          req.Customer.GetIDNumber(),
		},
	})
	if err != nil {
		return nil, handleCiCoError(err)
	}

	if ci == nil || ci.Result == nil {
		return nil, coreerror.NewCoreError(codes.NotFound, "not found")
	}

	return &cio.CiCoValidateResponse{
		Code:    int32(ci.Code),
		Message: ci.Message,
		Result: &cio.CicoValidateResult{
			PetnetTrackingno:   ci.Result.PetnetTrackingno,
			TrxDate:            ci.Result.TrxDate,
			TrxType:            ci.Result.TrxType,
			Provider:           ci.Result.Provider,
			ProviderTrackingno: ci.Result.ProviderTrackingno,
			ReferenceNumber:    ci.Result.ReferenceNumber,
			PrincipalAmount:    int32(ci.Result.PrincipalAmount),
			Charges:            int32(ci.Result.Charges),
			TotalAmount:        int32(ci.Result.TotalAmount),
			Timestamp:          ci.Result.Timestamp,
		},
	}, nil
}
