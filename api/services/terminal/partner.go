package terminal

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetPartnerByTxnID(ctx context.Context, req *tpb.GetPartnerByTxnIDRequest) (*tpb.GetPartnerByTxnIDResponse, error) {
	log := logging.FromContext(ctx)
	if req.GetTransactionID() == "" {
		log.Error("transaction id is required")
		return nil, util.HandleServiceErr(status.Error(codes.InvalidArgument, "transaction id is required"))
	}

	ptnr, err := s.remit.GetPartnerByTxnID(ctx, req.GetTransactionID())
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	if ptnr == "" {
		log.Error("partner doesn't exist")
		return nil, util.HandleServiceErr(status.Error(codes.NotFound, "partner doesn't exist"))
	}

	return &tpb.GetPartnerByTxnIDResponse{
		Partner: ptnr,
	}, nil
}
