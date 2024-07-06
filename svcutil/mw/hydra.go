package mw

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	phmw "brank.as/petnet/api/perahub-middleware"
	trmSvc "brank.as/petnet/api/services/terminal"
	apiutil "brank.as/petnet/api/util"
	"brank.as/petnet/gunk/drp/v1/terminal"
	dppb "brank.as/petnet/gunk/dsa/v1/profile"
	ptnrLst "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfSvc "brank.as/petnet/gunk/dsa/v2/service"
	spb "brank.as/petnet/gunk/v1/session"
	session "brank.as/petnet/profile/services/rbsession"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	"github.com/gorilla/sessions"
	"github.com/knq/jwt"
	"golang.org/x/oauth2"
)

const (
	sessionCookieName     = "petnet-session"
	sessionCookieState    = "state"
	sessionCookieEmail    = "email"
	sessionCookieUsername = "username"
	sessionCookieToken    = "token"
	sessionCookieRefresh  = "refresh"
	sessionCookieExpiry   = "expiry"
	sessionCookieUserID   = "user-id"
	sessionCookieOrgID    = "org-id"
	sourcePath            = "source-path"
	authCodeURL           = "somerandomstring"
	timeLayout            = "Jan. 2, 2006 15:04:05 MST"
	SessionAuthType       = "authtype"
)

const (
	BasicAuthType  = "basic"
	OAuth2AuthType = "oauth2"
)

const (
	DefaultLogin    = "/login"
	DefaultLogout   = "/logout"
	DefaultCallback = "/oauth2/callback"
	DefaultError    = "/error"
	DefaultRedirect = "/"
)

type (
	userIDKey      struct{}
	orgTypeKey     struct{}
	orgIDKey       struct{}
	petnetOwnerKey struct{}
	providerKey    struct{}
	sessionExpiry  struct{}
)

type Hydra struct {
	*oauth2.Config
	*sessions.CookieStore
	cl     *hydra.Service
	cookie string
	// path for the login handler
	loginPath string
	// path to redirect after login, if session is empty.
	redirPath string
	// path to redirect after auth error.
	errorPath string
	logout    *url.URL
	scl       spb.SessionServiceClient

	ignore         map[string]bool
	ignorePrefix   []string
	optional       map[string]bool
	optionalPrefix []string
}

type Config struct {
	LoginPath         string
	PostLoginRedirect string
	ErrorPath         string
	IgnorePaths       []string
	IgnorePrefix      []string
	OptionalPaths     []string
	OptionalPrefix    []string
}

func NewHydra(config *viper.Viper, mwConfig Config, pfConn *grpc.ClientConn) (*Hydra, error) {
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
	cookies := sessions.NewCookieStore([]byte(cookieSecret))
	cookies.Options.HttpOnly = true
	cookies.Options.MaxAge = config.GetInt("auth.cookieMaxAge")
	cookies.Options.Secure = config.GetBool("auth.cookieSecure")

	// OAuth2
	issuer, err := url.Parse(config.GetString("auth.issuer"))
	if issuer.String() == "" || err != nil {
		return nil, errors.New("missing valid auth issuer")
	}
	issuer.Path = path.Join(issuer.Path, "oauth2", "auth")
	logout, err := issuer.Parse(path.Join("/", "oauth2", "sessions", "logout"))
	if logout.String() == "" || err != nil {
		return nil, errors.New("invalid logout")
	}
	token, err := url.Parse(config.GetString("auth.token"))
	if token.String() == "" || err != nil {
		return nil, errors.New("missing valid auth token URL")
	}
	token.Path = path.Join(token.Path, "oauth2", "token")
	clientID := config.GetString("auth.clientID")
	if clientID == "" {
		return nil, errors.New("missing client ID")
	}
	clientSecret := config.GetString("auth.clientSecret")
	if clientSecret == "" {
		return nil, errors.New("missing client secret")
	}
	oauthURL, err := url.Parse(config.GetString("auth.oauthURL"))
	if oauthURL.String() == "" || err != nil {
		return nil, errors.New("missing valid oauth URL")
	}
	oauthURL.Path = path.Join(oauthURL.Path, "callback")

	opt := func(conf, env, def string) string {
		switch {
		case conf != "":
			return conf
		case env != "":
			return env
		}
		return def
	}

	loginPath := opt(mwConfig.LoginPath, config.GetString("auth.loginPath"), DefaultLogin)
	errorPath := opt(mwConfig.ErrorPath, config.GetString("auth.errorPath"), DefaultError)
	redirPath := opt(mwConfig.PostLoginRedirect,
		config.GetString("auth.postLoginRedirect"), DefaultRedirect)
	cookieName := opt("", config.GetString("auth.cookieName"), sessionCookieName)

	ig := make(map[string]bool, len(mwConfig.IgnorePaths))
	for _, p := range mwConfig.IgnorePaths {
		ig[p] = true
	}

	opp := make(map[string]bool, len(mwConfig.OptionalPaths))
	for _, p := range mwConfig.OptionalPaths {
		opp[p] = true
	}

	return &Hydra{
		cl:        h,
		scl:       spb.NewSessionServiceClient(pfConn),
		loginPath: loginPath,
		errorPath: errorPath,
		redirPath: redirPath,
		logout:    logout,
		cookie:    cookieName,

		ignore:         ig,
		ignorePrefix:   mwConfig.IgnorePrefix,
		optional:       opp,
		optionalPrefix: mwConfig.OptionalPrefix,

		CookieStore: cookies,
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  oauthURL.String(),
			Scopes:       []string{"offline_access", "openid"},
			Endpoint: oauth2.Endpoint{
				AuthURL:   issuer.String(),
				TokenURL:  token.String(),
				AuthStyle: oauth2.AuthStyleInHeader,
			},
		},
	}, nil
}

