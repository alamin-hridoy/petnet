package role

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	ppb "brank.as/rbac/gunk/v1/permissions"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
)

func (s *Svc) AssignRolePermission(ctx context.Context, req *ppb.AssignRolePermissionRequest) (*ppb.AssignRolePermissionResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.roles.assignrolepermission")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.RoleID, validation.Required, is.UUIDv4),
		validation.Field(&req.Permission, validation.Required, is.UUIDv4)); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r, err := s.perm.RoleGrant(ctx, core.Grant{
		RoleID:  req.RoleID,
		GrantID: req.Permission,
	})
	if err != nil {
		logging.WithError(err, log).Error("grant role")
		if status.Code(err) != codes.Internal {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "assignment failed")
	}

	return &ppb.AssignRolePermissionResponse{
		Role: &ppb.Role{
			ID:          r.ID,
			OrgID:       r.OrgID,
			Name:        r.Name,
			Description: r.Desc,
			Members:     r.Members,
			Permissions: r.Permissions,
		},
		Updated: tspb.Now(),
	}, nil
}
