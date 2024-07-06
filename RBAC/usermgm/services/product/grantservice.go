package product

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	ppb "brank.as/rbac/gunk/v1/permissions"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
)

func (s *Svc) GrantService(ctx context.Context, req *ppb.GrantServiceRequest) (*ppb.GrantServiceResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.product.grantservice")

	reqUUID := []validation.Rule{validation.Required, is.UUIDv4}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, reqUUID...),
		validation.Field(&req.ServiceID, reqUUID...),
		validation.Field(&req.Environment, validation.In(s.env...)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.pm.GrantService(ctx, core.Grant{
		RoleID:      req.OrgID,
		GrantID:     req.ServiceID,
		Environment: req.Environment,
	})
	if err != nil {
		logging.WithError(err, log).Error("core grant")
		return nil, err
	}

	return &ppb.GrantServiceResponse{ID: id}, nil
}
