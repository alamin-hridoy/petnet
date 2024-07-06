package rtamb

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Inquire(ctx context.Context, req *rta.RTAInquireRequest) (res *rta.RTAInquireResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.remitAcc.RTAMetrobankInquire(ctx, rtai.RTAMetrobankInquireRequest{
		ReferenceNumber: req.GetReferenceNumber(),
		LocationID:      req.GetLocationID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA MB inquire failed.")
		return nil, handleMBError(err)
	}

	res = &rta.RTAInquireResponse{
		Message: rs.Message,
		Result: &rta.RTAInquireResult{
			Description: rs.Result,
		},
		BankCode: rs.BankCode,
	}

	return res, nil
}
