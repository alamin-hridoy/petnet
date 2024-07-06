package rtaub

import (
	"context"
	"strconv"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/util"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Payment(ctx context.Context, req *rta.RTAPaymentRequest) (res *rta.RTAPaymentResponse, err error) {
	log := logging.FromContext(ctx)
	trxType, _ := strconv.Atoi(req.GetTrxType())
	bnkID, _ := strconv.Atoi(req.GetBankID())
	cusID, _ := strconv.Atoi(req.GetCustomerID())
	orgID := phmw.GetDSAOrgID(ctx)
	defer func() {
		_, err := util.RecordRTA(ctx, orgID, s.st, req, res, err)
		if err != nil {
			log.Error(err)
		}
	}()

	rs, err := s.remitAcc.RTAUBCashin(ctx, rtai.RTAUBCashinRequest{
		TrxType:          trxType,
		BillerName:       req.GetBillerName(),
		BankID:           bnkID,
		LocationID:       string(req.GetLocationID()),
		UserID:           string(req.GetUserID()),
		RemoteLocationID: string(req.GetRemoteLocationID()),
		RemoteUserID:     string(req.GetRemoteUserID()),
		CurrencyID:       req.GetCurrencyID(),
		FormType:         req.GetFormType(),
		FormNumber:       req.GetFormNumber(),
		CustomerID:       cusID,
		ReferenceNumber:  req.GetReferenceNumber(),
		TrxDate:          req.GetTrxDate(),
		TrxTime:          req.GetTrxTime(),
		AccountNumber:    req.GetAccountNumber(),
		Currency:         req.GetCurrency(),
		PrincipalAmount:  req.GetPrincipalAmount(),
		ServiceCharge:    req.GetServiceCharge(),
		TotalAmount:      req.GetTotalAmount(),
		Remarks:          req.GetRemarks(),
		Particulars:      req.GetParticulars(),
		MerchantName:     req.GetMerchantName(),
		Notification:     req.GetNotification(),
		AccountName:      req.GetAccountName(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA UB payment failed.")
		return nil, handleUBError(err)
	}

	res = &rta.RTAPaymentResponse{
		Message: rs.Message,
		Result: &rta.RTAPaymentResult{
			Code:            rs.Result.Code,
			SenderRefID:     rs.Result.SenderRefID,
			State:           rs.Result.State,
			UUID:            rs.Result.UUID,
			Description:     rs.Result.Description,
			Type:            rs.Result.Type,
			Amount:          rs.Result.Amount,
			UbpTranID:       rs.Result.UbpTranID,
			TranRequestDate: rs.Result.TranRequestDate,
			TranFinacleDate: rs.Result.TranFinacleDate,
		},
	}

	return res, nil
}
