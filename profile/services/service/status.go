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

func (s *Svc) SetStatusUploadSvcRequest(ctx context.Context, req *spb.SetStatusUploadSvcRequestRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Partners, required),
		validation.Field(&req.SvcName, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	for _, vv := range req.GetPartners() {
		if err := s.st.SetStatusSvcRequest(ctx, storage.ServiceRequest{
			OrgID:   req.GetOrgID(),
			Partner: vv,
			SvcName: req.GetSvcName(),
			Status:  req.GetStatus().String(),
		}); err != nil {
			if err != storage.NotFound {
				logging.WithError(err, log).Error("set status request failed")
				return nil, err
			}
		}
	}
	return &emptypb.Empty{}, nil
}
