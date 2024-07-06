package hydra

import (
	"context"
	"errors"
	"strings"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

const (
	// ClientIDKey is the context metadata's key for client IDs.
	// Value is being set by 'ValidateToken', always set for valid tokens.
	ClientIDKey = "client-id"

	// AudienceKey is the context metadata's key for audience.
	// Value is being set by 'ValidateToken', only set if service is initialized
	// with known audience.
	AudienceKey = "audience"

	TokenScopesKey = "token-scopes"
)

var errUnknownAudience = errors.New("unauthorized audience")

// Interface method
type AuthenticateMethod interface {
	AuthenticateMethod(methodName string) bool
}

// Option is option to configure UnaryServerInterceptor method.
type Option func(*Service)

// UnaryServerInterceptor returns a new unary server interceptor for hydra that validates tokens
// and injects claim scopes and clientID into context metadata
func (s *Service) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// TODO: remove after OB-2988 is done
		// For backward compatible. Historically manually defined method/scope maps have no leading "/" character but
		// gunk/scopegen auto generated map has leading "/" in method name.
		fullMethod := strings.TrimPrefix(info.FullMethod, "/")
		if _, ok := s.ignoredMethods[fullMethod]; ok {
			return handler(ctx, req)
		}

		// Check if a service implements the AuthenticateMethod interface to see whether an endpoint
		// requires authentication or not
		if h, ok := info.Server.(AuthenticateMethod); ok && !h.AuthenticateMethod(info.FullMethod) {
			return handler(ctx, req)
		}

		claims, err := s.ValidateTokenFromCtx(ctx)
		if err != nil {
			return nil, status.Error(codes.Unauthenticated, "unauthorized")
		}

		// extract clientID from token
		// TODO: verify why this is the case.
		// The introspect call returns a client ID, why aren't we using that?
		clientID := claims.Subject
		if clientID == "" {
			return nil, status.Error(codes.Unauthenticated, errMissingClientID.Error())
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// Should not happen, metadata is required for token auth
			return nil, status.Error(codes.Unauthenticated, "not authorized")
		}

		// set clientID and scopes as metadata values
		mdDup := md.Copy()
		mdDup.Set(ClientIDKey, clientID)
		claimedScopes := strings.Fields(claims.Scope)
		mdDup.Set(TokenScopesKey, claimedScopes...)
		mdDup.Set(AudienceKey, claims.Audience...)

		ctxWithMD := metadata.NewIncomingContext(ctx, mdDup)

		// set clientID for logging
		grpc_ctxtags.Extract(ctxWithMD).Set(ClientIDKey, clientID)

		// TODO:
		// this check will probably be invalidated once OB-3148 rolled out;
		// adding it as quick fix for now to quickly deliver the business requirements
		// that tokens for sandbox and live shouldn't be interchangeable,
		// see OB-3540 and OB-3584 - deprecate as need be.
		if ok := s.validateAudience(claims.Audience); !ok {
			return nil, status.Error(codes.PermissionDenied, errUnknownAudience.Error())
		}

		// only check matching auth scopes if provided
		if len(s.authScopes) == 0 {
			return handler(ctxWithMD, req)
		}

		// check scope
		requiredScope := s.authScopes[fullMethod]
		for _, scope := range claimedScopes {
			if scope == requiredScope {
				return handler(ctxWithMD, req)
			}
		}

		// scope not found
		return nil, status.Error(codes.PermissionDenied, "missing required scope")
	}
}

// TODO:
// this check will probably be invalidated once OB-3148 rolled out;
// adding it as quick fix for now to quickly deliver the business requirements
// that tokens for sandbox and live shouldn't be interchangeable,
// see OB-3540 and OB-3584 - deprecate as need be.
func (s *Service) validateAudience(claims []string) bool {
	// only check claim if known audience if pre-configured
	if len(s.knownAudience) == 0 {
		return true
	}

	for _, claim := range claims {
		// accept if the token was verified from the default namespace
		// as we currently only have superusers on the AMD
		if claim == defaultNsAudience {
			return true
		}
		// loop through the known audiences and validate claim's audience
		for _, knownAud := range s.knownAudience {
			if claim == knownAud {
				return true
			}
		}
	}
	return false
}

// StreamServerInterceptor returns a new stream server interceptor for hydra that validates tokens
// and injects claim scopes and clientID into context metadata
func (s *Service) StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		newStream := grpc_middleware.WrapServerStream(ss)

		// TODO: remove after OB-2988 is done
		// For backward compatible. Historically manually defined method/scope maps have no leading "/" character but
		// wrap server stream so we can return with context
		// gunk/scopegen auto generated map has leading "/" in method name.
		fullMethod := strings.TrimPrefix(info.FullMethod, "/")
		if _, ok := s.ignoredMethods[fullMethod]; ok {
			return handler(srv, ss)
		}

		claims, err := s.ValidateTokenFromCtx(ctx)
		if err != nil {
			return status.Error(codes.Unauthenticated, "unauthorized")
		}

		// extract clientID from token
		// TODO: verify why this is the case.
		// The introspect call returns a client ID, why aren't we using that?
		clientID := claims.Subject
		if clientID == "" {
			return status.Error(codes.Unauthenticated, "token is missing client ID")
		}

		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			// Should not happen, metadata is required for token auth
			return status.Error(codes.Unauthenticated, "not authorized")
		}

		// set clientID and scopes as metadata values
		mdDup := md.Copy()
		mdDup.Set(ClientIDKey, clientID)
		claimedScopes := strings.Fields(claims.Scope)
		mdDup.Set(TokenScopesKey, claimedScopes...)
		mdDup.Set(AudienceKey, claims.Audience...)

		ctxWithMD := metadata.NewIncomingContext(ctx, mdDup)

		// set clientID for logging
		grpc_ctxtags.Extract(ctxWithMD).Set(ClientIDKey, clientID)

		if ok := s.validateAudience(claims.Audience); !ok {
			return status.Error(codes.Unauthenticated, errUnknownAudience.Error())
		}

		// only check matching auth scopes if provided
		if len(s.authScopes) == 0 {
			newStream.WrappedContext = ctxWithMD
			return handler(srv, newStream)
		}

		requiredScope := s.authScopes[fullMethod]
		for _, scope := range claimedScopes {
			if scope == requiredScope {
				newStream.WrappedContext = ctxWithMD
				return handler(srv, newStream)
			}
		}

		// scope not found
		return status.Error(codes.PermissionDenied, "missing required scope")
	}
}

// MetadataClaimExtractor extracts and returns the granted scopes from context.
func MetadataClaimExtractor(ctx context.Context) []string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil
	}
	return md[TokenScopesKey]
}
