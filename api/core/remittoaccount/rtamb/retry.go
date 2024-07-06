package rtamb

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) Retry(ctx context.Context, req *rta.RTARetryRequest) (res *rta.RTARetryResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.remitAcc.RTAMetrobankRetry(ctx, rtai.RTAMetrobankRetryRequest{
		ReferenceNumber: req.GetReferenceNumber(),
		ID:              req.GetID(),
		LocationID:      int(req.GetLocationID()),
		PrincipalAmount: req.GetPrincipalAmount(),
		FormNumber:      req.GetFormNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA MB retry failed.")
		return nil, handleMBError(err)
	}

	res = &rta.RTARetryResponse{
		Message: rs.Message,
		Result: &rta.RTAPaymentResult{
			Message: rs.Result.Message,
		},
		BankCode: rs.BankCode,
	}

	return res, nil
}
