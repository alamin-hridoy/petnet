package cicopartnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DeleteCICOPartnerList(ctx context.Context, req *spb.DeleteCICOPartnerListRequest) (*spb.DeleteCICOPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stype, err := s.st.DeleteCICOPartnerList(ctx, req.Stype)
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("CICOPartnerList doesn't exists")
			return nil, status.Error(codes.NotFound, "failed to delete CICO PartnerList record")
		}
		logging.WithError(err, log).Error("delete CICO PartnerList")
		return nil, status.Error(codes.Internal, "failed to delete CICO PartnerList record")
	}
	return &spb.DeleteCICOPartnerListResponse{
		Stype: Stype,
	}, nil
}
