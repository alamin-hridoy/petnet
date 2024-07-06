package service

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) AddUploadSvcRequest(ctx context.Context, req *spb.AddUploadSvcRequestRequest) (*spb.AddUploadSvcRequestResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Partner, required),
		validation.Field(&req.SvcName, required),
		validation.Field(&req.FileType, required),
		validation.Field(&req.FileID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	lRes, err := s.st.CreateUploadSvcRequest(ctx, storage.UploadServiceRequest{
		OrgID:    req.GetOrgID(),
		Partner:  req.GetPartner(),
		SvcName:  req.GetSvcName(),
		FileType: req.GetFileType(),
		FileID:   req.GetFileID(),
		CreateBy: req.GetCreateBy(),
	})
	if err != nil {
		if err == storage.Conflict {
			lRes, err = s.st.UpdateUploadSvcRequest(ctx, storage.UploadServiceRequest{
				OrgID:    req.GetOrgID(),
				Partner:  req.GetPartner(),
				SvcName:  req.GetSvcName(),
				FileType: req.GetFileType(),
				FileID:   req.GetFileID(),
				CreateBy: req.GetCreateBy(),
				VerifyBy: req.GetCreateBy(),
			})
			if err != nil {
				logging.WithError(err, log).Error("update Upload Svc Request failed")
				return nil, err
			}
		}
	}
	return &spb.AddUploadSvcRequestResponse{
		ID:       lRes.ID,
		OrgID:    lRes.OrgID,
		Partner:  lRes.Partner,
		SvcName:  lRes.SvcName,
		Status:   lRes.Status,
		FileType: lRes.FileType,
		FileID:   lRes.FileID,
		CreateBy: lRes.CreateBy,
		VerifyBy: lRes.VerifyBy,
		Total:    lRes.Total,
		Created:  timestamppb.New(lRes.Created),
		Verified: timestamppb.New(lRes.Verified.Time),
	}, nil
}
