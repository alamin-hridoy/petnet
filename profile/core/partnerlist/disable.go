package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DisablePartnerList(ctx context.Context, req *spb.DisablePartnerListRequest) (*spb.DisablePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stype, err := s.st.DisablePartnerList(ctx, req.Stype, req.DisableReason)
	if err != nil {
		logging.WithError(err, log).WithField("Stype", req.Stype).Error("Disable Partner List")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "partner List not found")
		}
		return nil, err
	}
	return &spb.DisablePartnerListResponse{
		Stype: Stype,
	}, nil
}
