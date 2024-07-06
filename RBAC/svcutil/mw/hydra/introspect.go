package hydra

import (
	"context"
	"fmt"
	"strings"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"github.com/ory/hydra-client-go/client/admin"
)

type payload struct {
	Active    bool
	Extra     map[string]string
	Scope     []string
	Subject   string
	Audience  []string
	TokenType string
	Username  string
}

func (s *Service) Introspect(ctx context.Context) (*payload, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}
	params := admin.NewIntrospectOAuth2TokenParams().WithToken(token)
	pl, err := s.cl.IntrospectOAuth2Token(params)
	if err != nil {
		return nil, fmt.Errorf("hydra instropection failed: %w", err)
	}

	// Hydra returns an empty payload an no error if token not found
	if pl.Payload.Sub == "" {
		return nil, ErrInvalidToken
	}
	// ensure map is not nil
	ext := map[string]string{}
	if e, ok := pl.Payload.Ext.(map[string]string); ok {
		ext = e
	}
	return &payload{
		Audience:  pl.Payload.Aud,
		Active:    *pl.Payload.Active,
		Extra:     ext,
		Scope:     strings.Split(pl.Payload.Scope, " "),
		Subject:   pl.Payload.Sub,
		TokenType: pl.Payload.TokenType,
		Username:  pl.Payload.Username,
	}, nil
}
