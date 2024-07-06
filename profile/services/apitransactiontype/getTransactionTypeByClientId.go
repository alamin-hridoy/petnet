package apitransactiontype

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/transactiontype"
)

func (s *Svc) GetTransactionTypeByClientId(ctx context.Context, req *ppb.GetTransactionTypeByClientIdRequest) (*ppb.GetTransactionTypeByClientIdResponse, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ClientID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.GetTransactionTypeByClientId(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