type refreshToken struct{}

// Middleware authenticates requests and refreshes session.
func (hy *Hydra) DSASimMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if hy.ignore[path] {
			h.ServeHTTP(w, r)
			return
		}
		if hy.ignorePrefix != nil {
			for _, ignoredPath := range hy.ignorePrefix {
				if strings.HasPrefix(path, ignoredPath) {
					h.ServeHTTP(w, r)
					return
				}
			}
		}
		ctx := r.Context()
		log := logging.FromContext(ctx).WithField("Path", path)

		sess, err := hy.Get(r, hy.cookie)
		if err != nil {
			if strings.HasPrefix(path, "/auth-option") || path == "/" {
				h.ServeHTTP(w, r)
				return
			}
			logging.WithError(err, log).Error("missing session")
			http.Redirect(w, r, "/auth-option", http.StatusTemporaryRedirect)
			return
		}
		authType, ok := sess.Values[SessionAuthType].(string)
		if !ok {
			log.Error("missing auth type")
			http.Redirect(w, r, "/auth-option", http.StatusTemporaryRedirect)
			return
		}
		if authType == BasicAuthType {
			h.ServeHTTP(w, r)
			return
		}

		idnt, err := hy.validateToken(ctx, sess, path)
		if err != nil {
			if strings.HasPrefix(path, "/auth-option") || path == "/" {
				h.ServeHTTP(w, r)
				return
			}
			logging.WithError(err, log).Error("invalid token")
			http.Redirect(w, r, "/auth-option", http.StatusTemporaryRedirect)
			return
		}

		ctx = context.WithValue(ctx, userIDKey{}, idnt.Subject)
		ctx = context.WithValue(ctx, orgTypeKey{}, idnt.Extra[session.OrgType])
		ctx = context.WithValue(ctx, orgIDKey{}, idnt.Extra[session.OrgID])
		ctx = context.WithValue(ctx, petnetOwnerKey{}, idnt.Extra[session.PetnetOwner])

		tok, err := hy.refresh(ctx, sess)
		if err != nil {
			logging.WithError(err, log).Error("refreshing session")
			http.Redirect(w, r, "/auth-option", http.StatusTemporaryRedirect)
			return
		}
		log.WithField("expiry", tok.Expiry).Trace("valid token")
		sess.Values[sessionCookieToken] = tok.AccessToken
		sess.Values[sessionCookieRefresh] = tok.RefreshToken
		sess.Values[sessionCookieExpiry] = tok.Expiry.Format(time.RFC3339)
		ctx = context.WithValue(ctx, &refreshToken{}, tok)

		if _, err = hy.scl.SetSessionExpiry(ctx, &spb.SetSessionExpiryRequest{
			IDType: spb.IDType_USERID,
			ID:     idnt.Subject,
			Expiry: timestamppb.New(tok.Expiry),
		}); err != nil {
			logging.WithError(err, log).Error("setting session expiration")
		}

		idnt, err = hy.validateToken(ctx, sess, path)
		if err != nil {
			logging.WithError(err, log).Error("invalid token")
			http.Redirect(w, r, "/auth-option", http.StatusTemporaryRedirect)
			return
		}

		ctx = context.WithValue(ctx, userIDKey{}, idnt.Subject)
		ctx = context.WithValue(ctx, orgTypeKey{}, idnt.Extra[session.OrgType])
		ctx = context.WithValue(ctx, orgIDKey{}, idnt.Extra[session.OrgID])
		ctx = context.WithValue(ctx, petnetOwnerKey{}, idnt.Extra[session.PetnetOwner])

		r = r.WithContext(ctx)

		if err := sess.Save(r, w); err != nil {
			logging.WithError(err, log).Error("save refresh session")
		}
		h.ServeHTTP(w, r)
	})
}

