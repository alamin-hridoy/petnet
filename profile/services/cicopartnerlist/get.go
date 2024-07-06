package cicopartnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetCICOPartnerList(ctx context.Context, req *spb.GetCICOPartnerListRequest) (*spb.GetCICOPartnerListResponse, error) {
	svc, err := s.core.GetCICOPartnerList(ctx, req)
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to Get CICO Partner List")
	}
	return svc, nil
}
