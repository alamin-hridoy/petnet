package oauthclient

import (
	"context"

	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

// CreateClient ...
func (s *Svc) CreateClient(ctx context.Context, cl core.AuthCodeClient) (*core.AuthCodeClient, error) {
	log := logging.FromContext(ctx).WithField("method", "core.oauthclient.createclient")
	id, err := random.String(16)
	if err != nil {
		return nil, err
	}

	bufSec, err := random.String(48)
	if err != nil {
		logging.WithError(err, log).Error("random secret")
		return nil, err
	}
	if cl.AuthBackend != "" && cl.AuthConfig.IdentitySource == "" {
		cl.AuthConfig.IdentitySource = cl.AuthBackend
	}

	gr := []string{"authorization_code"}
	rsp := []string{"code"}
	for _, s := range cl.Scopes {
		switch s {
		case "offline", "offline_access":
			gr = append(gr, "refresh_token") // refresh token is necessary for offline access scope
		case "openid":
			rsp = append(rsp, "id_token") // openid connect token
		}
	}

	c, err := s.hy.CreateClient(ctx, client.AuthClient{
		OwnerID:                cl.OrgID,
		ClientID:               id,
		ClientName:             cl.ClientName,
		RedirectURIs:           cl.RedirectURLs,
		CORS:                   cl.CORS,
		PostLogoutRedirectURIs: []string{cl.LogoutRedirect},
		LogoURL:                cl.Logo,
		GrantTypes:             gr,
		ResponseTypes:          rsp,
		Scopes:                 cl.Scopes,
		Audience:               cl.Audience,
		Secret:                 bufSec,
		SubjectType:            "public",
		AuthMethod:             cl.AuthMethod,
		AuthConfig: client.AuthConfig{
			LoginTmpl:       cl.AuthConfig.LoginTmpl,
			OTPTmpl:         cl.AuthConfig.OTPTmpl,
			ConsentTmpl:     cl.AuthConfig.ConsentTmpl,
			RememberConsent: !cl.AuthConfig.ForceConsent,
			SessionDuration: cl.AuthConfig.SessionDuration,
			Authenticator:   cl.AuthConfig.IdentitySource,
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("hydra create")
		return nil, err
	}

	st, err := s.st.CreateOauthClient(ctx, storage.OAuthClient{
		OrgID:        cl.OrgID,
		ClientID:     id,
		ClientName:   cl.ClientName,
		CreateUserID: cl.CreatedBy,
		Environment:  cl.Environment,
	})
	if err != nil {
		if err := s.hy.DeleteClient(ctx, id); err != nil {
			logging.WithError(err, log).Error("hydra rollback")
		}
		return nil, err
	}
	cl.ClientID = id
	cl.ClientSecret = c.Secret
	cl.Created = st.Created
	cl.Updated = st.Updated
	return &cl, nil
}
