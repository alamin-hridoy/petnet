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
	Optional bool
}

// Metadata loads user and ord identifiers into the request metadata.
func (ud *UserDetails) Metadata(ctx context.Context) (context.Context, error) {
	extra := GetExtra(ctx)
	if extra == nil {
		if ud.Optional {
			return ctx, nil
		}
		return nil, status.Error(codes.Unauthenticated, "session details not found")
	}

	userInfo, err := ud.authViaToken(extra)
	if err != nil {
		if ud.Optional {
			return ctx, nil
		}
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}
	md := metautils.ExtractIncoming(ctx)
	md.Add(UsernameKey, userInfo.Username)
	md.Add(OrgIDKey, userInfo.OrgID)
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

func (ud *UserDetails) authViaToken(extra map[string]string) (*UserInfo, error) {
	username, err := extractString(extra, "username")
	if err != nil {
		return nil, err
	}

	oID, err := extractString(extra, "org_id")
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		OrgID:    oID,
		Username: username,
	}, nil
}
