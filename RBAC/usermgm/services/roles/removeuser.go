package role

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/usermgm/core"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) RemoveUser(ctx context.Context, req *ppb.RemoveUserRequest) (*ppb.RemoveUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.roles.removeuser")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.RoleID, validation.Required, is.UUIDv4),
		validation.Field(&req.UserID, validation.Required, is.UUIDv4)); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r, err := s.perm.UnassignRole(ctx, core.Grant{
		RoleID:  req.RoleID,
		GrantID: req.UserID,
	})
	if err != nil {
		logging.WithError(err, log).Error("grant role")
		if status.Code(err) != codes.Internal {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "removal failed")
	}

	return &ppb.RemoveUserResponse{
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
