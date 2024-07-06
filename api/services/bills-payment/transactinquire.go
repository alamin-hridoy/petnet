package bills_payment

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) BPTransactInquire(ctx context.Context, req *bp.BPTransactInquireRequest) (*bp.BPTransactInquireResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	if pn == "" {
		return nil, util.HandleServiceErr(status.Error(codes.InvalidArgument, "partner is required."))
	}

	_, err := s.validators[pn].BPTransactInquireValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.bpStore.BPTransactInquire(ctx, req, pn)
	if err != nil {
		logging.WithError(err, log).Error("failed to get inquire")
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *ECPBPVal) BPTransactInquireValidate(ctx context.Context, req *bp.BPTransactInquireRequest) (*bp.BPTransactInquireRequest, error) {
	logging.FromContext(ctx).Info("No validation required for Ecpay")
	return req, nil
}

func (s *BYCBPVal) BPTransactInquireValidate(ctx context.Context, req *bp.BPTransactInquireRequest) (*bp.BPTransactInquireRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Code, required),
		validation.Field(&req.ClientReference, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *MLPBPVal) BPTransactInquireValidate(ctx context.Context, req *bp.BPTransactInquireRequest) (*bp.BPTransactInquireRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.AccountNumber, required),
		validation.Field(&req.Amount, required),
		validation.Field(&req.ContactNumber, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}
