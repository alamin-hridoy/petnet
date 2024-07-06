package hydra

import (
	"context"
	"strings"

	session "brank.as/petnet/profile/services/rbsession"
	"brank.as/petnet/svcutil/mw/meta"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	XApiKey = "x-api-key"

	// ClientIDKey is the context metadata's key for client IDs.
	// Value is being set by 'ValidateToken', always set for valid tokens.
	ClientIDKey = "client-id"

	// AudienceKey is the context metadata's key for audience.
	// Value is being set by 'ValidateToken', only set if service is initialized
	// with known audience.
	AudienceKey = "audience"

	OrgTypeKey = "org-type"

	TokenScopesKey = "token-scopes"

	// EnvironmentKey the environment the authentication is permitted
	EnvironmentKey = "environment"
)

type extraData struct{}

// GetExtra returns the extra info from hydra session.
func GetExtra(ctx context.Context) map[string]string {
	ex := ctx.Value(&extraData{})
	if ex != nil {
		return ex.(map[string]string)
	}
	return nil
}

// OrgType extracts the organization type from context.
func OrgType(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(OrgTypeKey)
}

// ClientID extracts the client ID from the context.
func ClientID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(ClientIDKey)
}

// Audience extracts the audience from the context.
func Audience(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(AudienceKey)
}

// Metadata loads hydra session info into metadata
func (s *Service) GRPC() meta.RequestMeta {
	return meta.MetaFunc(func(ctx context.Context) (context.Context, error) {
		md := metautils.ExtractIncoming(ctx)

		// token authorization
		p, err := s.Introspect(ctx)
		if err != nil {
			if s.optional {
				return ctx, nil
			}
			return nil, status.Error(codes.Unauthenticated, codes.Unauthenticated.String())
		}
		if len(p.Extra) > 0 {
			ctx = context.WithValue(ctx, &extraData{}, p.Extra)
			if env, ok := p.Extra[EnvironmentKey]; ok {
				metadata.MD(md).Set(EnvironmentKey, env)
			}
		}

		// claims should be backwards compatible with dedicated hydra interceptor
		md.Add(ClientIDKey, p.Subject)
		metadata.MD(md).Set(AudienceKey, p.Audience...)
		metadata.MD(md).Set(TokenScopesKey, p.Scope...)

		return md.ToIncoming(ctx), nil
	})
}

// InternalGRPC loads hydra session info into metadata for internal use
func (s *Service) InternalGRPC() meta.RequestMeta {
	return meta.MetaFunc(func(ctx context.Context) (context.Context, error) {
		md := metautils.ExtractIncoming(ctx)
		e := GetExtra(ctx)
		if len(e) > 0 {
			ot, ok := e[session.OrgType]
			if ok {
				md.Add(OrgTypeKey, ot)
			}
		}
		return md.ToIncoming(ctx), nil
	})
}

func (s *Service) PublicEndpoint(method string) bool {
	fullMethod := strings.TrimPrefix(method, "/")
	_, ok := s.ignoredMethods[fullMethod]
	return ok
}