// MockMiddleware adds user data to context for local development.
func MockMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ctx = context.WithValue(ctx, userIDKey{}, "10000000-0000-0000-0000-000000000000")
		ctx = context.WithValue(ctx, orgIDKey{}, "20000000-0000-0000-0000-000000000000")
		r = r.WithContext(ctx)
		h.ServeHTTP(w, r)
	})
}

func getExpiry(ctx context.Context) time.Time {
	expiry, ok := ctx.Value(&sessionExpiry{}).(time.Time)
	if !ok {
		return time.Time{}
	}
	return expiry
}

func GetUserID(ctx context.Context) string {
	userID, ok := ctx.Value(userIDKey{}).(string)
	if !ok {
		return ""
	}
	return userID
}

func GetOrgType(ctx context.Context) string {
	orgType, ok := ctx.Value(orgTypeKey{}).(string)
	if !ok {
		return ""
	}
	return orgType
}

func GetOrgID(ctx context.Context) string {
	orgID, ok := ctx.Value(orgIDKey{}).(string)
	if !ok {
		return ""
	}
	return orgID
}

func IsPetnetOwner(ctx context.Context) bool {
	v, ok := ctx.Value(petnetOwnerKey{}).(string)
	if !ok || v != "true" {
		return false
	}
	return true
}

func IsProvider(ctx context.Context) bool {
	v, ok := ctx.Value(providerKey{}).(string)
	if !ok || v != "true" {
		return false
	}
	return true
}

// Exchange code and load session state.
func (hy *Hydra) Exchange(w http.ResponseWriter, r *http.Request) error {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	sess, err := hy.Get(r, sessionCookieName)
	if err != nil {
		logging.WithError(err, log).Error("missing cookie")
		http.Error(w, "cookie error", http.StatusInternalServerError)
		return err
	}

	token, err := hy.Config.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		logging.WithError(err, log).Error("exchange")
		http.Error(w, "cookie Exchange error", http.StatusInternalServerError)
		return err
	}
	rawIDToken, ok := token.Extra("id_token").(string)
	ui := &userInfo{}
	if ok {
		ui = decodeOpenID(r.Context(), sess, rawIDToken)
	} else {
		log.Error("missing openid token")
	}

	sess.Values[sessionCookieUserID] = ui.UserID
	sess.Values[sessionCookieToken] = token.AccessToken
	sess.Values[sessionCookieRefresh] = token.RefreshToken
	sess.Values[sessionCookieExpiry] = token.Expiry.Format(time.RFC3339)

	if _, err = hy.scl.SetSessionExpiry(ctx, &spb.SetSessionExpiryRequest{
		IDType: spb.IDType_USERID,
		ID:     ui.UserID,
		Expiry: timestamppb.New(token.Expiry),
	}); err != nil {
		logging.WithError(err, log).Error("setting session expiration")
	}

	if err := sess.Save(r, w); err != nil {
		logging.WithError(err, log).Error("session save")
		http.Error(w, "cookie error", http.StatusInternalServerError)
		return err
	}

	return nil
}

type userInfo struct {
	UserID   string
	Email    string
	Username string
}

