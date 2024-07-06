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
)

func (s *Svc) UpdateServiceRequestByOrgID(ctx context.Context, req *spb.UpdateServiceRequestByOrgIDRequest) (*spb.UpdateServiceRequestByOrgIDResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OldOrgID, required, is.UUID),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.st.UpdateServiceRequestByOrgID(ctx, storage.UpdateServiceRequestOrgID{
		OldOrgID: req.GetOldOrgID(),
		NewOrgID: req.GetNewOrgID(),
		Status:   req.GetStatus(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Validate Service Access")
		return &spb.UpdateServiceRequestByOrgIDResponse{
			ID: res,
		}, nil
	}
	return &spb.UpdateServiceRequestByOrgIDResponse{
		ID: res,
	}, nil
}
