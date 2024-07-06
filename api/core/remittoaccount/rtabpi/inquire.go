package rtabpi

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Inquire(ctx context.Context, req *rta.RTAInquireRequest) (res *rta.RTAInquireResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.remitAcc.RTABPIInquire(ctx, rtai.RTABPIInquireRequest{
		ReferenceNumber: req.GetReferenceNumber(),
		LocationID:      req.GetLocationID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA BPI inquire failed.")
		return nil, handleBPIError(err)
	}

	res = &rta.RTAInquireResponse{
		Message: rs.Message,
		Result: &rta.RTAInquireResult{
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
