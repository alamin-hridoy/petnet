package microinsurance

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/storage"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
	"brank.as/petnet/serviceutil/logging"
)

// GetReprint ...
func (s *MICoreSvc) GetReprint(ctx context.Context, req *migunk.GetReprintRequest) (*migunk.Insurance, error) {
	log := logging.FromContext(ctx)
	res, err := s.storage.GetMicroInsuranceHistoryByTraceNumber(ctx, req.TraceNumber)
	if err != nil {
		if err == storage.ErrNotFound {
			return nil, coreerror.NewCoreError(codes.NotFound, coreerror.MsgNotFound)
		}

		logging.WithError(err, log).Error("getting microinsurance history")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	var ins migunk.Insurance
	err = json.Unmarshal(res.InsuranceDetails, &ins)
	if err != nil {
		logging.WithError(err, log).WithField("traceNo", req.TraceNumber).Error("unmarshaling insurance")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return &ins, nil
}
