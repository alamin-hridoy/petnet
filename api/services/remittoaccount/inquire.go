package remittoaccount

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) RTAInquire(ctx context.Context, req *rta.RTAInquireRequest) (*rta.RTAInquireResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	_, err := s.validators[pn].RTAInquireValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	_, err = s.rtaStore.RTAInquire(ctx, req, pn)
	if err != nil {
		logging.WithError(err, log).Error("failed to get inquire")
		return nil, util.HandleServiceErr(err)
	}

	return &rta.RTAInquireResponse{}, nil
}

func (s *MBRtaVal) RTAInquireValidate(ctx context.Context, req *rta.RTAInquireRequest) (*rta.RTAInquireRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.LocationID, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *UBRtaVal) RTAInquireValidate(ctx context.Context, req *rta.RTAInquireRequest) (*rta.RTAInquireRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *BPIRtaVal) RTAInquireValidate(ctx context.Context, req *rta.RTAInquireRequest) (*rta.RTAInquireRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.LocationID, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}