func decodeOpenID(ctx context.Context, sess *sessions.Session, payload string) *userInfo {
	log := logging.FromContext(ctx)
	oid := &jwt.UnverifiedToken{}
	if err := jwt.DecodeUnverifiedToken([]byte(payload), oid); err != nil {
		logging.WithError(err, log).Error("decode openid Token")
		return nil
	}
	buf, err := base64.RawStdEncoding.DecodeString(string(oid.Payload))
	if err != nil {
		logging.WithError(err, log).WithField("payload", string(oid.Payload)).
			Error("decode oid payload")
		return nil
	}
	pl := map[string]interface{}{}
	if err := json.Unmarshal(buf, &pl); err != nil {
		logging.WithError(err, log).WithField("payload", string(oid.Payload)).Error("oid payload")
		return nil
	}
	em := fmt.Sprintf("%v", pl["email"])
	un := fmt.Sprintf("%v", pl["username"])
	uid := fmt.Sprintf("%v", pl["userid"])
	sess.Values[sessionCookieEmail] = em
	sess.Values[sessionCookieUsername] = un
	return &userInfo{
		UserID:   uid,
		Email:    em,
		Username: un,
	}
}

// refresh ...
func (hy *Hydra) refresh(ctx context.Context, sess *sessions.Session) (*oauth2.Token, error) {
	log := logging.FromContext(ctx)

	at, _ := sess.Values[sessionCookieToken].(string)
	rt, _ := sess.Values[sessionCookieRefresh].(string)
	tm, _ := sess.Values[sessionCookieExpiry].(string)
	switch "" {
	case at, rt, tm:
		log.Trace("access", at, "refresh", rt, "expiry", tm)
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
		log := logging.FromContext(r.Context())
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
		sess, err := hy.Get(r, sessionCookieName)
		if err != nil {
			logging.WithError(err, log).Error("missing session")
			http.Redirect(w, r, "/root", http.StatusTemporaryRedirect)
			return
		}
		token, ok := sess.Values[sessionCookieToken].(string)
		if !ok {
			log.WithField("session", sess.Values).Error("missing token")
			http.Redirect(w, r, "/root", http.StatusTemporaryRedirect)
			return
		}
		if tok, ok := ctx.Value(&refreshToken{}).(*oauth2.Token); ok {
			token = tok.AccessToken
		}
		log.WithField("auth_token", token).Trace("loaded")
		ctx = context.WithValue(ctx, &authToken{}, token)
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthForwarder() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		token, ok := ctx.Value(&authToken{}).(string)
		if !ok || token == "" {
			return status.Error(codes.Unauthenticated, "missing forward auth token")
		}
		ctx = metautils.ExtractOutgoing(ctx).Add("Authorization", "Bearer "+token).ToOutgoing(ctx)
		ctx, cancel := context.WithTimeout(ctx, 20*time.Second)
		defer cancel()
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

type identity struct {
	Active    bool
	Extra     map[string]string
	Scope     []string
	Subject   string
	Audience  []string
	TokenType string
	Username  string
	Expiry    time.Time
}

func (hy *Hydra) Session(r *http.Request) (*sessions.Session, error) { return hy.Get(r, hy.cookie) }

func (hy *Hydra) storePath(w http.ResponseWriter, r *http.Request, sess *sessions.Session) {
	if r.URL.Path == "/" {
		return
	}
	log := logging.FromContext(r.Context()).WithField("method", "hydra.storePath")
	log.WithField("source path", r.URL.Path).Debug("storing path")
	sess.Values[sourcePath] = r.URL.Path
	if err := hy.Save(r, w, sess); err != nil {
		logging.WithError(err, log).Error("save session")
		return
	}
}

func (hy *Hydra) skip(r *http.Request) bool {
	if hy.ignore[r.URL.Path] {
		return true
	}
	for _, p := range hy.ignorePrefix {
		if strings.HasPrefix(r.URL.Path, p) {
			return true
		}
	}
	return false
}

func (hy *Hydra) optionalReq(r *http.Request) bool {
	if hy.optional[r.URL.Path] {
		return true
	}
	for _, p := range hy.optionalPrefix {
		if strings.HasPrefix(r.URL.Path, p) {
			return true
		}
	}
	return false
}

// Middleware authenticates requests and refreshes session.
func (hy *Hydra) Middleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hy.skip(r) {
			h.ServeHTTP(w, r)
			return
		}
		ctx, opt := r.Context(), hy.optionalReq(r)
		log := logging.FromContext(ctx).WithFields(logrus.Fields{
			"Path": r.URL.Path, "method": "hydra.Middleware",
		})

		sess, err := hy.Session(r)
		if err != nil {
			log = logging.WithError(err, log)
			if opt {
				log.Debug("missing session")
				h.ServeHTTP(w, r)
				return
			}
			log.Error("missing session")
			http.Redirect(w, r, hy.loginPath, http.StatusTemporaryRedirect)
			return
		}

		rq, err := hy.auth(w, r, sess)
		if err == nil {
			h.ServeHTTP(w, rq)
			return
		}
		log = logging.WithError(err, log)
		if opt {
			log.Debug("auth session")
			h.ServeHTTP(w, r)
			return
		}
		uid, ok := sess.Values[sessionCookieUserID].(string)
		if !ok {
			log.Error("missing user id")
			h.ServeHTTP(w, r)
			return
		}
		if _, err = hy.scl.SetSessionExpiry(ctx, &spb.SetSessionExpiryRequest{
			IDType: spb.IDType_USERID,
			ID:     uid,
			Expiry: timestamppb.New(getExpiry(ctx)),
		}); err != nil {
			logging.WithError(err, log).Error("setting session expiration")
		}
		log.Error("auth session")
		hy.storePath(w, r, sess)
		http.Redirect(w, r, hy.loginPath, http.StatusTemporaryRedirect)
	})
}

