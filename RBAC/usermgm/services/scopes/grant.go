package scopes

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	cpb "brank.as/rbac/gunk/v1/consent"
)

func (s *Svc) Grant(ctx context.Context, req *cpb.GrantRequest) (*cpb.GrantResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.scopes.grant")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.ClientID, validation.Required),
		validation.Field(&req.OwnerID, validation.Required, is.UUIDv4),
		validation.Field(&req.Scopes, validation.Each(validation.Required, is.ASCII)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sc, err := s.gr.RecordGrant(ctx, core.ConsentGrant{
		UserID:   req.UserID,
		ClientID: req.ClientID,
		OwnerID:  req.OwnerID,
		Scopes:   req.Scopes,
	})
	if err != nil {
		logging.WithError(err, log).Error("fetch scopes")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to record scope")
	}
	return &cpb.GrantResponse{GrantID: sc.ID, Grants: sc.Scopes}, nil
}
