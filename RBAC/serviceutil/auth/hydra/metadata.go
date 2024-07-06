package hydra

import (
	"context"
	"strings"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	XApiKey = "x-api-key"
)

type extraData struct{}

// GetExtra returns the extra info from hydra session.
func GetExtra(ctx context.Context) map[string]interface{} {
	ex := ctx.Value(&extraData{})
	if ex != nil {
		return ex.(map[string]interface{})
	}
	return nil
}

// ClientID extracts the client ID from the context.
func ClientID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(ClientIDKey)
}

// Metadata loads hydra session info into metadata
func (s *Service) Metadata(ctx context.Context) (context.Context, error) {
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
	}

	// claims should be backwards compatible with dedicated hydra interceptor
	md.Add(ClientIDKey, p.Subject)
	metadata.MD(md).Set(AudienceKey, p.Audience...)
	metadata.MD(md).Set(TokenScopesKey, p.Scope...)

	return md.ToIncoming(ctx), nil
}

func (s *Service) PublicEndpoint(method string) bool {
	fullMethod := strings.TrimPrefix(method, "/")
	_, ok := s.ignoredMethods[fullMethod]
	return ok
}
