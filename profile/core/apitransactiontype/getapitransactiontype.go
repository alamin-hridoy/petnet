package apitransactiontype

import (
	"context"

	apt "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetAPITransactionType(ctx context.Context, req *apt.GetAPITransactionTypeRequest) (*apt.ApiKeyTransactionType, error) {
	log := logging.FromContext(ctx)

	res, err := s.st.GetAPITransactionType(ctx, &storage.ApiKeyTransactionType{
		UserID:          req.UserID,
		OrgID:           req.OrgID,
		Environment:     req.Environment,
		TransactionType: req.TransactionType.String(),
	})
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "api transaction type not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "api transaction type failed")
	}
	return &apt.ApiKeyTransactionType{
		ID:              res.ID,
		UserID:          res.UserID,
		OrgID:           res.OrgID,
		ClientID:        res.ClientID,
		Environment:     res.Environment,
		TransactionType: apt.TransactionType(apt.TransactionType_value[res.TransactionType]),
	}, nil
}
