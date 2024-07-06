package apitransactiontype

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) ListUserAPIKeyTransactionType(ctx context.Context, req *ppb.ListUserAPIKeyTransactionTypeRequest) (*ppb.ListUserAPIKeyTransactionTypeResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, required, is.UUID),
		validation.Field(&req.OrgID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("list api transaction type validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.ListUserAPIKeyTransactionType(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
