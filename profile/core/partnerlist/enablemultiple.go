package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) EnableMultiplePartnerList(ctx context.Context, req *spb.EnableMultiplePartnerListRequest) (*spb.EnableMultiplePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stypes, err := s.st.EnableMultiplePartnerList(ctx, req.Stypes, req.UpdatedBy)
	if err != nil {
		logging.WithError(err, log).WithField("Stype", req.Stypes).Error("EnableMultiple Partner List")
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "partner List not found")
		}
		return nil, err
	}
	return &spb.EnableMultiplePartnerListResponse{
		Stypes: Stypes,
	}, nil
}
