package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) UpdateRevenueSharing(ctx context.Context, req *rc.UpdateRevenueSharingRequest) (*rc.UpdateRevenueSharingResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, required, is.UUID),
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.UserID, required, is.UUID),
		validation.Field(&req.Partner, required),
		validation.Field(&req.BoundType, required),
		validation.Field(&req.TransactionType, required),
		validation.Field(&req.RemitType, required),
	); err != nil {
		logging.WithError(err, log).Error("update revenue sharing validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.UpdateRevenueSharing(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
