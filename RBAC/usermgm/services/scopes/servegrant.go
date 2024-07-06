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

func (s *Svc) ServeGrant(ctx context.Context, req *cpb.ServeGrantRequest) (*cpb.ServeGrantResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.scopes.servegrant")
	log.WithField("req", req).Debug("received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, is.UUIDv4),
		validation.Field(&req.ClientID, validation.Required),
		validation.Field(&req.OwnerID, is.UUIDv4),
		validation.Field(&req.Requested, validation.Each(validation.Required, is.ASCII)),
		validation.Field(&req.Granted, validation.Each(validation.Required, is.ASCII)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	sc, err := s.gr.GetGrant(ctx, core.ConsentGrant{
		UserID:   req.UserID,
		ClientID: req.ClientID,
		OwnerID:  req.OwnerID,
		Scopes:   append(req.Requested, req.Granted...),
	})
	if err != nil {
		logging.WithError(err, log).Error("fetch scopes")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to record scope")
	}
	sd := map[string]*cpb.ScopeDetail{}
	gr := map[string]*cpb.GroupDetail{}
	for _, v := range sc.Scopes {
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

	nw := map[string]*cpb.ScopeDetail{}
	prv := map[string]*cpb.ScopeDetail{}
	for _, g := range req.Granted {
		prv[g] = sd[g]
	}
	for _, g := range req.Requested {
		nw[g] = sd[g]
	}

	return &cpb.ServeGrantResponse{
		OrgID:         sc.OrgID,
		OrgName:       sc.OrgName,
		NewScopes:     nw,
		GrantedScopes: prv,
		Groups:        gr,
		Skip:          sc.Skip,
	}, nil
}
