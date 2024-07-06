package bayadcenter

import (
	"context"

	bpi "brank.as/petnet/api/integration/bills-payment"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) TransactInquire(ctx context.Context, req *bp.BPTransactInquireRequest) (res *bp.BPTransactInquireResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BCTransactInquire(ctx, bpi.BCTransactInquireRequest{
		Code:            req.GetCode(),
		ClientReference: req.GetClientReference(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment bayadcenter transact inquire failed.")
		return nil, handleBayadcenterError(err)
	}
	res = &bp.BPTransactInquireResponse{
		Code:    int32(rs.Code),
		Message: rs.Message,
		Result: &bp.BPTransactInquireResult{
			TransactionID:   rs.Result.TransactionID,
			ReferenceNumber: rs.Result.ReferenceNumber,
			ClientReference: rs.Result.ClientReference,
			BillerReference: rs.Result.BillerReference,
			PaymentMethod:   rs.Result.PaymentMethod,
			Amount:          rs.Result.Amount,
			OtherCharges:    rs.Result.OtherCharges,
			Status:          rs.Result.Status,
			Message: &bp.BPTransactInquireMessage{
				Header:  rs.Result.Message.Header,
				Message: rs.Result.Message.Message,
				Footer:  rs.Result.Message.Header,
			},
			Details:   rs.Result.Details,
			CreatedAt: rs.Result.CreatedAt,
		},
		RemcoID: int32(rs.RemcoID),
	}

	return res, nil
}
