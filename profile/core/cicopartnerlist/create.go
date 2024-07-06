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

func (s *Svc) CreateCICOPartnerList(ctx context.Context, req *spb.CreateCICOPartnerListRequest) (*spb.CreateCICOPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.CreateCICOPartnerList(ctx, &storage.CICOPartnerList{
		Stype:   req.GetCICOPartnerList().GetStype(),
		Name:    req.GetCICOPartnerList().GetName(),
		Created: req.GetCICOPartnerList().GetCreated().AsTime(),
		Updated: req.GetCICOPartnerList().GetUpdated().AsTime(),
		Status:  req.GetCICOPartnerList().GetStatus(),
	})
	if err != nil {
		if err == storage.Conflict {
			logging.WithError(err, log).Error("cico partner already exists")
			return nil, status.Error(codes.AlreadyExists, "failed to Create CICO Partner List")
		}
		logging.WithError(err, log).Error("store cico partner")
		return nil, status.Error(codes.Internal, "failed to Create CICO Partner List")
	}

	return &spb.CreateCICOPartnerListResponse{
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
