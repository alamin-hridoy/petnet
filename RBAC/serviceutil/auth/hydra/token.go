package hydra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
)

const defaultNsAudience string = "default"

// Claims extracted returned by hydra token introspect endpoint
type Claims struct {
	Active bool `json:"active"`
	// TODO:
	// this field will probably be invalidated once OB-3148 rolled out;
	// adding it as quick fix for now to quickly deliver the business requirements
	// that tokens for sandbox and live shouldn't be interchangeable,
	// see OB-3540 and OB-3584 - deprecate as need be.
	Audience   []string    `json:"audience,omitempty"`
	Scope      string      `json:"scope,omitempty"`
	ClientID   string      `json:"client_id,omitempty"`
	Subject    string      `json:"sub,omitempty"`
	Expiration int64       `json:"exp,omitempty"`
	IssuedAt   int64       `json:"iat,omitempty"`
	Issuer     string      `json:"iss,omitempty"`
	TokenType  string      `json:"token_type,omitempty"`
	Extra      ExtraClaims `json:"ext,omitempty"`
}

// ExtraClaims that we set in the IDP during login
type ExtraClaims struct {
	Company   string `json:"company,omitempty"`
	Email     string `json:"email,omitempty"`
	Subdomain string `json:"subdomain,omitempty"`
	Username  string `json:"username,omitempty"`
}

// ValidateToken validate the token using hydra token introspections endpoint
// and return the claims if the token is active.
func (s *Service) ValidateToken(ctx context.Context, token string) (*Claims, error) {
	if token == "" {
		return nil, errMissingToken
	}

	body := []byte(fmt.Sprintf(`token=%s`, token))
	claims, err := s.dispatchRequest(ctx, http.MethodPost, s.hydraIntrospectURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	// Check if token is active.
	// Only token that is not expired will return as true,
	// expired token will mark as inactive token.
	if claims == nil || !claims.Active {
		return nil, ErrInvalidToken
	}
	if claims.ClientID == "" {
		return nil, errMissingClientID
	}

	// only retrieve claim's audience if service has pre-configured known audiences
	// and only from the namespaced hydra server
	if len(s.knownAudience) == 0 || s.hydraGetClientURL == "" {
		return claims, nil
	}

	u := fmt.Sprintf(s.hydraGetClientURL, claims.ClientID)
	res, err := s.dispatchRequest(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	claims.Audience = res.Audience
	return claims, nil
}

func (s *Service) dispatchRequest(ctx context.Context, method, url string, body io.Reader) (*Claims, error) {
	logger := logging.FromContext(ctx)

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	// X-Forwaded-Proto use to connect to proxy or load balancer
	// as hydra server run on k8s and we will connect to the loadbalancer
	// through https protocol.
	req.Header.Set("X-Forwarded-Proto", "https")

	resp, err := s.HttpClient.Do(req)
	if err != nil {
		logging.WithError(err, logger).Error("unable to make request to Auth Server")
		return nil, err
	}
	defer func() {
		if err := resp.Body.Close(); err != nil {
			logging.WithError(err, logger).Error("error close the response body")
		}
	}()

	response := &Claims{}
	if err := json.NewDecoder(resp.Body).Decode(response); err != nil {
		logging.WithError(err, logger).Error("unable to decode the response")
		return nil, err
	}

	return response, nil
}

// ValidateTokenFromCtx validate the token from its context by
// extracting the authorization information
func (s *Service) ValidateTokenFromCtx(ctx context.Context) (*Claims, error) {
	token, err := grpc_auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "missing token")
	}
	return s.ValidateToken(ctx, token)
}
