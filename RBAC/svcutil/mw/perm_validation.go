package mw

import (
	"context"
	"net"
	"path"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"

	ppb "brank.as/rbac/gunk/v1/permissions"
)

func NewPermissionValidator(config *viper.Viper, env string) (*Mapper, error) {
	u := net.JoinHostPort(config.GetString("server.host"), config.GetString("server.adminPort"))
	conn, err := grpc.Dial(u, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Mapper{cl: ppb.NewValidationServiceClient(conn), env: env}, nil
}

func NewPermissionValidatorFromClient(cl ppb.ValidationServiceClient, env string) *Mapper {
	return &Mapper{cl: cl, env: env}
}

type Mapper struct {
	env string
	cl  ppb.ValidationServiceClient
}

// Service is avaliable to all authenticated users.
type PublicPermission struct{}

func (PublicPermission) Permission(context.Context, string) (string, string, bool) {
	return "", "", true
}

type Permissioner interface {
	Permission(ctx context.Context, method string) (resource, action string, public bool)
}

const (
	AuthIDKey = hydra.ClientIDKey
	OrgIDKey  = hydra.OrgIDKey
	IDNameKey = hydra.UsernameKey
	EnvKey    = "auth-env"
)

type orgGetter interface {
	GetOrgID() string
}

type userGetter interface {
	GetID() string
}

func (m *Mapper) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logging.FromContext(ctx).WithField("interceptor", "permissions")
		if p, ok := info.Server.(PublicEndpoint); ok {
			if p.PublicEndpoint(info.FullMethod) {
				return handler(ctx, req)
			}
		}
		res, act := path.Base(info.FullMethod), "call"
		if p, ok := info.Server.(Permissioner); ok {
			r, a, p := p.Permission(ctx, res)
			if p {
				// public endpoint - all authenticated users allowed.
				return handler(ctx, req)
			}
			if res == "GetUser" {
				u, exist := req.(userGetter)
				equal := u.GetID() == hydra.ClientID(ctx)
				if exist && equal {
					log.Trace("user accessing their own record")
					return handler(ctx, req)
				}
			}
			if r != "" && a != "" {
				res, act = r, a
			}
		}
		orgID := hydra.OrgID(ctx)
		if o, ok := req.(orgGetter); ok {
			if o.GetOrgID() != "" {
				orgID = o.GetOrgID()
			}
		}
		log = log.WithFields(logrus.Fields{
			"auth_org": hydra.OrgID(ctx),
			"org_id":   orgID,
			"resource": res,
			"action":   act,
			"env":      m.env,
		})
		log.Trace("checking permission")
		idn, err := m.cl.ValidatePermission(ctx, &ppb.ValidatePermissionRequest{
			ID:          hydra.ClientID(ctx),
			Action:      act,
			Resource:    res,
			OrgID:       orgID,
			Environment: m.env,
		})
		if err != nil {
			log.Trace("permission denied")
			return nil, status.Error(codes.PermissionDenied, "not authorized")
		}
		log.WithField("identity", idn).Trace("permission valid")
		ctx = metautils.ExtractIncoming(ctx).Set(AuthIDKey, idn.GetID()).
			Set(OrgIDKey, idn.GetOrgID()).Set(IDNameKey, idn.GetName()).
			ToIncoming(ctx)
		return handler(ctx, req)
	}
}
