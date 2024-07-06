package permissions

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/go-ozzo/ozzo-validation/v4/is"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) ListPermission(ctx context.Context, req *ppb.ListPermissionRequest) (*ppb.ListPermissionResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.permissions.list")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, is.UUIDv4),
		validation.Field(&req.ID,
			validation.When(len(req.GetID()) > 0, validation.Each(is.UUID)),
		),
		validation.Field(&req.Environment, is.Alphanumeric),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	orgID := hydra.OrgID(ctx)
	if req.GetOrgID() != orgID && req.GetOrgID() != "" {
		// Check admin/bootstrap
		r, a, _ := s.Permission(ctx, "ListPermission")
		if _, err := s.val.Validate(ctx, core.Validation{
			Environment: req.GetEnvironment(),
			Action:      a,
			Resource:    r,
			ID:          hydra.ClientID(ctx),
		}); err != nil {
			err := errors.New("invalid OrgID")
			logging.WithError(err, log).Error("invalid request")
			return nil, status.Error(codes.InvalidArgument, err.Error())
		}
	}

	perms, err := s.perm.ListPermission(ctx, core.ListPermissionFilter{
		ID:          req.GetID(),
		OrgID:       orgID,
		Environment: req.GetEnvironment(),
	})
	if err != nil {
		logging.WithError(err, log).Error("process failed")
		return nil, err
	}

	var ps []*ppb.Permission
	for _, p := range perms {
		ps = append(ps, &ppb.Permission{
			ID:          p.ID,
			ServiceName: p.SvcName,
			Name:        p.Name,
			Description: p.Description,
			Environment: p.Environment,
			Restrict:    !p.Allow,
			Action:      p.Action,
			Resource:    p.Resource,
			Groups:      p.Groups,
		})
	}

	return &ppb.ListPermissionResponse{
		Permissions: ps,
	}, nil
}
