package scopes

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	cpb "brank.as/rbac/gunk/v1/consent"
)

func (s *Svc) GetScope(ctx context.Context, req *cpb.GetScopeRequest) (*cpb.GetScopeResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.scopes.upsertscope")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Scopes, validation.Required, validation.Length(1, 0),
			validation.Each(validation.Required, is.ASCII)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sc, err := s.sc.GetScopes(ctx, req.Scopes)
	if err != nil {
		logging.WithError(err, log).Error("fetch scopes")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to record scope")
	}
	sd := map[string]*cpb.ScopeDetail{}
	gr := map[string]*cpb.GroupDetail{}
	for _, v := range sc {
		if len(v.Scopes) == 0 {
			continue
		}
		for _, sp := range v.Scopes {
			sd[sp.ID] = &cpb.ScopeDetail{
				Scope:       sp.ID,
				Name:        sp.Name,
				Group:       sp.Group,
				Description: sp.Desc,
			}
		}
		gr[v.Name] = &cpb.GroupDetail{
			Name:        v.Name,
			Description: v.Desc,
		}
	}

	return &cpb.GetScopeResponse{Scopes: sd, Groups: gr}, nil
}
