package rtaub

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) Retry(ctx context.Context, req *rta.RTARetryRequest) (res *rta.RTARetryResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.remitAcc.RTAUBRetry(ctx, rtai.RTAUBRetryRequest{
		ReferenceNumber: req.GetReferenceNumber(),
		ID:              req.GetID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA UB retry failed.")
		return nil, handleUBError(err)
	}

	res = &rta.RTARetryResponse{
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
			Created:         &timestamppb.Timestamp{},
			Updated:         &timestamppb.Timestamp{},
		},
	}

	return res, nil
}
