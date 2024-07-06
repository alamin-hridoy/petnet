package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Svc) DeleteRevenueSharingTier(ctx context.Context, req *rc.DeleteRevenueSharingTierRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RevenueSharingID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("delete revenue tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.DeleteRevenueSharingTier(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}

func (s *Svc) DeleteRevenueSharingTierById(ctx context.Context, req *rc.DeleteRevenueSharingTierByIdRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("delete revenue sharing tier validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.core.DeleteRevenueSharingTierById(ctx, req); err != nil {
		return nil, err
	}
	return nil, nil
}
