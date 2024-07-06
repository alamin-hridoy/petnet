package revenuesharingreport

import (
	"context"

	rsp "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CreateRevenueSharingReport(ctx context.Context, req *rsp.CreateRevenueSharingReportRequest) (*rsp.CreateRevenueSharingReportResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
	); err != nil {
		logging.WithError(err, log).Error("create revenue sharing report validation error")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.core.CreateRevenueSharingReport(ctx, req)
	if err != nil {
		return nil, err
	}
	return res, nil
}
