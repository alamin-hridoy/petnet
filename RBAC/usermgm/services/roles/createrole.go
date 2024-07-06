package role

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) CreateRole(ctx context.Context, req *ppb.CreateRoleRequest) (*ppb.CreateRoleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.role.createrole")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required, validation.Length(1, 0)),
		validation.Field(&req.Description),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := s.perm.CreateRole(ctx, core.Role{
		OrgID:      hydra.OrgID(ctx),
		Name:       req.Name,
		Desc:       req.Description,
		CreateUID:  hydra.ClientID(ctx),
		UpdatedUID: hydra.ClientID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("processing")
		return nil, err
	}

	return &ppb.CreateRoleResponse{ID: id}, nil
}
