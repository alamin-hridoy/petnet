package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DisableMultiplePartnerList(ctx context.Context, req *spb.DisableMultiplePartnerListRequest) (*spb.DisableMultiplePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stypes, err := s.st.DisableMultiplePartnerList(ctx, req.Stypes, req.DisableReason, req.UpdatedBy)
	if err != nil {
		logging.WithError(err, log).WithField("Stype", req.Stypes).Error("DisableMultiple Partner List")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "partner List not found")
		}
		return nil, err
	}
	return &spb.DisableMultiplePartnerListResponse{
		Stypes: Stypes,
	}, nil
}
