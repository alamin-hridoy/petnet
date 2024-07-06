package oauthclient

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

func (s *Svc) GetClient(ctx context.Context, org, clientID string, listDisable bool) ([]core.AuthCodeClient, error) {
	log := logging.FromContext(ctx).WithField("method", "core.oauthclient.getclient")

	var lst []storage.OAuthClient
	if clientID == "" {
		c, err := s.st.GetOauthClientByOrgID(ctx, org, listDisable)
		if err != nil {
			logging.WithError(err, log).Error("storage list clients")
			return nil, status.Error(codes.Internal, "failed to list clients")
		}
		lst = c
	} else {
		c, err := s.st.GetOauthClientByID(ctx, clientID)
		if err != nil {
			logging.WithError(err, log).Error("storage list clients")
			return nil, status.Error(codes.Internal, "failed to list clients")
		}
		if c.OrgID != org && org != "" {
			log.WithField("request_org", org).
				WithField("target_org", c.OrgID).
				Error("org mismatch")
			return nil, status.Error(codes.NotFound, "client not found")
		}
		lst = []storage.OAuthClient{*c}
	}

	cl := make([]core.AuthCodeClient, len(lst))
	for i, c := range lst {
		if c.Deleted.Valid {
			cl[i] = core.AuthCodeClient{
				OrgID:       c.OrgID,
				ClientID:    c.ClientID,
				Environment: c.Environment,
				CreatedBy:   c.CreateUserID,
				UpdatedBy:   c.UpdateUserID,
				DeletedBy:   c.DeleteUserID,
				Created:     c.Created,
				Deleted:     c.Deleted.Time,
				Updated:     c.Updated,
			}
			continue
		}
		hc, err := s.hy.GetClient(ctx, c.ClientID)
		if err != nil {
			logging.WithError(err, log).Error("hydra list")
			return nil, status.Error(codes.Internal, "failed to list clients")
		}
		logout := ""
		if len(hc.PostLogoutRedirectURIs) > 0 {
			logout = hc.PostLogoutRedirectURIs[0]
		}
		cl[i] = core.AuthCodeClient{
			OrgID:          c.OrgID,
			ClientID:       c.ClientID,
			ClientName:     hc.ClientName,
			Environment:    c.Environment,
			CORS:           hc.CORS,
			RedirectURLs:   hc.RedirectURIs,
			LogoutRedirect: logout,
			Logo:           hc.LogoURL,
			GrantTypes:     hc.GrantTypes,
			ResponseTypes:  hc.ResponseTypes,
			Scopes:         hc.Scopes,
			Audience:       hc.Audience,
			SubjectType:    hc.SubjectType,
			AuthMethod:     hc.AuthMethod,
			AuthBackend:    hc.AuthBackend,
			AuthConfig: core.AuthClientConfig{
				LoginTmpl:       hc.AuthConfig.LoginTmpl,
				OTPTmpl:         hc.AuthConfig.OTPTmpl,
				ConsentTmpl:     hc.AuthConfig.ConsentTmpl,
				ForceConsent:    !hc.AuthConfig.RememberConsent,
				SessionDuration: hc.AuthConfig.SessionDuration,
				IdentitySource:  hc.AuthConfig.Authenticator,
			},
			CreatedBy: c.CreateUserID,
			UpdatedBy: c.UpdateUserID,
			DeletedBy: c.DeleteUserID,
			Created:   c.Created,
			Updated:   c.Updated,
			Deleted:   c.Deleted.Time,
		}
	}
	return cl, nil
}
