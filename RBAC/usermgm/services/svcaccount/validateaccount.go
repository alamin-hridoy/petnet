package svcaccount

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

func (h *Svc) ValidateAccount(ctx context.Context, req *sapb.ValidateAccountRequest) (*sapb.ValidateAccountResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.svcaccount.ValidateAccount")
	log.Trace("request received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.ClientID, validation.Required, is.ASCII),
	); err != nil {
		logging.WithError(err, log).Error("validation")
		return nil, status.Error(codes.InvalidArgument, "id invalid")
	}

	acct, err := h.store.ValidateSvcAccount(ctx, req.GetClientID())
	if err != nil {
		return nil, err
	}
	return &sapb.ValidateAccountResponse{
		OrgID:       acct.OrgID,
		Environment: acct.Environment,
		ClientName:  acct.ClientName,
	}, nil
}
