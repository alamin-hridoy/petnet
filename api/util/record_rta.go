package util

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RecordRTA(ctx context.Context, orgID string, st *postgres.Storage, req *rta.RTAPaymentRequest, res *rta.RTAPaymentResponse, er error) (*storage.RemitToAccountHistory, error) {
	log := logging.FromContext(ctx)
	var errCode string
	var errMsg string
	var errType perahub.ErrType
	bps := storage.SuccessStatus
	if er != nil {
		bps = storage.FailStatus
		switch t := er.(type) {
		case *perahub.Error:
			errMsg = t.Msg
			if t.UnknownErr != "" {
				errMsg = t.UnknownErr
			}
			if t.Errors != nil {
				b, err := json.Marshal(t.Errors)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				errMsg = errMsg + ", " + string(b)
			}
			errCode = t.Code
			errType = t.Type
		default:
			errMsg = er.Error()
			errType = perahub.BillerError
		}
	}
	resDetails, err := json.Marshal(res)
	if err != nil {
		resDetails = []byte{}
	}

	bnkID, _ := strconv.Atoi(req.GetBankID())
	notification, _ := strconv.ParseBool(req.GetNotification())
	trxdt, err := time.Parse("2006-01-01", req.GetTrxDate())
	if err != nil {
		trxdt = time.Time{}
	}
	crtaHistory := storage.RemitToAccountHistory{
		OrgID:                       orgID,
		Partner:                     req.GetPartner(),
		ReferenceNumber:             req.GetReferenceNumber(),
		TrxDate:                     trxdt,
		AccountNumber:               req.GetAccountNumber(),
		Currency:                    req.GetCurrency(),
		ServiceCharge:               req.GetServiceCharge(),
		Remarks:                     req.GetRemarks(),
		Particulars:                 req.GetParticulars(),
		MerchantName:                req.GetMerchantName(),
		BankID:                      bnkID,
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
		BeneficiaryBirthDate:        req.GetBeneficiaryBirthdate(),
		BeneficiaryCity:             req.GetBeneficiaryCity(),
		BeneficiaryCivil:            req.GetBeneficiaryCivil(),
		BeneficiaryCountry:          req.GetBeneficiaryCountry(),
		BeneficiaryCustomerType:     req.GetBeneficiaryCustomertype(),
		BeneficiaryFirstName:        req.GetBeneficiaryFirstname(),
		BeneficiaryLastName:         req.GetBeneficiaryLastname(),
		BeneficiaryMiddleName:       req.GetBeneficiaryMiddlename(),
		BeneficiaryTin:              req.GetBeneficiaryTin(),
		BeneficiarySex:              req.GetBeneficiarySex(),
		BeneficiaryState:            req.GetBeneficiaryState(),
		CurrencyCodePrincipalAmount: req.GetCurrencyCodePrincipalAmount(),
		PrincipalAmount:             req.GetPrincipalAmount(),
		RecordType:                  req.GetRecordType(),
		RemitterAddress:             req.GetRemitterAddress(),
		RemitterBirthDate:           req.GetRemitterBirthdate(),
		RemitterCity:                req.GetRemitterCity(),
		RemitterCivil:               req.GetRemitterCivil(),
		RemitterCountry:             req.GetRemitterCountry(),
		RemitterCustomerType:        req.GetRemitterCustomerType(),
		RemitterFirstName:           req.GetRemitterFirstname(),
		RemitterGender:              req.GetRemitterGender(),
		RemitterID:                  int(req.GetRemitterID()),
		RemitterLastName:            req.GetRemitterLastname(),
		RemitterMiddleName:          req.GetRemitterMiddlename(),
		RemitterState:               req.GetRemitterState(),
		SettlementMode:              req.GetSettlementMode(),
		Notification:                notification,
		BeneZipCode:                 req.GetBeneZipCode(),
	}

	crtaHistory.Details = resDetails
	crtaHistory.ErrorCode = errCode
	crtaHistory.ErrorMessage = errMsg
	crtaHistory.ErrorType = string(errType)
	if errCode != "" {
		crtaHistory.ErrorTime = time.Now().String()
	}
	crtaHistory.TxnStatus = string(bps)
	rs, err := st.CreateRTAHistory(ctx, crtaHistory)
	if err != nil {
		logging.WithError(err, log).Error("creating remit to account transaction")
		if err == storage.Conflict {
			return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("this transaction already exists"))
		}
		return nil, err
	}
	return rs, nil
}

func UpdateRTA(ctx context.Context, orgID string, st *postgres.Storage, req *rta.RTARetryRequest, res *rta.RTARetryResponse, er error) (*storage.RemitToAccountHistory, error) {
	log := logging.FromContext(ctx)
	var errCode string
	var errMsg string
	var errType perahub.ErrType
	bps := storage.SuccessStatus
	if er != nil {
		bps = storage.FailStatus
		switch t := er.(type) {
		case *perahub.Error:
			errMsg = t.Msg
			if t.UnknownErr != "" {
				errMsg = t.UnknownErr
			}
			if t.Errors != nil {
				b, err := json.Marshal(t.Errors)
				if err != nil {
					log.Error(err)
					return nil, err
				}
				errMsg = errMsg + ", " + string(b)
			}
			errCode = t.Code
			errType = t.Type
		default:
			errMsg = er.Error()
			errType = perahub.BillerError
		}
	}
	resDetails, err := json.Marshal(res)
	if err != nil {
		resDetails = []byte{}
	}

	updRTaHis := storage.RemitToAccountHistory{
		Partner:         req.GetPartner(),
		ReferenceNumber: req.GetReferenceNumber(),
		ID:              orgID,
		LocationID:      int(req.GetLocationID()),
		PrincipalAmount: req.GetPrincipalAmount(),
		FormNumber:      req.GetFormNumber(),
	}
	updRTaHis.Details = resDetails
	updRTaHis.ErrorCode = errCode
	updRTaHis.ErrorMessage = errMsg
	updRTaHis.ErrorType = string(errType)
	if errCode != "" {
		updRTaHis.ErrorTime = time.Now().String()
	}
	updRTaHis.TxnStatus = string(bps)
	var rs *storage.RemitToAccountHistory
	var rserr error
	if req.ReferenceNumber != "" {
		rs, rserr = st.UpdateRTAHistory(ctx, updRTaHis)
		if rserr != nil {
			logging.WithError(rserr, log).Error("updating remit to account transaction")
			if rserr == storage.Conflict {
				return nil, status.Error(codes.AlreadyExists, fmt.Sprintf("this transaction, can not be update"))
			}
			return nil, rserr
		}
	}
	return rs, nil
}
