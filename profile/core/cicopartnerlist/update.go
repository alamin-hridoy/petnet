package cicopartnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) UpdateCICOPartnerList(ctx context.Context, req *spb.UpdateCICOPartnerListRequest) (*spb.UpdateCICOPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.UpdateCICOPartnerList(ctx, &storage.CICOPartnerList{
		ID:      req.GetCICOPartnerList().GetID(),
		Stype:   req.GetCICOPartnerList().GetStype(),
		Name:    req.GetCICOPartnerList().GetName(),
		Created: req.GetCICOPartnerList().GetCreated().AsTime(),
		Updated: req.GetCICOPartnerList().GetUpdated().AsTime(),
		Status:  req.GetCICOPartnerList().GetStatus(),
	})
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("cico partner Not Found")
			return nil, status.Error(codes.NotFound, "failed to update cico partner record")
		}
		logging.WithError(err, log).Error("store cico partner")
		return nil, status.Error(codes.Internal, "failed to update cico partner record")
	}

	return &spb.UpdateCICOPartnerListResponse{
		CICOPartnerList: &spb.CICOPartnerList{
			ID:      res.ID,
			Stype:   res.Stype,
			Name:    res.Name,
			Created: timestamppb.New(res.Created),
			Updated: timestamppb.New(res.Updated),
			Status:  res.Status,
		},
	}, nil
}
