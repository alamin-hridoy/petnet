package partnerlist

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) CreatePartnerList(ctx context.Context, req *spb.CreatePartnerListRequest) (*spb.CreatePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.CreatePartnerList(ctx, &storage.PartnerList{
		Stype:            req.GetPartnerList().GetStype(),
		Name:             req.GetPartnerList().GetName(),
		Created:          req.GetPartnerList().GetCreated().AsTime(),
		Updated:          req.GetPartnerList().GetUpdated().AsTime(),
		Status:           req.GetPartnerList().GetStatus(),
		ServiceName:      req.GetPartnerList().GetServiceName(),
		Platform:         req.GetPartnerList().GetPlatform(),
		IsProvider:       req.GetPartnerList().GetIsProvider(),
		UpdatedBy:        req.GetPartnerList().GetUpdatedBy(),
		PerahubPartnerID: req.GetPartnerList().GetPerahubPartnerID(),
		RemcoID:          req.GetPartnerList().GetRemcoID(),
	})
	if err != nil {
		if err == storage.Conflict {
			logging.WithError(err, log).Error("partner already exists")
			return nil, status.Error(codes.AlreadyExists, "failed to Create Partner List")
		}
		logging.WithError(err, log).Error("store partner")
		return nil, status.Error(codes.Internal, "failed to Create Partner List")
	}

	return &spb.CreatePartnerListResponse{
		PartnerList: &spb.PartnerList{
			ID:               res.ID,
			Stype:            res.Stype,
			Name:             res.Name,
			Created:          timestamppb.New(res.Created),
			Updated:          timestamppb.New(res.Updated),
			Status:           res.Status,
			ServiceName:      res.ServiceName,
			Platform:         res.Platform,
			IsProvider:       res.IsProvider,
			PerahubPartnerID: res.PerahubPartnerID,
			RemcoID:          res.RemcoID,
		},
	}, nil
}
