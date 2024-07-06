package revenuesharing

import (
	"context"

	rc "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

func (s *Svc) GetRevenueSharingList(ctx context.Context, req *rc.GetRevenueSharingListRequest) (*rc.GetRevenueSharingListResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.RemitType, required),
	); err != nil {
		logging.WithError(err, log).Error("get revenue validation error")
		return nil, err
	}
	res, err := s.core.GetRevenueSharingList(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
