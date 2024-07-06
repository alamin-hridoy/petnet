package rtaub

import (
	"context"

	rtai "brank.as/petnet/api/integration/remittoaccount"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) Inquire(ctx context.Context, req *rta.RTAInquireRequest) (res *rta.RTAInquireResponse, err error) {
	log := logging.FromContext(ctx)
	rs, err := s.remitAcc.RTAUBInquire(ctx, rtai.RTAUBInquireRequest{
		ReferenceNumber: req.GetReferenceNumber(),
	})
	if err != nil {
		logging.WithError(err, log).Error("RTA UB inquire failed.")
		return nil, handleUBError(err)
	}

	res = &rta.RTAInquireResponse{
		Message: rs.Message,
		Result: &rta.RTAInquireResult{
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
			CreatedAt:       &timestamppb.Timestamp{},
			UpdatedAt:       &timestamppb.Timestamp{},
		},
	}

	return res, nil
}
