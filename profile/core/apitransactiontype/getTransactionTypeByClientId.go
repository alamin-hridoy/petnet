package apitransactiontype

import (
	"context"

	apt "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetTransactionTypeByClientId(ctx context.Context, req *apt.GetTransactionTypeByClientIdRequest) (*apt.GetTransactionTypeByClientIdResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.GetTransactionTypeByClientId(ctx, req.ClientID)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "api transaction type not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "api transaction type failed")
	}
	return &apt.GetTransactionTypeByClientIdResponse{
		Environment:     res.Environment,
		TransactionType: res.TransactionType,
	}, nil
}
