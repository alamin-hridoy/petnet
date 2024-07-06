package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdatePartnerList(ctx context.Context, req *spb.UpdatePartnerListRequest) (*spb.UpdatePartnerListResponse, error) {
	if err := validation.ValidateStruct(req.PartnerList,
		validation.Field(&req.PartnerList.Stype, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.UpdatePartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Update Partner List")
	}
	return res, nil
}
