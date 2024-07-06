package oauthclient

import (
	"context"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) UpdateClient(ctx context.Context, cl core.AuthCodeClient) (*core.AuthCodeClient, error) {
	log := logging.FromContext(ctx).WithField("method", "core.oauthclient.updateclient")
	c, err := s.st.GetOauthClientByID(ctx, cl.ClientID)
	if err != nil {
		if err == storage.NotFound {
			return nil, status.Error(codes.NotFound, "client not found")
		}
		return nil, err
	}

	if cl.OrgID != c.OrgID {
		log.WithField("request_org", cl.OrgID).
			WithField("target_org", c.OrgID).
			Error("org mismatch")
		return nil, status.Error(codes.NotFound, "client not found")
	}

	hc, err := s.hy.GetClient(ctx, cl.ClientID)
	if err != nil {
		logging.WithError(err, log).Error("hydra fetch")
		return nil, err
	}

	if hc.ClientName != cl.ClientName {
		log.WithField("request_name", cl.ClientName).
			WithField("target_name", hc.ClientName).
			Error("name mismatch")
		return nil, status.Error(codes.InvalidArgument, "client name is immutable")
	}

	hc.AuthConfig = configMerge(hc.AuthConfig, cl.AuthConfig)

	// cors
	if cl.CORS != nil {
		hc.CORS = cl.CORS
	}
	// redir
	if cl.RedirectURLs != nil {
		hc.RedirectURIs = cl.RedirectURLs
	}
	// logout
	if cl.LogoutRedirect != "" {
		hc.PostLogoutRedirectURIs = []string{cl.LogoutRedirect}
	}
	// logo
	if cl.Logo != "" {
		hc.LogoURL = cl.Logo
	}
	// scopes
	if cl.Scopes != nil {
		hc.Scopes = cl.Scopes
	}

	if err := s.hy.UpdateClient(ctx, *hc); err != nil {
		return nil, err
	}

	c.UpdateUserID = cl.UpdatedBy
	stcl, err := s.st.UpdateOauthClient(ctx, *c)
	if err != nil {
		logging.WithError(err, log).Error("storage update")
		return nil, err
	}

	lgo := ""
	if len(hc.PostLogoutRedirectURIs) != 0 {
		lgo = hc.PostLogoutRedirectURIs[0]
	}

	return &core.AuthCodeClient{
		OrgID:          stcl.OrgID,
		ClientID:       stcl.ClientID,
		ClientName:     hc.ClientName,
		CORS:           hc.CORS,
		RedirectURLs:   hc.RedirectURIs,
		LogoutRedirect: lgo,
		Logo:           hc.LogoURL,
		GrantTypes:     hc.GrantTypes,
		ResponseTypes:  hc.ResponseTypes,
		Scopes:         hc.Scopes,
		Audience:       hc.Audience,
		SubjectType:    hc.SubjectType,
		AuthMethod:     hc.AuthMethod,
		AuthBackend:    hc.AuthBackend,
		CreatedBy:      stcl.CreateUserID,
		UpdatedBy:      stcl.UpdateUserID,
		Created:        stcl.Created,
		Updated:        stcl.Updated,
	}, nil
}

func configMerge(dst client.AuthConfig, src core.AuthClientConfig) client.AuthConfig {
	m := func(d, n string) string {
		if n != "" {
			return n
		}
		return d
	}
	d := func(d, n time.Duration) time.Duration {
		if n != 0 {
			return n
		}
		return d
	}
	return client.AuthConfig{
		LoginTmpl:       m(dst.LoginTmpl, src.LoginTmpl),
		OTPTmpl:         m(dst.LoginTmpl, src.LoginTmpl),
		ConsentTmpl:     m(dst.LoginTmpl, src.LoginTmpl),
		RememberConsent: dst.RememberConsent,
		SessionDuration: d(dst.SessionDuration, src.SessionDuration),
		Authenticator:   m(dst.LoginTmpl, src.LoginTmpl),
	}
}