func (hy *Hydra) auth(w http.ResponseWriter, r *http.Request, sess *sessions.Session) (*http.Request, error) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithFields(logrus.Fields{
		"Path": r.URL.Path, "method": "hydra.auth",
	})
	t, ok := sess.Values[sessionCookieToken].(string)
	if !ok {
		return nil, fmt.Errorf("missing token")
	}
	ctx, err := hy.introspect(ctx, t)
	if err != nil {
		logging.WithError(err, log).Error("validate error")
	}

	tok, err := hy.refresh(ctx, sess)
	if err != nil {
		return nil, fmt.Errorf("refresh error: %w", err)
	}
	sess.Values[sessionCookieToken] = tok.AccessToken
	sess.Values[sessionCookieRefresh] = tok.RefreshToken
	sess.Values[sessionCookieExpiry] = tok.Expiry.Format(time.RFC3339)
	sess.Values[sessionCookieOrgID] = GetOrgID(ctx)

	ctx = context.WithValue(ctx, &refreshToken{}, tok)
	ctx, err = hy.introspect(ctx, tok.AccessToken)
	if err != nil {
		logging.WithError(err, log).Error("refresh token error")
	}

	r = r.WithContext(ctx)
	if err := sess.Save(r, w); err != nil {
		logging.WithError(err, log).Error("save refresh session")
	}
	return r, nil
}

func (hy *Hydra) introspect(ctx context.Context, tok string) (context.Context, error) {
	log := logging.FromContext(ctx).WithField("method", "hydra.introspect")

	cl, err := hy.cl.IntrospectToken(ctx, tok)
	switch {
	case err != nil:
		return ctx, err
	case cl == nil:
		return ctx, fmt.Errorf("token introspect: %w", err)
	case !cl.Active:
		return ctx, fmt.Errorf("token inactive: %w", err)
	case time.Until(cl.Expiry) < 10*time.Second:
		return ctx, fmt.Errorf("token expired: %w", err)
	}
	log.WithField("expiry", cl.Expiry).Trace("valid token")

	ctx = context.WithValue(ctx, userIDKey{}, cl.Subject)
	if len(cl.Extra) == 0 {
		return ctx, nil
	}

	ctx = context.WithValue(ctx, orgTypeKey{}, cl.Extra[session.OrgType])
	ctx = context.WithValue(ctx, orgIDKey{}, cl.Extra[session.OrgID])
	ctx = context.WithValue(ctx, petnetOwnerKey{}, cl.Extra[session.PetnetOwner])
	ctx = context.WithValue(ctx, providerKey{}, cl.Extra[session.Provider])
	return ctx, nil
}

func (hy *Hydra) validateToken(ctx context.Context, sess *sessions.Session, path string) (*identity, error) {
	log := logging.FromContext(ctx).WithField("Path", path)

	token, ok := sess.Values[sessionCookieToken].(string)
	if !ok {
		log.Error("missing token")
		return nil, fmt.Errorf("missing token")
	}
	idnt, err := hy.cl.IntrospectToken(ctx, token)
	if err != nil {
		logging.WithError(err, log).Error("introspecting token")
		return nil, err
	}
	if len(idnt.Extra) == 0 {
		return nil, fmt.Errorf("missing extra")
	}

	_, ok = idnt.Extra[session.OrgType]
	if !ok {
		return nil, fmt.Errorf("missing org type")
	}
	_, ok = idnt.Extra[session.OrgID]
	if !ok {
		return nil, fmt.Errorf("missing org id")
	}
	_, ok = idnt.Extra[session.PetnetOwner]
	if !ok {
		return nil, fmt.Errorf("missing petnet owner")
	}

	return &identity{
		Active:    idnt.Active,
		Extra:     idnt.Extra,
		Scope:     idnt.Scope,
		Subject:   idnt.Subject,
		Audience:  idnt.Audience,
		TokenType: idnt.TokenType,
		Username:  idnt.Username,
		Expiry:    idnt.Expiry,
	}, nil
}

