package remittoaccount

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RTAPayment(ctx context.Context, req *rta.RTAPaymentRequest) (*rta.RTAPaymentResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	r, err := s.validators[pn].RTAPaymentValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	_, err = s.rtaStore.RTAPayment(ctx, r, pn)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return &rta.RTAPaymentResponse{}, nil
}

func (s *MBRtaVal) RTAPaymentValidate(ctx context.Context, req *rta.RTAPaymentRequest) (*rta.RTAPaymentRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.TrxDate, required),
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.Currency, required),
		validation.Field(&req.ServiceCharge, required),
		validation.Field(&req.Remarks, required),
		validation.Field(&req.Particulars, required),
		validation.Field(&req.MerchantName, required),
		validation.Field(&req.BankID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.UserID, required),
		validation.Field(&req.CurrencyID, required),
		validation.Field(&req.CustomerID, required),
		validation.Field(&req.FormType, required),
		validation.Field(&req.FormNumber, required),
		validation.Field(&req.TrxType, required),
		validation.Field(&req.RemoteLocationID, required),
		validation.Field(&req.RemoteUserID, required),
		validation.Field(&req.BillerName, required),
		validation.Field(&req.TrxTime, required),
		validation.Field(&req.TotalAmount, required),
		validation.Field(&req.AccountName, required),
		validation.Field(&req.BeneficiaryAddress, required),
		validation.Field(&req.BeneficiaryBirthdate, required),
		validation.Field(&req.BeneficiaryCity, required),
		validation.Field(&req.BeneficiaryCivil, required),
		validation.Field(&req.BeneficiaryCountry, required),
		validation.Field(&req.BeneficiaryCustomertype, required),
		validation.Field(&req.BeneficiaryFirstname, required),
		validation.Field(&req.BeneficiaryLastname, required),
		validation.Field(&req.BeneficiaryMiddlename, required),
		validation.Field(&req.BeneficiaryTin, required),
		validation.Field(&req.BeneficiarySex, required),
		validation.Field(&req.BeneficiaryState, required),
		validation.Field(&req.CurrencyCodePrincipalAmount, required),
		validation.Field(&req.PrincipalAmount, required),
		validation.Field(&req.RecordType, required),
		validation.Field(&req.RemitterAddress, required),
		validation.Field(&req.RemitterBirthdate, required),
		validation.Field(&req.RemitterCity, required),
		validation.Field(&req.RemitterCivil, required),
		validation.Field(&req.RemitterCountry, required),
		validation.Field(&req.RemitterCustomerType, required),
		validation.Field(&req.RemitterFirstname, required),
		validation.Field(&req.RemitterGender, required),
		validation.Field(&req.RemitterID, required),
		validation.Field(&req.RemitterLastname, required),
		validation.Field(&req.RemitterMiddlename, required),
		validation.Field(&req.RemitterState, required),
		validation.Field(&req.SettlementMode, required),
		validation.Field(&req.BeneZipCode, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return req, nil
}

func (s *UBRtaVal) RTAPaymentValidate(ctx context.Context, req *rta.RTAPaymentRequest) (*rta.RTAPaymentRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.TrxType, required),
		validation.Field(&req.BillerName, required),
		validation.Field(&req.BankID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.UserID, required),
		validation.Field(&req.RemoteLocationID, required),
		validation.Field(&req.RemoteUserID, required),
		validation.Field(&req.CurrencyID, required),
		validation.Field(&req.FormType, required),
		validation.Field(&req.FormNumber, required),
		validation.Field(&req.CustomerID, required),
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.TrxDate, required),
		validation.Field(&req.TrxTime, required),
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.Currency, required),
		validation.Field(&req.PrincipalAmount, required),
		validation.Field(&req.ServiceCharge, required),
		validation.Field(&req.TotalAmount, required),
		validation.Field(&req.Remarks, required),
		validation.Field(&req.Particulars, required),
		validation.Field(&req.MerchantName, required),
		validation.Field(&req.Notification, required),
		validation.Field(&req.AccountName, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return req, nil
}

func (s *BPIRtaVal) RTAPaymentValidate(ctx context.Context, req *rta.RTAPaymentRequest) (*rta.RTAPaymentRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.TrxDate, required),
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.Currency, required),
		validation.Field(&req.ServiceCharge, required),
		validation.Field(&req.Remarks, required),
		validation.Field(&req.Particulars, required),
		validation.Field(&req.MerchantName, required),
		validation.Field(&req.BankID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.UserID, required),
		validation.Field(&req.CurrencyID, required),
		validation.Field(&req.CustomerID, required),
		validation.Field(&req.FormType, required),
		validation.Field(&req.FormNumber, required),
		validation.Field(&req.TrxType, required),
		validation.Field(&req.RemoteLocationID, required),
		validation.Field(&req.RemoteUserID, required),
		validation.Field(&req.BillerName, required),
		validation.Field(&req.TrxTime, required),
		validation.Field(&req.TotalAmount, required),
		validation.Field(&req.AccountName, required),
		validation.Field(&req.BeneficiaryAddress, required),
		validation.Field(&req.BeneficiaryBirthdate, required),
		validation.Field(&req.BeneficiaryCity, required),
		validation.Field(&req.BeneficiaryCivil, required),
		validation.Field(&req.BeneficiaryCountry, required),
		validation.Field(&req.BeneficiaryCustomertype, required),
		validation.Field(&req.BeneficiaryFirstname, required),
		validation.Field(&req.BeneficiaryLastname, required),
		validation.Field(&req.BeneficiaryMiddlename, required),
		validation.Field(&req.BeneficiaryTin, required),
		validation.Field(&req.BeneficiarySex, required),
		validation.Field(&req.BeneficiaryState, required),
		validation.Field(&req.CurrencyCodePrincipalAmount, required),
		validation.Field(&req.PrincipalAmount, required),
		validation.Field(&req.RecordType, required),
		validation.Field(&req.RemitterAddress, required),
		validation.Field(&req.RemitterBirthdate, required),
		validation.Field(&req.RemitterCity, required),
		validation.Field(&req.RemitterCivil, required),
		validation.Field(&req.RemitterCountry, required),
		validation.Field(&req.RemitterCustomerType, required),
		validation.Field(&req.RemitterFirstname, required),
		validation.Field(&req.RemitterGender, required),
		validation.Field(&req.RemitterID, required),
		validation.Field(&req.RemitterLastname, required),
		validation.Field(&req.RemitterMiddlename, required),
		validation.Field(&req.RemitterState, required),
		validation.Field(&req.SettlementMode, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}
