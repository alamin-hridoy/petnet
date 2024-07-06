package service

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

func (s *Svc) AcceptUploadSvcRequest(ctx context.Context, req *spb.AcceptUploadSvcRequestRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Partner, required),
		validation.Field(&req.SvcName, required),
		validation.Field(&req.FileType, required),
		validation.Field(&req.VerifyBy, required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.st.AcceptUploadSvcRequest(ctx, storage.UploadServiceRequest{
		OrgID:    req.GetOrgID(),
		Partner:  req.GetPartner(),
		SvcName:  req.GetSvcName(),
		FileType: req.GetFileType(),
		VerifyBy: req.GetVerifyBy(),
	}); err != nil {
		logging.WithError(err, log).Error("accept request failed")
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
