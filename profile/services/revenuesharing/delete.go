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

func (s *Svc) DeleteRevenueSharing(ctx context.Context, req *rc.DeleteRevenueSharingRequest) (*emptypb.Empty, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.UserID, required, is.UUID),
		validation.Field(&req.Partner, required),
		validation.Field(&req.TransactionType, required),
		validation.Field(&req.RemitType, required),
		validation.Field(&req.BoundType, required),
	); err != nil {
		logging.WithError(err, log).Error("delete revenue sharing validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	if err := s.core.DeleteRevenueSharing(ctx, req); err != nil {
		return nil, err
	}
	return &emptypb.Empty{}, nil
}
