package role

import (
	"context"
	"errors"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) ListUserRoles(ctx context.Context, req *ppb.ListUserRolesRequest) (*ppb.ListUserRolesResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.roles.ListUserRoles")
	roles, err := s.perm.ListUserRoles(ctx, core.ListUserRolesRequest{
		UserID: req.GetUserID(),
	})
	if err != nil {
		logging.WithError(err, log).Error("process failed")
		return nil, err
	}
	return &ppb.ListUserRolesResponse{
		Roles: roles,
	}, nil
}

func (s *Svc) ListRole(ctx context.Context, req *ppb.ListRoleRequest) (*ppb.ListRoleResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.roles.list")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, is.UUIDv4),
		validation.Field(&req.UserID, is.UUIDv4),
		validation.Field(&req.ID,
			validation.When(len(req.GetID()) > 0, validation.Each(is.UUIDv4)),
		),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orgID := hydra.OrgID(ctx)
	if req.GetOrgID() != orgID && req.GetOrgID() != "" {
		err := errors.New("invalid OrgID")
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sortBy := "ASC"
	if req.GetSortBy() == ppb.SortBy_DESC {
		sortBy = "DESC"
	}

	roles, err := s.perm.ListRole(ctx, core.ListRoleFilter{
		ID:     req.GetID(),
		OrgID:  orgID,
		SortBy: sortBy,
		UserID: req.GetUserID(),
		Name:   req.GetName(),
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
	})
	if err != nil {
		logging.WithError(err, log).Error("process failed")
		return nil, err
	}

	ts := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}
	var ppbRoles []*ppb.Role
	var total int32
	for _, role := range roles {
		ppbRoles = append(ppbRoles, &ppb.Role{
			ID:          role.ID,
			OrgID:       role.OrgID,
			Name:        role.Name,
			Description: role.Desc,
			Members:     role.Members,
			Permissions: role.Permissions,
			CreateUID:   role.CreateUID,
			DeleteUID:   role.DeleteUID,
			UpdatedUID:  role.UpdatedUID,
			Created:     ts(role.Created),
			Updated:     ts(role.Updated),
		})
	}
	if len(roles) > 0 {
		total = int32(roles[0].Count)
	}
	return &ppb.ListRoleResponse{
		Roles: ppbRoles,
		Total: total,
	}, nil
}
