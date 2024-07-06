package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) EnableMultiplePartnerList(ctx context.Context, req *spb.EnableMultiplePartnerListRequest) (*spb.EnableMultiplePartnerListResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Stypes, validation.Required),
		validation.Field(&req.UpdatedBy, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.EnableMultiplePartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to EnableMultiple Partner List")
	}
	return res, nil
}
