package scopes

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	cpb "brank.as/rbac/gunk/v1/consent"
)

func (s *Svc) UpdateGroup(ctx context.Context, req *cpb.UpdateGroupRequest) (*cpb.UpdateGroupResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.scopes.updategroup")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.GroupName, validation.Required, is.ASCII),
		validation.Field(&req.Description, validation.Required, is.ASCII),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sc, err := s.sc.UpdateGroup(ctx, core.ScopeGroup{
		Name: req.GroupName,
		Desc: req.Description,
	})
	if err != nil {
		logging.WithError(err, log).Error("storage date")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update group")
	}

	return &cpb.UpdateGroupResponse{Updated: tspb.New(sc.Updated)}, nil
}
