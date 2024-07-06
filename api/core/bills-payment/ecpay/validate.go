package ecpay

import (
	"context"

	bpi "brank.as/petnet/api/integration/bills-payment"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Validate(ctx context.Context, req *bp.BPValidateRequest) (res *bp.BPValidateResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BillsPaymentEcpayValidate(ctx, bpi.BillsPaymentEcpayValidateRequest{
		AccountNo:  req.GetAccountNo(),
		Identifier: req.GetIdentifier(),
		BillerTag:  req.GetBillerTag(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Bills payment ecpay validate failed.")
		return nil, handleEcpayError(err)
	}

	res = &bp.BPValidateResponse{
		Code:    int32(rs.Code),
		Message: rs.Message,
		Reason:  rs.Result,
		RemcoID: int32(rs.RemcoID),
	}

	return res, nil
}
