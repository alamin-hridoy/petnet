package apitransactiontype

import (
	"context"

	apt "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListUserAPIKeyTransactionType(ctx context.Context, req *apt.ListUserAPIKeyTransactionTypeRequest) (*apt.ListUserAPIKeyTransactionTypeResponse, error) {
	res, err := s.st.ListUserAPIKeyTransactionType(ctx, req.OrgID, req.UserID)
	if err != nil {
		return nil, err
	}
	var aptList []*apt.ApiKeyTransactionType
	for _, v := range res {
		aptList = append(aptList, &apt.ApiKeyTransactionType{
			ID:              v.ID,
			UserID:          v.UserID,
			OrgID:           v.OrgID,
			ClientID:        v.ClientID,
			Environment:     v.Environment,
			TransactionType: apt.TransactionType(apt.TransactionType_value[v.TransactionType]),
			Created:         timestamppb.New(v.Created),
		})
	}
	return &apt.ListUserAPIKeyTransactionTypeResponse{
		ApiTransactionType: aptList,
	}, nil
}
