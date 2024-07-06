package multipay

import (
	"context"
	"strconv"

	bpi "brank.as/petnet/api/integration/bills-payment"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) TransactInquire(ctx context.Context, req *bp.BPTransactInquireRequest) (res *bp.BPTransactInquireResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BillsPaymentMultiPayBillerInquire(ctx, bpi.BillsPaymentMultiPayBillerInquireRequest{
		AccountNumber: req.GetAccountNumber(),
		Amount:        int(req.GetAmount()),
		ContactNumber: req.GetContactNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment multipay transact inquire failed.")
		return nil, handleMultipayError(err)
	}

	res = &bp.BPTransactInquireResponse{
		Status: int32(rs.Status),
		Reason: rs.Reason,
		Result: &bp.BPTransactInquireResult{
			Amount:        strconv.Itoa(rs.Data.Amount),
			AccountNumber: rs.Data.AccountNumber,
			Biller:        rs.Data.Biller,
		},
	}

	return res, nil
}
