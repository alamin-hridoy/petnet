package cashincashout

import (
	"context"

	"brank.as/petnet/api/util"
	cbp "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) CICOTransactList(ctx context.Context, req *cbp.CICOTransactListRequest) (*cbp.CICOTransactListResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.cashincashout.CICOTransactList")
	r, err := s.CICOTransactListValidate(ctx, req)
	if err != nil {
		logging.WithError(err, log).Error("validate request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.core.CICOTransactList(ctx, r)
	if err != nil {
		logging.WithError(err, log).Error("failed to get CICO Transact List")
		return nil, util.HandleServiceErr(err)
	}

	return res, nil
}

func (s *Svc) CICOTransactListValidate(ctx context.Context, req *cbp.CICOTransactListRequest) (*cbp.CICOTransactListRequest, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID)); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	logging.FromContext(ctx).Info("CICO Transact List validation")
	return req, nil
}
