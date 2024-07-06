package bills_payment

import (
	"context"

	"brank.as/petnet/api/util"
	bp "brank.as/petnet/gunk/drp/v1/bills-payment"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) BPBillerList(ctx context.Context, req *bp.BPBillerListRequest) (*bp.BPBillerListResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	if pn == "" {
		return nil, util.HandleServiceErr(status.Error(codes.InvalidArgument, "partner is required."))
	}

	_, err := s.validators[pn].BPBillerListValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.bpStore.BPBillerList(ctx, req, pn)
	if err != nil {
		logging.WithError(err, log).Error("failed to get inquire")
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *ECPBPVal) BPBillerListValidate(ctx context.Context, req *bp.BPBillerListRequest) (*bp.BPBillerListRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Partner, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	logging.FromContext(ctx).Info("No validation required for Ecpay")
	return req, nil
}

func (s *BYCBPVal) BPBillerListValidate(ctx context.Context, req *bp.BPBillerListRequest) (*bp.BPBillerListRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Partner, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	logging.FromContext(ctx).Info("No validation required for Bayadcenter")
	return req, nil
}

func (s *MLPBPVal) BPBillerListValidate(ctx context.Context, req *bp.BPBillerListRequest) (*bp.BPBillerListRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Partner, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	logging.FromContext(ctx).Info("No validation required for Multipay")
	return req, nil
}
