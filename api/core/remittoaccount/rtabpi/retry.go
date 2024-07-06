package rtabpi

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Retry(ctx context.Context, req *rta.RTARetryRequest) (res *rta.RTARetryResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.remitAcc.RTABPIRetry(ctx, rtai.RTABPIRetryRequest{
		ReferenceNumber: req.GetReferenceNumber(),
		ID:              req.GetID(),
		LocationID:      int(req.GetLocationID()),
		FormNumber:      req.GetFormNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA BPI retry failed.")
		return nil, handleBPIError(err)
	}

	res = &rta.RTARetryResponse{
		Message: rs.Message,
		Result: &rta.RTAPaymentResult{
			BeneAmount:                     rs.Result.BeneAmount,
			BeneficiaryAddress:             rs.Result.BeneficiaryAddress,
			BeneficiaryBankAccountno:       rs.Result.BeneficiaryBankAccountno,
			BeneficiaryCity:                rs.Result.BeneficiaryCity,
			BeneficiaryCountry:             rs.Result.BeneficiaryCountry,
			BeneficiaryFirstName:           rs.Result.BeneficiaryFirstName,
			BeneficiaryLastName:            rs.Result.BeneficiaryLastName,
			BeneficiaryMiddleName:          rs.Result.BeneficiaryMiddleName,
			BeneficiaryStateOrProvince:     rs.Result.BeneficiaryStateOrProvince,
			BpiBranchCode:                  rs.Result.BpiBranchCode,
			CurrencyCodeOfFundingAmount:    rs.Result.CurrencyCodeOfFundingAmount,
			CurrencyCodeOfSettlementAmount: rs.Result.CurrencyCodeOfSettlementAmount,
			TxnDistributionDate:            rs.Result.TxnDistributionDate,
			FundingAmount:                  rs.Result.FundingAmount,
			Reason:                         rs.Result.Reason,
			RemitterCity:                   rs.Result.RemitterCity,
			RemitterCountry:                rs.Result.RemitterCountry,
			RemitterFirstName:              rs.Result.RemitterFirstName,
			RemitterLastName:               rs.Result.RemitterLastName,
			RemitterMiddleName:             rs.Result.RemitterMiddleName,
			RemitterStateOrProvince:        rs.Result.RemitterStateOrProvince,
			SettlementMode:                 rs.Result.SettlementMode,
			StatusCode:                     rs.Result.StatusCode,
			TransactionDate:                rs.Result.TransactionDate,
			TransactionReferenceNo:         rs.Result.TransactionReferenceNo,
		},
		BankCode: rs.BankCode,
	}

	return res, nil
}
