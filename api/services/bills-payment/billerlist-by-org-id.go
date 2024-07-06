package bills_payment

import (
	"context"

	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) BillsPaymentTransactList(ctx context.Context, req *bp.BillsPaymentTransactListRequest) (*bp.BillsPaymentTransactListResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.billspayment.BillsPaymentTransactList")
	r, err := s.BillsPaymentTransactListValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.bpStore.BillsPaymentTransactList(ctx, r)
	if err != nil {
		logging.WithError(err, log).Error("failed to get Bills Payment Transact List")
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) BillsPaymentTransactListValidate(ctx context.Context, req *bp.BillsPaymentTransactListRequest) (*bp.BillsPaymentTransactListRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	logging.FromContext(ctx).Info("Bills Payment Transact List validation")
	return req, nil
}
