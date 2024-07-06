package rtamb

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/util"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Payment(ctx context.Context, req *rta.RTAPaymentRequest) (res *rta.RTAPaymentResponse, err error) {
	log := logging.FromContext(ctx)
	orgID := phmw.GetDSAOrgID(ctx)
	defer func() {
		_, err := util.RecordRTA(ctx, orgID, s.st, req, res, err)
		if err != nil {
			log.Error(err)
		}
	}()

	rs, err := s.remitAcc.RTAMetrobankPayment(ctx, rtai.RTAMetrobankPaymentRequest{
		ReferenceNumber:             req.GetReferenceNumber(),
		TrxDate:                     req.GetTrxDate(),
		AccountNumber:               req.GetAccountNumber(),
		Currency:                    req.GetCurrency(),
		ServiceCharge:               req.GetServiceCharge(),
		Remarks:                     req.GetRemarks(),
		Particulars:                 req.GetParticulars(),
		MerchantName:                req.GetMerchantName(),
		BankID:                      req.GetBankID(),
		LocationID:                  int(req.GetLocationID()),
		UserID:                      int(req.GetUserID()),
		CurrencyID:                  req.GetCurrencyID(),
		CustomerID:                  req.GetCustomerID(),
		FormType:                    req.GetFormType(),
		FormNumber:                  req.GetFormNumber(),
		TrxType:                     req.GetTrxType(),
		RemoteLocationID:            int(req.GetRemoteLocationID()),
		RemoteUserID:                int(req.GetRemoteUserID()),
		BillerName:                  req.GetBillerName(),
		TrxTime:                     req.GetTrxTime(),
		TotalAmount:                 req.GetTotalAmount(),
		AccountName:                 req.GetAccountName(),
		BeneficiaryAddress:          req.GetBeneficiaryAddress(),
		BeneficiaryBirthdate:        req.GetBeneficiaryBirthdate(),
		BeneficiaryCity:             req.GetBeneficiaryCity(),
		BeneficiaryCivil:            req.GetBeneficiaryCivil(),
		BeneficiaryCountry:          req.GetBeneficiaryCountry(),
		BeneficiaryCustomertype:     req.GetBeneficiaryCustomertype(),
		BeneficiaryFirstname:        req.GetBeneficiaryFirstname(),
		BeneficiaryLastname:         req.GetBeneficiaryLastname(),
		BeneficiaryMiddlename:       req.GetBeneficiaryMiddlename(),
		BeneficiaryTin:              req.GetBeneficiaryTin(),
		BeneficiarySex:              req.GetBeneficiarySex(),
		BeneficiaryState:            req.GetBeneficiaryState(),
		CurrencyCodePrincipalAmount: req.GetCurrencyCodePrincipalAmount(),
		PrincipalAmount:             req.GetPrincipalAmount(),
		RecordType:                  req.GetRecordType(),
		RemitterAddress:             req.GetRemitterAddress(),
		RemitterBirthdate:           req.GetRemitterBirthdate(),
		RemitterCity:                req.GetRemitterCity(),
		RemitterCivil:               req.GetRemitterCivil(),
		RemitterCountry:             req.GetRemitterCountry(),
		RemitterCustomerType:        req.GetRemitterCustomerType(),
		RemitterFirstname:           req.GetRemitterFirstname(),
		RemitterGender:              req.GetRemitterGender(),
		RemitterID:                  int(req.GetRemitterID()),
		RemitterLastname:            req.GetRemitterLastname(),
		RemitterMiddlename:          req.GetRemitterMiddlename(),
		RemitterState:               req.GetRemitterState(),
		SettlementMode:              req.GetSettlementMode(),
		BeneZipCode:                 req.GetBeneZipCode(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA MB payment failed.")
		return nil, handleMBError(err)
	}

	res = &rta.RTAPaymentResponse{
		Message: rs.Message,
		Result: &rta.RTAPaymentResult{
			Message: rs.Result.Message,
		},
	}

	return res, nil
}
