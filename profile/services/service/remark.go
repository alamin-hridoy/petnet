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

func (s *Svc) AddRemarkSvcRequest(ctx context.Context, req *spb.AddRemarkSvcRequestRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.UpdatedBy, required),
		validation.Field(&req.SvcName, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.st.AddRemarkSvcRequest(ctx, storage.ServiceRequest{
		OrgID:     req.GetOrgID(),
		UpdatedBy: req.GetUpdatedBy(),
		SvcName:   req.GetSvcName(),
		Remarks:   req.GetRemark(),
	}); err != nil {
		if err != storage.NotFound {
			logging.WithError(err, log).Error("add remark request failed")
			return nil, err
		}
	}
	return &emptypb.Empty{}, nil
}
