package microinsurance

import (
	"context"

	"brank.as/petnet/api/util"
	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/microinsurance"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
	"brank.as/petnet/serviceutil/logging"
)

// RetryTransaction ...
func (s *MICoreSvc) RetryTransaction(ctx context.Context, req *migunk.RetryTransactionRequest) (*migunk.Insurance, error) {
	log := logging.FromContext(ctx)
	res, err := s.cl.RetryTransaction(ctx, &microinsurance.RetryTransactionRequest{
		ID: req.ID,
	})
	if err != nil {
		logging.WithError(err, log).WithField("ID", req.ID).Error("microinsurance retry error")
		return nil, coreerror.ToCoreError(err)
	}

	ins := toInsurance(res)
	if ins == nil {
		log.WithField("ID", req.ID).Error("empty insurance response")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	_, recErr := util.UpdateMicroInsurance(ctx, s.storage, ins)
	if recErr != nil {
		logging.WithError(recErr, log).
			WithField("ID", req.ID).
			WithField("insurance", ins).
			Error("update mi history")
	}

	return ins, nil
}
