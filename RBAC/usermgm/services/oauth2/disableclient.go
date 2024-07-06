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
	"brank.as/rbac/usermgm/core"

	oapb "brank.as/rbac/gunk/v1/oauth2"
)

func (s *Svc) DisableClient(ctx context.Context, req *oapb.DisableClientRequest) (*oapb.DisableClientResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.oauth2.updateclient")

	log.WithField("request", req).Debug("received")
	if err := validation.ValidateStruct(req,
		validation.Field(&req.ClientID, validation.Required, is.Alphanumeric),
	); err != nil {
	}

	if _, err := s.cl.DeleteClient(ctx, core.AuthCodeClient{
		OrgID:          hydra.OrgID(ctx),
		ClientID:       req.ClientID,
		ClientName:     "",
		ClientSecret:   "",
		CORS:           []string{},
		RedirectURLs:   []string{},
		LogoutRedirect: "",
		Logo:           "",
		GrantTypes:     []string{},
		ResponseTypes:  []string{},
		Scopes:         []string{},
		Audience:       []string{},
		SubjectType:    "",
		AuthMethod:     "",
		AuthBackend:    "",
		CreatedBy:      "",
		UpdatedBy:      "",
		DeletedBy:      hydra.ClientID(ctx),
		Created:        time.Time{},
		Updated:        time.Time{},
		Deleted:        time.Time{},
	}); err != nil {
		logging.WithError(err, log).Error("client delete")
		switch status.Code(err) {
		case codes.InvalidArgument, codes.NotFound:
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to disable client")
	}

	return &oapb.DisableClientResponse{}, nil
}
