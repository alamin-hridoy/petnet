package revenuesharingreport

import (
	"context"

	rsr "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) GetRevenueSharingReportList(ctx context.Context, req *rsr.GetRevenueSharingReportListRequest) (*rsr.GetRevenueSharingReportListResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("get revenue sharing validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.GetRevenueSharingReportList(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
