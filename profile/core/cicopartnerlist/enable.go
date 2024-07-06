package cicopartnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) EnableCICOPartnerList(ctx context.Context, req *spb.EnableCICOPartnerListRequest) (*spb.EnableCICOPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stype, err := s.st.EnableCICOPartnerList(ctx, req.Stype)
	if err != nil {
		logging.WithError(err, log).WithField("Stype", req.Stype).Error("Enable CICO Partner List")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "CICO partner List not found")
		}
		return nil, err
	}
	return &spb.EnableCICOPartnerListResponse{
		Stype: Stype,
	}, nil
}
