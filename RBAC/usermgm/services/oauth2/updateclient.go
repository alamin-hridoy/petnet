package oauth2

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	oapb "brank.as/rbac/gunk/v1/oauth2"
)

func (s *Svc) UpdateClient(ctx context.Context, req *oapb.UpdateClientRequest) (*oapb.UpdateClientResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.oauth2.updateclient")
	log.WithField("request", req).Trace("received")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.ClientID, validation.Required, is.Alphanumeric),
		validation.Field(&req.Name, is.Alphanumeric),
		validation.Field(&req.LogoURL, is.URL),
		validation.Field(&req.LogoutRedirectURL, is.URL),
		validation.Field(&req.CORS, validation.Each(validation.Required, is.URL)),
		validation.Field(&req.RedirectURL, validation.Each(validation.Required, is.URL)),
		validation.Field(&req.Scopes, validation.Each(validation.Required, is.Alphanumeric)),
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

	if _, err := s.cl.UpdateClient(ctx, core.AuthCodeClient{
		OrgID:          hydra.OrgID(ctx),
		ClientID:       req.ClientID,
		ClientName:     req.Name,
		CORS:           req.CORS,
		RedirectURLs:   req.RedirectURL,
		LogoutRedirect: req.LogoutRedirectURL,
		AuthConfig: core.AuthClientConfig{
			LoginTmpl:       req.GetConfig().GetLoginTemplate(),
			OTPTmpl:         req.GetConfig().GetLoginTemplate(),
			ConsentTmpl:     req.GetConfig().GetLoginTemplate(),
			ForceConsent:    req.GetConfig().GetForceConsent(),
			SessionDuration: req.GetConfig().GetSessionDuration().AsDuration(),
			IdentitySource:  req.GetConfig().GetLoginTemplate(),
		},
		Logo:      req.LogoURL,
		Scopes:    req.Scopes,
		UpdatedBy: hydra.ClientID(ctx),
	}); err != nil {
		logging.WithError(err, log).Error("client update")
		switch status.Code(err) {
		case codes.InvalidArgument, codes.NotFound:
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update client")
	}

	return &oapb.UpdateClientResponse{}, nil
}