const (
	TerminalService      string = "terminal.TerminalService"
	CashInCashOutService string = "cashincashout.CashInCashOutService"
	BillsPaymentService  string = "bills_payment.BillspaymentService"
	RuralNetService      string = "microinsurance.MicroInsuranceService"
	RuralNetParter       string = "RuralNet"
	Provider             string = "provider"
)

var excludeValMap = map[string]string{
	pfSvc.ServiceType_REMITTANCE.String(): "ListRemit",
}

func getPartnerFormat(ptnr string) string {
	if len(ptnr) > 0 && ptnr[0] == '"' {
		ptnr = ptnr[1:]
	}
	if len(ptnr) > 0 && ptnr[len(ptnr)-1] == '"' {
		ptnr = ptnr[:len(ptnr)-1]
	}
	return ptnr
}

// getTerminalPartner for getting terminal partner from request
func getTerminalPartner(log *logrus.Entry, req interface{}, key string) string {
	b, err := json.Marshal(req)
	if err != nil {
		logging.WithError(err, log).Error("marshaling")
		return ""
	}
	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(b, &objmap); err != nil {
		logging.WithError(err, log).Error("unmarshaling")
		return ""
	}

	ptnr := getPartnerFormat(string(objmap[key]))

	return ptnr
}

// getCiCoPartner for getting cash in cash out partner from request
func getCiCoPartner(log *logrus.Entry, req interface{}) string {
	b, err := json.Marshal(req)
	if err != nil {
		logging.WithError(err, log).Error("marshaling")
		return ""
	}
	var objmap map[string]json.RawMessage
	if err := json.Unmarshal(b, &objmap); err != nil {
		logging.WithError(err, log).Error("unmarshaling")
		return ""
	}

	ptnr := getPartnerFormat(string(objmap["provider"]))
	if ptnr != "" {
		return ptnr
	}

	if trx, ok := objmap["trx"]; ok {
		var trxObj map[string]json.RawMessage
		if err := json.Unmarshal(trx, &trxObj); err != nil {
			logging.WithError(err, log).Error("unmarshaling")
			return ""
		}
		ptnr = getPartnerFormat(string(trxObj["provider"]))

		if ptnr != "" {
			return ptnr
		}
	}

	return ""
}

// getPartnerFormator for getting cash in cash out partner from request
func getPartnerFormator(dptnrs []*ptnrLst.DSAPartnerList) map[string][]string {
	ptnr := map[string][]string{}
	if len(dptnrs) == 0 {
		return ptnr
	}
	for _, v := range dptnrs {
		ptnr[v.GetTransactionType()] = append(ptnr[v.GetTransactionType()], v.GetPartner())
	}
	return ptnr
}

func InSlice(val interface{}, slce interface{}) (exists bool) {
	exists = false
	switch reflect.TypeOf(slce).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(slce)
		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) {
				exists = true
				return
			}
		}
	}
	return
}

