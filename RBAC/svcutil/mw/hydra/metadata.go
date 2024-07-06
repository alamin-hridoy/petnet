package hydra

import (
	"context"
	"errors"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	XApiKey = "x-api-key"

	// UserIDKey is the context metadata's key for client IDs.
	// Value is being set by 'ValidateToken', always set for valid tokens.
	UserIDKey = "client-id"

	// AudienceKey is the context metadata's key for audience.
	AudienceKey = "audience"

	TokenScopesKey = "token-scopes"
)

var errUnknownAudience = errors.New("unauthorized audience")

type extraData struct{}

// GetExtra returns the extra info from hydra session.
func GetExtra(ctx context.Context) map[string]string {
	ex := ctx.Value(&extraData{})
	if ex != nil {
		return ex.(map[string]string)
	}
	return nil
}

// ClientID extracts the client ID from the context.
func ClientID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(UserIDKey)
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
	md.Add(UserIDKey, p.Subject)
	metadata.MD(md).Set(AudienceKey, p.Audience...)
	metadata.MD(md).Set(TokenScopesKey, p.Scope...)

	return md.ToIncoming(ctx), nil
}
