package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreatePartnerList(ctx context.Context, req *spb.CreatePartnerListRequest) (*spb.CreatePartnerListResponse, error) {
	if err := validation.ValidateStruct(req.PartnerList,
		validation.Field(&req.PartnerList.Stype, validation.Required),
		validation.Field(&req.PartnerList.Name, validation.Required),
		validation.Field(&req.PartnerList.Status, validation.Required),
		validation.Field(&req.PartnerList.ServiceName, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.CreatePartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Create Partner List")
	}
	return res, nil
}
