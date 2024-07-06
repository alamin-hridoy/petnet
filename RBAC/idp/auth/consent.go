package auth

import (
	"context"

	"google.golang.org/grpc"

	"brank.as/rbac/serviceutil/logging"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	cpb "brank.as/rbac/gunk/v1/consent"
)

type ConsentGrantor interface {
	ServeGrant(ctx context.Context, g Grant) (*GrantDetail, error)
	Grant(ctx context.Context, g Grant) (*Grant, error)
}

func NewConsentGrantor(conn *grpc.ClientConn) ConsentGrantor {
	if conn == nil {
		return nil
	}
	return &ConsentSvc{cl: cpb.NewGrantServiceClient(conn)}
}

type ConsentSvc struct {
	cl cpb.GrantServiceClient
}

type Grant struct {
	Challenge string
	ID        string
	UserID    string
	ClientID  string
	OwnerID   string
	Requested []string
	Granted   []string
	Remember  bool
}

type GrantDetail struct {
	OrgID     string
	OrgName   string
	Skip      bool
	Requested map[string]Scope
	Granted   map[string]Scope
	Groups    map[string]Group
}

type Scope struct {
	ID          string
	Name        string
	Group       string
	Description string
}

type Group struct {
	Name        string
	Description string
}

// Serve grant request for user consent grant.
func (s *ConsentSvc) ServeGrant(ctx context.Context, g Grant) (*GrantDetail, error) {
	log := logging.FromContext(ctx).WithField("method", "auth.consent.servegrant")
	log.WithField("req", g).Trace("grant request")
	r, err := s.cl.ServeGrant(ctx, &cpb.ServeGrantRequest{
		UserID:    g.UserID,
		ClientID:  g.ClientID,
		OwnerID:   g.OwnerID,
		Requested: g.Requested,
		Granted:   g.Granted,
	})
	if err != nil {
		return nil, err
	}
	return &GrantDetail{
		OrgID:     r.GetOrgID(),
		OrgName:   r.GetOrgName(),
		Skip:      r.GetSkip(),
		Requested: sc(r.GetNewScopes()),
		Granted:   sc(r.GetGrantedScopes()),
		Groups:    gr(r.GetGroups()),
	}, nil
}

func sc(sc map[string]*cpb.ScopeDetail) map[string]Scope {
	if sc == nil {
		return nil
	}
	m := make(map[string]Scope, len(sc))
	for k, v := range sc {
		m[k] = Scope{
			ID:          v.Scope,
			Name:        v.Name,
			Group:       v.Group,
			Description: v.Description,
		}
	}
	return m
}

func gr(g map[string]*cpb.GroupDetail) map[string]Group {
	if g == nil {
		return nil
	}
	m := make(map[string]Group, len(g))
	for k, v := range g {
		m[k] = Group{Name: v.Name, Description: v.Description}
	}
	return m
}

func (s *ConsentSvc) Grant(ctx context.Context, g Grant) (*Grant, error) {
	log := logging.FromContext(ctx).WithField("method", "auth.consent.grant")
	log.WithField("req", g).Trace("consent grant")
	r, err := s.cl.Grant(ctx, &cpb.GrantRequest{
		UserID:    g.UserID,
		ClientID:  g.ClientID,
		OwnerID:   g.OwnerID,
		Scopes:    g.Granted,
		Timestamp: tspb.Now(),
	})
	if err != nil {
		return nil, err
	}
	g.ID = r.GetGrantID()
	g.Granted = r.GetGrants()
	return &g, nil
}
