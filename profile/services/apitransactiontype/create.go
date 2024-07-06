package apitransactiontype

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"brank.as/petnet/serviceutil/logging"
)

func (h *Svc) CreateApiKeyTransactionType(ctx context.Context, req *ppb.CreateApiKeyTransactionTypeRequest) (*ppb.CreateApiKeyTransactionTypeResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, required),
		validation.Field(&req.OrgID, required),
		validation.Field(&req.ClientID, required),
		validation.Field(&req.Environment, required),
		validation.Field(&req.TransactionType, required),
	); err != nil {
		logging.WithError(err, log).Error("create api key transaction type validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := h.core.CreateApiKeyTransactionType(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
