package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) DeletePartnerList(ctx context.Context, req *spb.DeletePartnerListRequest) (*spb.DeletePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	Stype, err := s.st.DeletePartnerList(ctx, req.Stype)
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("PartnerList doesn't exists")
			return nil, status.Error(codes.NotFound, "failed to delete PartnerList record")
		}
		logging.WithError(err, log).Error("delete PartnerList")
		return nil, status.Error(codes.Internal, "failed to delete PartnerList record")
	}
	return &spb.DeletePartnerListResponse{
		Stype: Stype,
	}, nil
}
