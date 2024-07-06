package service

import (
	"context"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListUploadSvcRequest(ctx context.Context, req *spb.ListUploadSvcRequestRequest) (*spb.ListUploadSvcRequestResponse, error) {
	log := logging.FromContext(ctx)
	rs, err := s.st.ListUploadSvcRequest(ctx, storage.UploadSvcRequestFilter{
		OrgID:   req.GetOrgID(),
		Status:  req.GetStatus(),
		SvcName: req.GetSvcNames(),
		Partner: req.GetPartners(),
	})
	if err != nil {
		logging.WithError(err, log).Error("List Upload Svc Request failed")
		return nil, err
	}
	us := []*spb.UploadSvcResponse{}
	for _, r := range rs {
		us = append(us, &spb.UploadSvcResponse{
			ID:       r.ID,
			OrgID:    r.OrgID,
			Partner:  r.Partner,
			SvcName:  r.SvcName,
			Status:   r.Status,
			FileType: r.FileType,
			FileID:   r.FileID,
			CreateBy: r.CreateBy,
			VerifyBy: r.VerifyBy,
			Total:    r.Total,
			Created:  timestamppb.New(r.Created),
			Verified: timestamppb.New(r.Verified.Time),
		})
	}
	return &spb.ListUploadSvcRequestResponse{
		Results: us,
	}, nil
}
