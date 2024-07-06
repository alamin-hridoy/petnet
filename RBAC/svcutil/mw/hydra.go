package mw

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/gorilla/sessions"
	"github.com/kenshaw/jwt"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
)

const (
	sessionCookieName     = "rbac-session"
	sessionCookieState    = "state"
	sessionCookieEmail    = "email"
	sessionCookieUsername = "username"
	sessionCookieToken    = "token"
	sessionCookieRefresh  = "refresh"
	sessionCookieExpiry   = "expiry"
)

type userIDKey struct{}

type Hydra struct {
	*oauth2.Config
	*sessions.CookieStore
	name         string
	cl           *hydra.Service
	ignore       map[string]bool
	ignorePrefix []string
}

func NewHydra(config *viper.Viper, ignorePath []string, ignorePrefix []string) (*Hydra, error) {
	h, err := hydra.NewService(&http.Client{
		Timeout: time.Second,
	}, config.GetString("hydra.adminUrl"))
	if err != nil {
		return nil, err
	}

	// cookie store
	cookieSecret := config.GetString("auth.cookieSecret")
	if cookieSecret == "" {
		return nil, errors.New("missing cookie secret")
	}
	cookieName := config.GetString("auth.cookieName")
	if cookieName == "" {
		return nil, errors.New("missing cookie name")
	}
	cookies := sessions.NewCookieStore([]byte(cookieSecret))

	cookies.Options.HttpOnly = true
	cookies.Options.MaxAge = config.GetInt("auth.cookieMaxAge")
	cookies.Options.Secure = config.GetBool("auth.cookieSecure")

	// OAuth2
	issuer := config.GetString("auth.issuer")
	if issuer == "" {
		return nil, errors.New("missing auth issuer")
	}
	token := config.GetString("auth.token")
	if token == "" {
		return nil, errors.New("missing auth token URL")
	}
	clientID := config.GetString("auth.clientID")
	if clientID == "" {
		return nil, errors.New("missing client ID")
	}
	clientSecret := config.GetString("auth.clientSecret")
	if clientSecret == "" {
		return nil, errors.New("missing client secret")
	}
	redirectURL := config.GetString("auth.redirectURL")
	if redirectURL == "" {
		return nil, errors.New("missing redirect URL")
	}

	ig := make(map[string]bool, len(ignorePath))
	for _, p := range ignorePath {
		ig[p] = true
	}

	return &Hydra{
		cl:           h,
		name:         cookieName,
		ignore:       ig,
		ignorePrefix: ignorePrefix,
		CookieStore:  cookies,
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURL,
			Scopes:       []string{"offline_access", "openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   issuer + "/oauth2/auth",
				TokenURL:  token + "/oauth2/token",
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
	}, nil
}

type refreshToken struct{}

// Middleware authenticates requests and refreshes session.
func (hy *Hydra) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hy.ignore[r.URL.Path] {
			h.ServeHTTP(w, r)
			return
		}
		if hy.ignorePrefix != nil {
			for _, ignoredPath := range hy.ignorePrefix {
				if strings.HasPrefix(r.URL.Path, ignoredPath) {
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		ctx := r.Context()
		log := logging.FromContext(ctx).WithField("Path", r.URL.Path).
			WithField("method", "hydra.AuthMW")
		sess, err := hy.Get(r, hy.name)
		if err != nil {
			log.WithError(err).WithField("name", hy.name).Error("missing session")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		token, ok := sess.Values[sessionCookieToken].(string)
		if !ok {
			log.Error("missing token")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		cl, err := hy.cl.ValidateToken(ctx, token)
		if err != nil {
			log.WithError(err).Error("validate error")
		}

		if cl != nil {
			ctx = context.WithValue(ctx, userIDKey{}, cl.Subject)
		}

		if cl != nil && cl.Active && time.Until(time.Unix(cl.Expiration, 0)) > 10*time.Second {
			tm, _ := sess.Values[sessionCookieExpiry].(string)
			log.WithField("expiry", tm).Info("valid token")
			h.ServeHTTP(w, r.WithContext(ctx))
			return
		}
		tok, err := hy.refresh(ctx, sess)
		if err != nil {
			log.WithError(err).Error("refresh error")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		sess.Values[sessionCookieToken] = tok.AccessToken
		sess.Values[sessionCookieRefresh] = tok.RefreshToken
		sess.Values[sessionCookieExpiry] = tok.Expiry.Format(time.RFC3339)
		ctx = context.WithValue(ctx, &refreshToken{}, tok)

		cl, err = hy.cl.ValidateToken(ctx, tok.AccessToken)
		if err != nil {
			log.WithError(err).Error("validate error")
		}

		if cl != nil {
			ctx = context.WithValue(ctx, userIDKey{}, cl.Subject)
		}
		r = r.WithContext(ctx)

		if err := sess.Save(r, w); err != nil {
			log.WithError(err).Error("save refresh session")
		}
		h.ServeHTTP(w, r)
	})
}

func GetUserID(r *http.Request) (string, error) {
	userID, ok := r.Context().Value(userIDKey{}).(string)
	if !ok {
		return "", errors.New("no userid set")
	}

	return userID, nil
}

// Exchange code and load session state.
func (hy *Hydra) Exchange(w http.ResponseWriter, r *http.Request) error {
	log := logging.FromContext(r.Context()).WithField("method", "hydra.exchangeMW")
	// returns new session on error
	sess, _ := hy.Get(r, hy.name)

	token, err := hy.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		log.WithError(err).Error("exchange")
		http.Error(w, "cookie Exchange error", http.StatusInternalServerError)
		return err
	}
	log.WithField("token", token).Info("exchanged")
	rawIDToken, ok := token.Extra("id_token").(string)
	if ok {
		oid := &jwt.UnverifiedToken{}
		if err := jwt.DecodeUnverifiedToken([]byte(rawIDToken), oid); err != nil {
			log.WithError(err).Error("decode openid Token")
		}
		pl := map[string]interface{}{}
		if err := json.Unmarshal(oid.Payload, &pl); err != nil {
			log.WithError(err).WithField("payload", string(oid.Payload)).Error("oid payload")
		} else {
			sess.Values[sessionCookieEmail] = pl["email"]
			sess.Values[sessionCookieUsername] = pl["username"]
		}
	} else {
		log.Info("missing openid token")
	}

	fmt.Println("session:", sess)

	sess.Values[sessionCookieToken] = token.AccessToken
	sess.Values[sessionCookieRefresh] = token.RefreshToken
	sess.Values[sessionCookieExpiry] = token.Expiry.Format(time.RFC3339)

	if err := sess.Save(r, w); err != nil {
		log.WithError(err).Error("session save")
		http.Error(w, "cookie error", http.StatusInternalServerError)
		return err
	}
	return nil
}

func (hy *Hydra) refresh(ctx context.Context, sess *sessions.Session) (*oauth2.Token, error) {
	log := logging.FromContext(ctx).WithField("method", "hydra.refreshToken")

	at, _ := sess.Values[sessionCookieToken].(string)
	rt, _ := sess.Values[sessionCookieRefresh].(string)
	tm, _ := sess.Values[sessionCookieExpiry].(string)
	switch "" {
	case at, rt, tm:
		log.Info("access", at, "refresh", rt, "expiry", tm)
		return nil, fmt.Errorf("invalid session")
	}
	exp, err := time.Parse(time.RFC3339, tm)
	if err != nil {
		return nil, err
	}
	ts := hy.TokenSource(ctx, &oauth2.Token{
		AccessToken:  at,
		RefreshToken: rt,
		Expiry:       exp,
		TokenType:    "Bearer",
	})
	return ts.Token()
}

type authToken struct{}

func (hy *Hydra) ForwardAuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logging.FromContext(r.Context()).WithField("method", "hydra.forwardAuthMW")
		if hy.ignore[r.URL.Path] {
			h.ServeHTTP(w, r)
			return
		}
		if hy.ignorePrefix != nil {
			for _, ignoredPath := range hy.ignorePrefix {
				if strings.HasPrefix(r.URL.Path, ignoredPath) {
					h.ServeHTTP(w, r)
					return
				}
			}
		}

		ctx := r.Context()
		sess, err := hy.Get(r, hy.name)
		if err != nil {
			log.WithError(err).WithField("name", hy.name).Error("missing session")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		token, ok := sess.Values[sessionCookieToken].(string)
		if !ok {
			log.WithField("session", sess.Values).Error("missing token")
			http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
			return
		}
		if tok, ok := ctx.Value(&refreshToken{}).(*oauth2.Token); ok {
			token = tok.AccessToken
		}
		log.WithField("auth_token", token).Info("loaded")
		ctx = context.WithValue(ctx, &authToken{}, token)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthForwarder() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, ok := ctx.Value(&authToken{}).(string)
		if !ok || token == "" {
			return status.Error(codes.Unauthenticated, "missing forward auth token")
		}
		ctx = metautils.ExtractOutgoing(ctx).Add("Authorization", "Bearer "+token).ToOutgoing(ctx)
		// TODO(PROX-369): make the timeout dynamic depending on endpoint
		// and make the mailer run in parallel
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
