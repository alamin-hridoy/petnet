package partnerlist

import (
	"context"
	"strings"

	ppb "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) GetPartnerList(ctx context.Context, req *ppb.GetPartnerListRequest) (*ppb.GetPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.GetPartnerList(ctx, &storage.PartnerList{
		ID:          req.GetID(),
		Stype:       req.GetStype(),
		Status:      req.GetStatus(),
		Name:        req.GetName(),
		ServiceName: req.GetServiceName(),
		Platform:    req.GetPlatform(),
		IsProvider:  req.GetIsProvider(),
	})
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("PartnerList not found")
			return nil, status.Error(codes.NotFound, "PartnerList not found")
		}
		logging.WithError(err, log).Error("get partner list")
		return nil, status.Error(codes.Internal, "failed to get partner list")
	}
	var pnr []*ppb.PartnerList
	for _, v := range res {
		pnr = append(pnr, &ppb.PartnerList{
			ID:            v.ID,
			Stype:         v.Stype,
			Name:          v.Name,
			Created:       timestamppb.New(v.Created),
			Updated:       timestamppb.New(v.Updated),
			Status:        v.Status,
			ServiceName:   v.ServiceName,
			UpdatedBy:     v.UpdatedBy,
			DisableReason: v.DisableReason.String,
			Platform:      v.Platform,
			IsProvider:    v.IsProvider,
		})
	}
	return &ppb.GetPartnerListResponse{
		PartnerList: pnr,
	}, nil
}

func (s *Svc) GetDSAPartnerList(ctx context.Context, req *ppb.DSAPartnerListRequest) (*ppb.GetDSAPartnerListResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.GetDSAPartnerList(ctx, &storage.GetDSAPartnerListRequest{
		TransactionTypes: strings.Split(req.GetTransactionType(), ","),
	})
	if err != nil {
		if err == storage.NotFound {
			logging.WithError(err, log).Error("DSA PartnerList not found")
			return nil, status.Error(codes.NotFound, "DSA PartnerList not found")
		}
		logging.WithError(err, log).Error("get DSA partner list")
		return nil, status.Error(codes.Internal, "failed to get DSA partner list")
	}
	var pnr []*ppb.DSAPartnerList
	for _, v := range res {
		pnr = append(pnr, &ppb.DSAPartnerList{
			Partner:         v.Partner,
			TransactionType: v.TransactionType,
		})
	}
	return &ppb.GetDSAPartnerListResponse{
		DSAPartnerList: pnr,
	}, nil
}

func (s *Svc) GetPartnerByStype(ctx context.Context, req *ppb.GetPartnerByStypeRequest) (*ppb.GetPartnerByStypeResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.st.GetPartnerByStype(ctx, req.Stype)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "get partner by stype not found")
		}
		logging.WithError(err, log).Error("fetch from storage")
		return nil, status.Error(codes.Internal, "get partner by stype failed")
	}
	return &ppb.GetPartnerByStypeResponse{
		PartnerList: &ppb.PartnerList{
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
			IsProvider:       res.IsProvider,
			PerahubPartnerID: res.PerahubPartnerID,
			RemcoID:          res.RemcoID,
		},
	}, nil
}
