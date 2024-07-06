package ecpay

import (
	"context"

	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) BillerList(ctx context.Context, req *bp.BPBillerListRequest) (res *bp.BPBillerListResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.billAcc.BillsPaymentEcpayBillerlist(ctx)
	if err != nil {
		logging.WithError(err, log).Error("Bills payment ecpay biller list failed.")
		return nil, handleEcpayError(err)
	}

	var bl []*bp.BPBillerListResult
	for _, v := range rs.Result {
		bl = append(bl, &bp.BPBillerListResult{
			BillerTag:         v.BillerTag,
			Description:       v.Description,
			FirstField:        v.FirstField,
			FirstFieldFormat:  v.FirstFieldFormat,
			FirstFieldWidth:   v.FirstFieldWidth,
			SecondField:       v.SecondField,
			SecondFieldFormat: v.SecondFieldFormat,
			SecondFieldWidth:  v.SecondFieldWidth,
			ServiceCharge:     int32(v.ServiceCharge),
		})
	}
	res = &bp.BPBillerListResponse{
		Code:    int32(rs.Code),
		Message: rs.Message,
		Result:  bl,
		RemcoID: int32(rs.RemcoID),
	}

	return res, nil
}
