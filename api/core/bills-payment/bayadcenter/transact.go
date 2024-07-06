package bayadcenter

import (
	"context"
	"strconv"

	bpi "brank.as/petnet/api/integration/bills-payment"
	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Transact(ctx context.Context, req *bp.BPTransactRequest) (res *bp.BPTransactResponse, err error) {
	log := logging.FromContext(ctx).WithField("method", "bills-payment.bayadcenter.transact")
	defer func() {
		_, err := util.RecordBillPayment(ctx, s.st, req, res, err)
		if err != nil {
			logging.WithError(err, log).Error("unable to record bill payment")
		}
	}()
	rs, err := s.billAcc.BCTransact(ctx, bpi.BCTransactRequest{
		UserID:                  req.GetUserID(),
		CustomerID:              int(req.GetCustomerID()),
		LocationID:              req.GetLocationID(),
		LocationName:            req.GetLocationName(),
		Coy:                     req.GetCoy(),
		CallbackURL:             req.GetCallbackURL(),
		BillID:                  req.GetBillID(),
		BillerTag:               req.GetBillerTag(),
		BillerName:              req.GetBillerName(),
		TrxDate:                 req.GetTrxDate(),
		Amount:                  req.GetAmount(),
		ServiceCharge:           req.GetServiceCharge(),
		PartnerCharge:           req.GetPartnerCharge(),
		TotalAmount:             int(req.GetTotalAmount()),
		Identifier:              req.GetIdentifier(),
		AccountNumber:           req.GetAccountNumber(),
		PaymentMethod:           req.GetPaymentMethod(),
		ClientReferenceNumber:   req.GetClientReferenceNumber(),
		ReferenceNumber:         req.GetReferenceNumber(),
		ValidationNumber:        req.GetValidationNumber(),
		ReceiptValidationNumber: req.GetReceiptValidationNumber(),
		TpaID:                   req.GetTpaID(),
		CurrencyID:              req.GetCurrencyID(),
		FormType:                req.GetFormType(),
		FormNumber:              req.GetFormType(),
		OtherInfo: bpi.BCOtherInfo{
			LastName:        req.OtherInfo.GetLastName(),
			FirstName:       req.OtherInfo.GetFirstName(),
			MiddleName:      req.OtherInfo.GetMiddleName(),
			PaymentType:     req.OtherInfo.GetPaymentType(),
			Course:          req.OtherInfo.GetCourse(),
			TotalAssessment: req.OtherInfo.GetTotalAssessment(),
			SchoolYear:      req.OtherInfo.GetSchoolYear(),
			Term:            req.OtherInfo.GetTerm(),
		},
		Type: req.GetType(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment bayadcenter transact failed.")
		return nil, handleBayadcenterError(err)
	}

	res = &bp.BPTransactResponse{
		Code:    strconv.Itoa(rs.Code),
		Message: rs.Message,
		Result: &bp.BPTransactResult{
			Status:          rs.Result.Status,
			Message:         rs.Result.Message,
			Timestamp:       rs.Result.Timestamp,
			ReferenceNumber: rs.Result.ReferenceNumber,
			TransactionID:   rs.Result.TransactionID,
			ClientReference: rs.Result.ClientReference,
			BillerReference: rs.Result.BillerReference,
			PaymentMethod:   rs.Result.PaymentMethod,
			Amount:          rs.Result.Amount,
			OtherCharges:    rs.Result.OtherCharges,
			Details:         rs.Result.Details,
			CreatedAt:       rs.Result.CreatedAt,
		},
		RemcoID: int32(rs.RemcoID),
	}

	return res, nil
}
