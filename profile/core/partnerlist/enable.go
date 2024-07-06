package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) EnablePartnerList(ctx context.Context, req *spb.EnablePartnerListRequest) (*spb.EnablePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stype, err := s.st.EnablePartnerList(ctx, req.Stype)
	if err != nil {
		logging.WithError(err, log).WithField("Stype", req.Stype).Error("Enable Partner List")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "partner List not found")
		}
		return nil, err
	}
	return &spb.EnablePartnerListResponse{
		Stype: Stype,
	}, nil
}
