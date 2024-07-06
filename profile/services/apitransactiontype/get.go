package apitransactiontype

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/serviceutil/logging"

	ppb "brank.as/petnet/gunk/dsa/v2/transactiontype"
)

func (s *Svc) GetAPITransactionType(ctx context.Context, req *ppb.GetAPITransactionTypeRequest) (*ppb.ApiKeyTransactionType, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, required, is.UUID),
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Environment, required),
		validation.Field(&req.TransactionType, required),
	); err != nil {
		logging.WithError(err, log).Error("get api transaction type validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.GetAPITransactionType(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
