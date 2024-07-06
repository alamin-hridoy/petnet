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

func (s *Svc) ApplyServiceRequest(ctx context.Context, req *spb.ApplyServiceRequestRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Type, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.st.ApplySvcRequest(ctx, storage.ServiceRequest{
		OrgID:   req.OrgID,
		SvcName: req.Type.String(),
	}); err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, err
	}
	return new(emptypb.Empty), nil
}
