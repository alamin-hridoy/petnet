package cicopartnerlist

import (
	"context"

	ppb "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetCICOPartnerList(ctx context.Context, req *ppb.GetCICOPartnerListRequest) (*ppb.GetCICOPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.GetCICOPartnerList(ctx, &storage.CICOPartnerList{
		ID:     req.GetID(),
		Stype:  req.GetStype(),
		Status: req.GetStatus(),
		Name:   req.GetName(),
	})
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("CICO PartnerList not found")
			return nil, status.Error(codes.NotFound, "CICO PartnerList not found")
		}
		logging.WithError(err, log).Error("get cico partner list")
		return nil, status.Error(codes.Internal, "failed to get cico partner list")
	}
	var pnr []*ppb.CICOPartnerList
	for _, v := range res {
		pnr = append(pnr, &ppb.CICOPartnerList{
			ID:      v.ID,
			Stype:   v.Stype,
			Name:    v.Name,
			Created: timestamppb.New(v.Created),
			Updated: timestamppb.New(v.Updated),
			Status:  v.Status,
		})
	}
	return &ppb.GetCICOPartnerListResponse{
		CICOPartnerList: pnr,
	}, nil
}
