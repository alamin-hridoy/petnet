package role

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) UpdateRole(ctx context.Context, req *ppb.UpdateRoleRequest) (*ppb.UpdateRoleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.role.updaterole")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.ID, validation.Required, is.UUIDv4),
		validation.Field(&req.Name, validation.Required, validation.Length(1, 0)),
		validation.Field(&req.Description),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r, err := s.perm.UpdateRole(ctx, core.Role{
		ID:         req.GetID(),
		Name:       req.GetName(),
		Desc:       req.GetDescription(),
		UpdatedUID: hydra.ClientID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("processing")
		return nil, err
	}

	return &ppb.UpdateRoleResponse{
		Role: &ppb.Role{
			ID:          r.ID,
			OrgID:       r.OrgID,
			Name:        r.Name,
			Description: r.Desc,
			Members:     r.Members,
			Permissions: r.Permissions,
		},
		Updated: timestamppb.New(time.Now()),
	}, nil
}
