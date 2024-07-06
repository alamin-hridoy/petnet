package apitransactiontype

import (
	"context"

	apt "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) CreateApiKeyTransactionType(ctx context.Context, req *apt.CreateApiKeyTransactionTypeRequest) (*apt.CreateApiKeyTransactionTypeResponse, error) {
	id, err := s.st.InsertApiKeyTransactionType(ctx, &storage.ApiKeyTransactionType{
		UserID:          req.GetUserID(),
		OrgID:           req.GetOrgID(),
		ClientID:        req.GetClientID(),
		Environment:     req.GetEnvironment(),
		TransactionType: req.GetTransactionType().String(),
	})
	if err != nil {
		return nil, err
	}
	return &apt.CreateApiKeyTransactionTypeResponse{
		ApiTransactionTypes: &apt.ApiKeyTransactionType{
			UserID: id,
		},
	}, nil
}
