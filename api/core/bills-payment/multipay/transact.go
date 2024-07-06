package multipay

import (
	"context"

	bpi "brank.as/petnet/api/integration/bills-payment"
	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Transact(ctx context.Context, req *bp.BPTransactRequest) (res *bp.BPTransactResponse, err error) {
	log := logging.FromContext(ctx).WithField("method", "bills-payment.multipay.transact")
	defer func() {
		_, err := util.RecordBillPayment(ctx, s.st, req, res, err)
		if err != nil {
			logging.WithError(err, log).Error("unable to record bill payment")
		}
	}()
	rs, err := s.billAcc.BillsPaymentMultiPayTransact(ctx, bpi.BillsPaymentMultiPayTransactRequest{
		Amount: req.GetAmount(),
		Txnid:  req.GetTxnid(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment multipay transact failed.")
		return nil, handleMultipayError(err)
	}

	res = &bp.BPTransactResponse{
		Message: rs.Data.URL,
	}

	return res, nil
}
