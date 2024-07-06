package oauth2

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/usermgm/core"

	oapb "brank.as/rbac/gunk/v1/oauth2"
)

func (s *Svc) CreateClient(ctx context.Context, req *oapb.CreateClientRequest) (*oapb.CreateClientResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.oauth2.createclient")
	log.WithField("request", req).Trace("received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Name, validation.Required, is.ASCII),
		validation.Field(&req.Audience, is.ASCII),
		validation.Field(&req.Env, is.Alphanumeric, validation.In(s.envs...)),
		validation.Field(&req.LogoURL, is.URL),
		validation.Field(&req.LogoutRedirectURL, validation.Required, is.URL),
		validation.Field(&req.IdentitySource, is.Alphanumeric),
		validation.Field(&req.CORS, validation.Required,
			validation.Each(validation.Required, is.URL)),
		validation.Field(&req.RedirectURL, validation.Required,
			validation.Each(validation.Required, is.URL)),
		validation.Field(&req.Scopes, validation.Required,
			validation.Each(validation.Required, is.ASCII)),
		validation.Field(&req.Config, validation.By(func(interface{}) error {
			r := req.Config
			if r == nil {
				return nil
			}
			return validation.ValidateStruct(r,
				validation.Field(&r.IdentitySource, is.Alphanumeric),
			)
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	mthd := client.MethodPrivate
	if req.ClientType == oapb.ClientType_Public {
		mthd = client.MethodPublic
	}

	cl, err := s.cl.CreateClient(ctx, core.AuthCodeClient{
		OrgID:          hydra.OrgID(ctx),
		ClientName:     req.Name,
		Environment:    req.Env,
		CORS:           req.CORS,
		RedirectURLs:   req.RedirectURL,
		LogoutRedirect: req.LogoutRedirectURL,
		Logo:           req.LogoURL,
		Scopes:         req.Scopes,
		Audience:       []string{req.Audience},
		AuthMethod:     mthd,
		AuthBackend:    req.IdentitySource,
		AuthConfig: core.AuthClientConfig{
			LoginTmpl:       req.GetConfig().GetLoginTemplate(),
			OTPTmpl:         req.GetConfig().GetOTPTemplate(),
			ConsentTmpl:     req.GetConfig().GetConsentTemplate(),
			ForceConsent:    req.GetConfig().GetForceConsent(),
			SessionDuration: req.GetConfig().GetSessionDuration().AsDuration(),
			IdentitySource:  req.GetConfig().GetIdentitySource(),
		},
		CreatedBy: hydra.ClientID(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("creating client")
		if status.Code(err) == codes.InvalidArgument {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to create client")
	}

	return &oapb.CreateClientResponse{
		ClientID: cl.ClientID,
		Secret:   cl.ClientSecret,
	}, nil
}
