package permissions

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func (s *Svc) CreatePermission(ctx context.Context, req *ppb.CreatePermissionRequest) (*ppb.CreatePermissionResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.permissions.create")
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ServiceName, validation.Required),
		validation.Field(&req.Description, validation.Required),
		validation.Field(&req.Permissions, validation.Required, validation.Each(
			validation.Required, validation.By(func(i interface{}) error {
				p, ok := i.(ppb.ServicePermission)
				if !ok {
					return fmt.Errorf("invalid permission")
				}
				return validation.ValidateStruct(&p,
					validation.Field(&p.Name, validation.Required, is.ASCII),
					validation.Field(&p.Description, validation.Required, is.ASCII),
					validation.Field(&p.Resource, validation.Required, validation.By(func(interface{}) error {
						return validation.Validate(strings.Split(p.GetResource(), ":"),
							validation.Required, validation.Each(validation.Required, is.Alphanumeric))
					})),
					validation.Field(&p.Actions, validation.Required,
						validation.Each(validation.Required, is.Alphanumeric),
					),
				)
			}),
		)),
		// validation.Field(&req.)
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pm := core.ServicePermission{
		Service: core.Service{
			Name:        req.ServiceName,
			Description: req.Description,
		},
		Res: make([]core.ServiceResource, len(req.Permissions)),
	}
	for i, p := range req.Permissions {
		pm.Res[i] = core.ServiceResource{
			Name:        p.Name,
			Description: p.Description,
			Resource:    p.Resource,
			Actions:     p.Actions,
		}
	}

	prm, err := s.perm.CreatePermission(ctx, pm)
	if err != nil {
		logging.WithError(err, log).Error("failed to create permission")
		return nil, status.Error(codes.Internal, "failed to create permission")
	}
	ids := map[string]string{}
	for _, r := range prm.Res {
		ids[r.Name] = r.ID
	}

	return &ppb.CreatePermissionResponse{
		ID:        ids,
		ServiceID: prm.Service.ID,
	}, nil
}
