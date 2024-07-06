package scopes

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"
)

type Svc struct {
	st   *postgres.Storage
	bsID BootstrapID
}

type BootstrapID interface {
	Org() string
}

func New(st *postgres.Storage, bootstrapOrg BootstrapID) *Svc {
	return &Svc{st: st, bsID: bootstrapOrg}
}

func (s *Svc) UpsertScope(ctx context.Context, sc core.Scope) (*core.Scope, error) {
	sp, err := s.st.UpsertScope(ctx, storage.Scope(sc))
	if err != nil {
		return nil, err
	}
	return (*core.Scope)(sp), nil
}

func (s *Svc) UpdateGroup(ctx context.Context, sc core.ScopeGroup) (*core.ScopeGroup, error) {
	sp, err := s.st.UpdateGroup(ctx, storage.ScopeGroup{
		Name: sc.Name,
		Desc: sc.Desc,
	})
	if err != nil {
		return nil, err
	}
	sc.Updated = sp.Updated
	return &sc, nil
}

func (s *Svc) GetScopes(ctx context.Context, sc []string) (map[string]core.ScopeGroup, error) {
	log := logging.FromContext(ctx).WithField("method", "core.scopes.getscopes")
	if len(sc) == 0 {
		return nil, nil
	}

	l, err := s.st.GetScopes(ctx, sc)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "scopes not registered")
		}
		return nil, err
	}
	log.WithField("requested", sc).WithField("scopes", l).Debug("scopes fetched")
	m := map[string]bool{}
	gn := []string{}
	for _, p := range l {
		if m[p.Group] {
			continue
		}
		gn = append(gn, p.Group)
		m[p.Group] = true
	}
	gr, err := s.st.GetScopeGroups(ctx, gn)
	if err != nil {
		return nil, err
	}
	sg := map[string]core.ScopeGroup{}
	for _, g := range gr {
		sg[g.Name] = core.ScopeGroup{
			Name:    g.Name,
			Desc:    g.Desc,
			Scopes:  []core.Scope{},
			Updated: g.Updated,
		}
	}
	for _, p := range l {
		scp := p
		grp := sg[p.Group]
		grp.Scopes = append(grp.Scopes, core.Scope(scp))
		sg[p.Group] = grp
	}
	return sg, nil
}

func (s *Svc) GetGrant(ctx context.Context, sc core.ConsentGrant) (*core.OfferGrant, error) {
	log := logging.FromContext(ctx).WithField("method", "core.scopes.getgrant")
	gr, err := s.GetScopes(ctx, sc.Scopes)
	if err != nil {
		return nil, err
	}
	og := &core.OfferGrant{OrgID: sc.OwnerID, Scopes: gr}
	o, err := s.st.GetOrgByID(ctx, sc.OwnerID)
	if err != nil {
		logging.WithError(err, log).Error("get org")
	} else {
		og.OrgName = o.OrgName
	}
	if sc.UserID != "" {
		u, err := s.st.GetUserByID(ctx, sc.UserID)
		if err != nil {
			return nil, err
		}
		og.Skip = u.OrgID == sc.OwnerID || sc.OwnerID == s.bsID.Org()
		return og, nil
	}
	og.Skip = true
	return og, nil
}

func (s *Svc) RecordGrant(ctx context.Context, sc core.ConsentGrant) (*core.ConsentGrant, error) {
	g, err := s.st.RecordGrant(ctx, storage.ConsentGrant(sc))
	if err != nil {
		return nil, err
	}
	return (*core.ConsentGrant)(g), nil
}
