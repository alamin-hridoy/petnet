package cicopartnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DisableCICOPartnerList(ctx context.Context, req *spb.DisableCICOPartnerListRequest) (*spb.DisableCICOPartnerListResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Stype, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.DisableCICOPartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Disable CICO Partner List")
	}
	return res, nil
}
