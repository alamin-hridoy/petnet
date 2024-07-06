package hydra

import (
	"context"
	"fmt"
	"strings"
	"time"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	adm "github.com/ory/hydra-client-go/client/admin"
)

type payload struct {
	Active    bool
	Extra     map[string]string
	Scope     []string
	Subject   string
	Audience  []string
	TokenType string
	Username  string
	Expiry    time.Time
}

type admin struct {
	cl adm.ClientService
}

func (h *admin) Introspect(ctx context.Context) (*payload, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}
	return h.IntrospectToken(ctx, token)
}

func convertClientMeta(m interface{}) map[string]string {
	md, ok := m.(map[string]interface{})
	if !ok {
		return map[string]string{}
	}
	mp := make(map[string]string, len(md))
	for k, v := range md {
		if val, ok := v.(string); ok {
			mp[k] = val
		}
	}
	return mp
}

func (h *admin) IntrospectToken(ctx context.Context, token string) (*payload, error) {
	params := adm.NewIntrospectOAuth2TokenParams().WithToken(token)
	pl, err := h.cl.IntrospectOAuth2Token(params)
	if err != nil {
		return nil, fmt.Errorf("hydra instropection failed: %w", err)
	}

	// Hydra returns an empty payload an no error if token not found
	if pl.Payload.Sub == "" {
		return nil, ErrInvalidToken
	}
	ext := convertClientMeta(pl.Payload.Ext)

	return &payload{
		Audience:  pl.Payload.Aud,
		Active:    *pl.Payload.Active,
		Extra:     ext,
		Scope:     strings.Split(pl.Payload.Scope, " "),
		Subject:   pl.Payload.Sub,
		TokenType: pl.Payload.TokenType,
		Username:  pl.Payload.Username,
		Expiry:    time.Unix(pl.Payload.Exp, 0),
	}, nil
}
