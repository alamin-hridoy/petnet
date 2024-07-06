package hydra

import (
	"context"
	"fmt"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	UsernameKey = "username"
	OrgIDKey    = "org-id"
)

// Username extracts the username from the context.
func Username(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(UsernameKey)
}

// OrgID extracts the org ID from the context.
func OrgID(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(OrgIDKey)
}

type UserDetails struct {
	// when set to true, the metadata function will not return a fatal error
	opt bool
}

// Metadata loads user and ord identifiers into the request metadata.
func (ud *UserDetails) Metadata(ctx context.Context) (context.Context, error) {
	extra := GetExtra(ctx)
	if extra == nil {
		if ud.opt {
			return ctx, nil
		}
		return nil, status.Error(codes.Unauthenticated, "session details not found")
	}

	md := metautils.ExtractIncoming(ctx)
	if u, err := extractString(extra, "username"); err == nil {
		md.Add(UsernameKey, u)
	}
	if o, err := extractString(extra, "org_id"); err == nil {
		md.Add(OrgIDKey, o)
	}
	if !ud.opt && (md.Get(OrgIDKey) == "" || md.Get(UsernameKey) == "") {
		return nil, status.Error(codes.Unauthenticated, "session invalid")
	}
	return md.ToIncoming(ctx), nil
}

// UserInfo contains a logged-in user, is returning by a instropector.
type UserInfo struct {
	Username string
	OrgID    string
}

// extractString extracts a string claim from Hydra Introspect result extra with given key in parameter.
func extractString(extra map[string]string, key string) (string, error) {
	result, ok := extra[key]
	if !ok {
		return "", fmt.Errorf("missing %q in auth token", key)
	}
	return result, nil
}
