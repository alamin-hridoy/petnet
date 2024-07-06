package permissions

import (
	"context"
	"fmt"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw/hydra"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) CreatePermission(ctx context.Context, p core.ServicePermission) (*core.ServicePermission, error) {
	log := logging.FromContext(ctx).WithField("method", "core.permission.createpermission")

	svc, err := s.store.UpsertService(ctx, storage.Service{
		Name:        p.Service.Name,
		Description: p.Service.Description,
	})
	if err != nil {
		return nil, err
	}

	sp := []storage.ServicePermission{}
	for i, r := range p.Res {
		for _, a := range r.Actions {
			ps, err := s.store.UpsertServicePermission(ctx, storage.ServicePermission{
				ServiceID:   svc.ID,
				Name:        r.Name,
				Description: r.Description,
				Resource:    r.Resource,
				Action:      a,
			})
			if err != nil {
				logging.WithError(err, log).Error("store permission")
				continue
			}
			p.Res[i].ID = ps.ID
			if ps.Created.Equal(ps.Updated) {
				sp = append(sp, *ps)
			}
		}
	}

	p.Service.ID = svc.ID
	if !svc.Default {
		return &p, nil
	}
	gr, err := s.store.ListServiceAssignID(ctx, svc.ID)
	if err != nil {
		logging.WithError(err, log).Error("list grants")
		return &p, nil
	}

	for _, prm := range sp {
		for _, g := range gr {
			_, err := s.CreateOrgPermission(ctx, core.OrgPermission{
				OrgID:       g.OrgID,
				SvcID:       g.ServiceID,
				SvcPermID:   prm.ID,
				GrantID:     g.GrantID,
				CreateUID:   hydra.ClientID(ctx),
				Name:        prm.Name,
				Description: prm.Description,
				Environment: g.Environment,
				Allow:       true,
				Action:      prm.Action,
				Resource:    prm.Resource,
			})
			if err != nil {
				logging.WithError(err, log).Error("org permission")
			}
		}
	}
	return &p, nil
}

func (s *Svc) CreateOrgPermission(ctx context.Context, p core.OrgPermission) (string, error) {
	pm := keto.Permission{
		Description: p.Description,
		Environment: p.Environment,
		Allow:       p.Allow,
		Actions:     []string{p.Action},
		Resources:   []string{fmt.Sprintf("org:%s:%s", p.OrgID, p.Resource)},
		Groups:      p.Groups,
	}

	id, err := s.keto.CreatePermission(ctx, pm)
	if err != nil {
		return "", err
	}

	if _, err := s.store.UpsertOrgPermission(ctx, storage.OrgPermission{
		ID:           id,
		OrgID:        p.OrgID,
		GrantID:      p.GrantID,
		ServiceID:    p.SvcID,
		PermissionID: p.SvcPermID,
		Name:         p.Name,
		Description:  p.Description,
		Resource:     p.Resource,
		Action:       p.Action,
		Environment:  p.Environment,
	}); err != nil {
		return "", err
	}

	return id, nil
}
