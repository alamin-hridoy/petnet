package validation

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) ValidatePermission(ctx context.Context, req *ppb.ValidatePermissionRequest) (*ppb.ValidatePermissionResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.permission.validatepermission")
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Action, validation.Required),
		validation.Field(&req.Resource, validation.Required),
		validation.Field(&req.ID, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	idn, err := s.v.Validate(ctx, core.Validation{
		Environment: req.Environment,
		Action:      req.Action,
		Resource:    req.Resource,
		ID:          req.ID,
	})
	if err != nil {
		logging.WithError(err, log).Error("request not authorized")
		return nil, status.Error(codes.PermissionDenied, "request not authorized")
	}
	return &ppb.ValidatePermissionResponse{
		ID:        idn.ID,
		OrgID:     idn.OrgID,
		Validated: tspb.Now(),
	}, nil
}
