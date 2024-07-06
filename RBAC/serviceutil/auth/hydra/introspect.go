package hydra

import (
	"context"
	"fmt"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/ory/hydra/sdk/go/hydra/client/admin"
)

type payload struct {
	Active    bool
	Extra     map[string]interface{}
	Scope     []string
	Subject   string
	Audience  []string
	TokenType string
	Username  string
}

type introspect struct {
	cl *admin.Client
}

func (h *introspect) Introspect(ctx context.Context) (*payload, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}
	params := admin.NewIntrospectOAuth2TokenParams().WithToken(token)
	pl, err := h.cl.IntrospectOAuth2Token(params, nil)
	if err != nil {
		return nil, fmt.Errorf("hydra instropection failed: %w", err)
	}

	// Hydra returns an empty payload an no error if token not found
	if pl.Payload.Subject == "" {
		return nil, ErrInvalidToken
	}

	return &payload{
		Audience:  pl.Payload.Audience,
		Active:    *pl.Payload.Active,
		Extra:     pl.Payload.Extra,
		Scope:     strings.Split(pl.Payload.Scope, " "),
		Subject:   pl.Payload.Subject,
		TokenType: pl.Payload.TokenType,
		Username:  pl.Payload.Username,
	}, nil
}
