package product

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) PublicService(ctx context.Context, req *ppb.PublicServiceRequest) (*ppb.PublicServiceResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.product.grantservice")

	reqUUID := []validation.Rule{validation.Required, is.UUIDv4}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ServiceID, reqUUID...),
		validation.Field(&req.Environment, validation.In(s.env...)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if req.Enabled {
		if err := s.pm.PublicService(ctx, core.Grant{
			RoleID:      hydra.ClientID(ctx),
			GrantID:     req.ServiceID,
			Environment: req.Environment,
			Default:     req.Enabled,
		}); err != nil {
			logging.WithError(err, log).Error("core grant")
			return nil, err
		}
		return &ppb.PublicServiceResponse{Processed: tspb.Now()}, nil
	}

	if err := s.pm.PrivateService(ctx, core.Grant{
		GrantID:     req.ServiceID,
		Environment: req.Environment,
	}); err != nil {
		logging.WithError(err, log).Error("core grant")
		return nil, err
	}

	return &ppb.PublicServiceResponse{Processed: tspb.Now()}, nil
}