// Todo(Kiran) Pull this out of svcutil and put it on API.
func ValidateAccess(cl pfSvc.ServiceServiceClient, pl ptnrLst.PartnerListServiceClient, tl *trmSvc.Svc) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := logging.FromContext(ctx).WithField("method", "hydra.ValidateAccess")
		env := phmw.GetEnv(ctx)
		apiEnv := phmw.GetAPIEnv(ctx)
		orgType := phmw.GetOrgInfo(ctx)

		if env != apiEnv {
			log.Error("env not match")
			return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
		}

		var needValidate bool
		var partner, service string
		switch {
		case strings.Contains(info.FullMethod, TerminalService):
			needValidate = true
			service = pfSvc.ServiceType_REMITTANCE.String()
			method := strings.Split(info.FullMethod, "/")
			if method[len(method)-1] == excludeValMap[service] {
				sr, err := cl.ValidateServiceAccess(ctx, &pfSvc.ValidateServiceAccessRequest{
					OrgID:               hydra.OrgID(ctx),
					SvcName:             service,
					IsAnyPartnerEnabled: true,
				})
				if err != nil {
					logging.WithError(err, log).Error("validate service access")
					return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
				}
				if !sr.Enabled {
					log.Error("partner not found in service access")
					return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
				}
				return handler(ctx, req)
			}
			ptnr := getTerminalPartner(log, req, "remit_partner")
			transactionID := getTerminalPartner(log, req, "transaction_id")
			if ptnr == "" && transactionID != "" {
				ptnrRes, err := tl.GetPartnerByTxnID(ctx, &terminal.GetPartnerByTxnIDRequest{
					TransactionID: transactionID,
				})
				if err != nil {
					logging.WithError(err, log).Error("get partner by txn id")
					return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
				}
				ptnr = ptnrRes.GetPartner()
			}

			if ptnr == "" {
				log.Error("partner not found")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}
			partner = ptnr
		case strings.Contains(info.FullMethod, CashInCashOutService):
			needValidate = true
			service = pfSvc.ServiceType_CASHINCASHOUT.String()
			ptnr := getCiCoPartner(log, req)
			if ptnr == "" {
				log.Error("partner not found")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}
			partner = ptnr
		case strings.Contains(info.FullMethod, BillsPaymentService):
			if orgType == Provider {
				log.Error("providers not allowed to access bills payment")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}
			needValidate = true
			service = pfSvc.ServiceType_BILLSPAYMENT.String()
			ptnr := getTerminalPartner(log, req, "partner")
			if ptnr == "" {
				log.Error("partner not found")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}
			partner = ptnr
		case strings.Contains(info.FullMethod, RuralNetService):
			if orgType == Provider {
				log.Error("providers not allowed to access rural net")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}
			needValidate = true
			service = pfSvc.ServiceType_MICROINSURANCE.String()
			partner = RuralNetParter
		default:
			needValidate = false
		}

		if !needValidate || GetOrgID(ctx) == dppb.OrgType_PetNet.String() {
			return handler(ctx, req)
		}

		log = log.WithFields(logrus.Fields{
			"method":   info.FullMethod,
			"env":      phmw.GetEnv(ctx),
			"apiEnv":   phmw.GetAPIEnv(ctx),
			"orgId":    hydra.OrgID(ctx),
			"clientId": hydra.ClientID(ctx),
			"partner":  partner,
		})

		ctx = metautils.ExtractIncoming(ctx).Set(phmw.Partner, partner).ToIncoming(ctx)

		lpr := &ptnrLst.GetPartnerListRequest{
			Status:      "ENABLED",
			ServiceName: service,
		}

		if service == pfSvc.ServiceType_MICROINSURANCE.String() {
			lpr.Name = partner
		} else {
			lpr.Stype = partner
		}

		gpl, err := pl.GetPartnerList(ctx, lpr)
		if err != nil {
			logging.WithError(err, log).Error("get partner list")
			return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
		}

		if gpl == nil {
			log.WithError(err).Error("partner list is nil")
			return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
		}

		if len(gpl.PartnerList) == 0 {
			log.WithError(err).Error("partner list is empty")
			return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
		}

		if service == pfSvc.ServiceType_REMITTANCE.String() {
			transactiontype := phmw.GetTransactionTypes(ctx)
			if transactiontype == "" {
				log.Error("transaction type not found")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}

			gDPLst, err := pl.GetDSAPartnerList(ctx, &ptnrLst.DSAPartnerListRequest{
				TransactionType: transactiontype,
			})
			if err != nil {
				logging.WithError(err, log).Error("get dsa partner list")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}

			getDsaPtnrs := getPartnerFormator(gDPLst.GetDSAPartnerList())
			if !InSlice(partner, getDsaPtnrs[transactiontype]) {
				log.Error("partner not found in dsa partner list")
				return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
			}
		}

		vSa, err := cl.ValidateServiceAccess(ctx, &pfSvc.ValidateServiceAccessRequest{
			OrgID:               hydra.OrgID(ctx),
			Partner:             partner,
			SvcName:             service,
			IsAnyPartnerEnabled: false,
		})
		if err != nil {
			logging.WithError(err, log).Error("validate service access")
			return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
		}

		if !vSa.Enabled {
			log.Error("partner not found in service access")
			return nil, apiutil.HandleServiceErr(status.Error(codes.PermissionDenied, "Forbidden partner"))
		}

		return handler(ctx, req)
	}
}
