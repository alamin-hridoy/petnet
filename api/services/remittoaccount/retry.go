package remittoaccount

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"

	rta "brank.as/petnet/gunk/drp/v1/remittoaccount"
)

func (s *Svc) RTARetry(ctx context.Context, req *rta.RTARetryRequest) (*rta.RTARetryResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetPartner()
	_, err := s.validators[pn].RTARetryValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}
	_, err = s.rtaStore.RTARetry(ctx, req, pn)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return &rta.RTARetryResponse{}, nil
}

func (s *MBRtaVal) RTARetryValidate(ctx context.Context, req *rta.RTARetryRequest) (*rta.RTARetryRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.ID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.PrincipalAmount, required),
		validation.Field(&req.FormNumber, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *BPIRtaVal) RTARetryValidate(ctx context.Context, req *rta.RTARetryRequest) (*rta.RTARetryRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.ID, required),
		validation.Field(&req.LocationID, required),
		validation.Field(&req.FormNumber, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}

func (s *UBRtaVal) RTARetryValidate(ctx context.Context, req *rta.RTARetryRequest) (*rta.RTARetryRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ReferenceNumber, required),
		validation.Field(&req.ID, required)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return req, nil
}
