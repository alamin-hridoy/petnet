package partnerlist

import (
	"context"
	"database/sql"

	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) UpdatePartnerList(ctx context.Context, req *spb.UpdatePartnerListRequest) (*spb.UpdatePartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.UpdatePartnerList(ctx, &storage.PartnerList{
		ID:          req.GetPartnerList().GetID(),
		Stype:       req.GetPartnerList().GetStype(),
		Name:        req.GetPartnerList().GetName(),
		Created:     req.GetPartnerList().GetCreated().AsTime(),
		Updated:     req.GetPartnerList().GetUpdated().AsTime(),
		Status:      req.GetPartnerList().GetStatus(),
		ServiceName: req.GetPartnerList().GetServiceName(),
		UpdatedBy:   req.GetPartnerList().GetUpdatedBy(),
		Platform:    req.GetPartnerList().GetPlatform(),
		IsProvider:  req.GetPartnerList().GetIsProvider(),
		DisableReason: sql.NullString{
			String: req.GetPartnerList().GetDisableReason(),
			Valid:  true,
		},
		PerahubPartnerID: req.GetPartnerList().GetPerahubPartnerID(),
		RemcoID:          req.GetPartnerList().GetRemcoID(),
	})
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("partner Not Found")
			return nil, status.Error(codes.NotFound, "failed to update partner record")
		}
		logging.WithError(err, log).Error("store partner")
		return nil, status.Error(codes.Internal, "failed to update partner record")
	}

	return &spb.UpdatePartnerListResponse{
		PartnerList: &spb.PartnerList{
			ID:               res.ID,
			Stype:            res.Stype,
			Name:             res.Name,
			Created:          timestamppb.New(res.Created),
			Updated:          timestamppb.New(res.Updated),
			Status:           res.Status,
			ServiceName:      res.ServiceName,
			UpdatedBy:        res.UpdatedBy,
			DisableReason:    res.DisableReason.String,
			Platform:         res.Platform,
			PerahubPartnerID: res.PerahubPartnerID,
			RemcoID:          res.RemcoID,
		},
	}, nil
}
