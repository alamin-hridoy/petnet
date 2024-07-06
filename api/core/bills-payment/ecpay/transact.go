package ecpay

import (
	"context"
	"strconv"

	bpi "brank.as/petnet/api/integration/bills-payment"
	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Transact(ctx context.Context, req *bp.BPTransactRequest) (res *bp.BPTransactResponse, err error) {
	log := logging.FromContext(ctx).WithField("method", "bills-payment.ecpay.transact")
	defer func() {
		_, err := util.RecordBillPayment(ctx, s.st, req, res, err)
		if err != nil {
			logging.WithError(err, log).Error("unable to record bill payment")
		}
	}()
	billID, err := strconv.Atoi(req.GetBillID())
	if err != nil {
		logging.WithError(err, log).Error("Cant convert bill id")
	}
	amount, err := strconv.Atoi(req.GetAmount())
	if err != nil {
		logging.WithError(err, log).Error("Cant convert amount")
	}
	serviceCharge, err := strconv.Atoi(req.GetServiceCharge())
	if err != nil {
		logging.WithError(err, log).Error("Cant convert bill id")
	}
	rs, err := s.billAcc.BillsPaymentEcpayTransact(ctx, bpi.BillsPaymentEcpayTransactRequest{
		BillID:                billID,
		BillerTag:             req.GetBillerTag(),
		TrxDate:               req.GetTrxDate(),
		UserID:                req.GetUserID(),
		RemoteUserID:          req.GetRemoteUserID(),
		CustomerID:            strconv.Itoa(int(req.GetCustomerID())),
		LocationID:            req.GetLocationID(),
		RemoteLocationID:      req.GetRemoteLocationID(),
		LocationName:          req.GetLocationName(),
		Coy:                   req.GetCoy(),
		CurrencyID:            req.GetCurrencyID(),
		FormType:              req.GetFormType(),
		FormNumber:            req.GetFormNumber(),
		AccountNumber:         req.GetAccountNumber(),
		Identifier:            req.GetIdentifier(),
		Amount:                amount,
		ServiceCharge:         serviceCharge,
		TotalAmount:           int(req.GetTotalAmount()),
		ClientReferenceNumber: req.GetClientReferenceNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment ecpay transact failed.")
		return nil, handleEcpayError(err)
	}

	res = &bp.BPTransactResponse{
		Code:    rs.Code,
		Message: rs.Message,
		Result: &bp.BPTransactResult{
			Status:          rs.Result.Status,
			Message:         rs.Result.Message,
			ServiceCharge:   int32(rs.Result.ServiceCharge),
			Timestamp:       rs.Result.Timestamp,
			ReferenceNumber: rs.Result.ReferenceNumber,
		},
		RemcoID: int32(rs.RemcoID),
	}

	return res, nil
}
