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

func (s *Svc) RejectServiceRequest(ctx context.Context, req *spb.ServiceStatusRequestRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)

	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.UpdatedBy, required),
		validation.Field(&req.Partner, required),
		validation.Field(&req.SvcName, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.st.RejectSvcRequest(ctx, storage.ServiceRequest{
		OrgID:     req.GetOrgID(),
		Partner:   req.GetPartner(),
		SvcName:   req.GetSvcName(),
		UpdatedBy: req.GetUpdatedBy(),
	}); err != nil {
		logging.WithError(err, log).Error("reject request failed")
		return nil, err
	}

	go func() {
		err := s.sendServiceNotification(ctx, req.GetOrgID(), req.GetSvcName())
		if err != nil {
			logging.WithError(err, log).Error("send service notification failed")
		}
	}()

	return &emptypb.Empty{}, nil
}
