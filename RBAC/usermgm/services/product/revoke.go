package product

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/rbac/gunk/v1/permissions"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
)

func (s *Svc) RevokeService(ctx context.Context, req *ppb.RevokeServiceRequest) (*ppb.RevokeServiceResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.product.revokeservice")

	reqUUID := []validation.Rule{validation.Required, is.UUIDv4}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, reqUUID...),
		validation.Field(&req.ServiceID, reqUUID...),
		validation.Field(&req.Environment, validation.Required, validation.In(s.env...)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.pm.RevokeService(ctx, core.Grant{
		RoleID:      req.OrgID,
		GrantID:     req.ServiceID,
		Environment: req.Environment,
	}); err != nil {
		logging.WithError(err, log).Error("core revoke")
		return nil, err
	}

	return &ppb.RevokeServiceResponse{Revoked: tspb.Now()}, nil
}
