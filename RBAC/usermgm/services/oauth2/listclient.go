package oauth2

import (
	"context"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	durpb "google.golang.org/protobuf/types/known/durationpb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	oapb "brank.as/rbac/gunk/v1/oauth2"
)

func (s *Svc) ListClients(ctx context.Context, req *oapb.ListClientsRequest) (*oapb.ListClientsResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.oauth2.listclients")

	log.WithField("request", req).Debug("received")
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ClientID, is.Alphanumeric),
		validation.Field(&req.OrgID, is.UUIDv4),
		validation.Field(&req.ListDisable, validation.In(true, false)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	org := req.GetOrgID()
	if org == "" {
		org = hydra.OrgID(ctx)
	}
	cl, err := s.cl.GetClient(ctx, org, req.ClientID, req.ListDisable)
	if err != nil {
		logging.WithError(err, log).Error("get clients")
		switch status.Code(err) {
		case codes.InvalidArgument, codes.NotFound:
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to list clients")
	}

	ts := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}
	e, lst := req.GetEnv(), make([]*oapb.Oauth2Client, 0, len(cl))
	for _, c := range cl {
		if e != "" && c.Environment != e {
			continue
		}
		lst = append(lst, &oapb.Oauth2Client{
			Name:              c.ClientName,
			Env:               c.Environment,
			OrgID:             c.OrgID,
			ClientID:          c.ClientID,
			LogoURL:           c.Logo,
			Scopes:            c.Scopes,
			RedirectURL:       c.RedirectURLs,
			LogoutRedirectURL: c.LogoutRedirect,
			Creator:           c.CreatedBy,
			Created:           ts(c.Created),
			Disabled:          ts(c.Deleted),
			Config: &oapb.ClientConfig{
				LoginTemplate:   c.AuthConfig.LoginTmpl,
				OTPTemplate:     c.AuthConfig.OTPTmpl,
				ConsentTemplate: c.AuthConfig.ConsentTmpl,
				ForceConsent:    c.AuthConfig.ForceConsent,
				SessionDuration: durpb.New(c.AuthConfig.SessionDuration),
				IdentitySource:  c.AuthConfig.IdentitySource,
			},
		})
	}
	return &oapb.ListClientsResponse{Clients: lst}, nil
}
