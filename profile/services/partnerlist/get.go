package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetPartnerList(ctx context.Context, req *spb.GetPartnerListRequest) (*spb.GetPartnerListResponse, error) {
	svc, err := s.core.GetPartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Get Partner List")
	}
	return svc, nil
}

func (s *Svc) GetDSAPartnerList(ctx context.Context, req *spb.DSAPartnerListRequest) (*spb.GetDSAPartnerListResponse, error) {
	svc, err := s.core.GetDSAPartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Get DSA Partner List")
	}
	return svc, nil
}

func (s *Svc) GetPartnerByStype(ctx context.Context, req *spb.GetPartnerByStypeRequest) (*spb.GetPartnerByStypeResponse, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Stype, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.GetPartnerByStype(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
