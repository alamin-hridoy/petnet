package handler

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"errors"
	tmpl "html/template"
	"io/fs"
	"net/http"
	"net/url"
	"path"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/benbjohnson/hashfs"
	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"

	openapi "github.com/go-openapi/runtime/client"
	"github.com/ory/hydra-client-go/client/admin"
	"github.com/ory/hydra-client-go/models"
	"google.golang.org/grpc"

	"brank.as/rbac/idp/auth"
	"brank.as/rbac/serviceutil/logging"

	ipb "brank.as/rbac/gunk/v1/invite"
	upb "brank.as/rbac/gunk/v1/user"
)

// genericErrMsg is the error to return to public client to avoid leaking internal implementation details.
const genericErrMsg = "unknown error"

type transport struct {
	Transport http.RoundTripper
}

type rt func(*http.Request) (*http.Response, error)

func (t rt) RoundTrip(req *http.Request) (*http.Response, error) { return t(req) }

func fakeTLS(r http.RoundTripper) http.RoundTripper {
	return rt(func(req *http.Request) (*http.Response, error) {
		rq := req.Clone(req.Context())
		rq.Header.Set("X-Forwarded-Proto", "https")
		return r.RoundTrip(rq)
	})
}

func newHydra(hydraURL *url.URL) admin.ClientService {
	ht := openapi.New(hydraURL.Host, hydraURL.Path, []string{hydraURL.Scheme})
	ht.Transport = fakeTLS(http.DefaultTransport)
	return admin.New(ht, nil)
}

type githubConf struct {
	ClientID     string
	ClientSecret string
	Scopes       string
}

type URLs struct {
	UserMmg string `json:"user_mgm"`
}

type ServerOption func(*server)

func noop() ServerOption { return func(*server) {} }

// WithAuthenticator sets the authenticator for the IDP handler.
func WithDefaultAuthenticator(authenticator auth.Authenticator) ServerOption {
	return func(s *server) { s.authClient[""] = authenticator }
}

// WithAuthenticator sets the authenticator for the IDP handler.
func WithAuthenticator(name string, authenticator auth.Authenticator) ServerOption {
	return func(s *server) { s.authClient[name] = authenticator }
}

// WithLoginTemplate sets the template to be used for the login page.
func WithLoginTemplate(loginTpl string) ServerOption {
	return func(s *server) { s.loginTpl = s.tmpl.Lookup(loginTpl) }
}

func WithAdminLoginTemplate(adminLoginTpl string) ServerOption {
	return func(s *server) { s.adminLoginTpl = s.tmpl.Lookup(adminLoginTpl) }
}

// WithConsentTemplate sets the template to be used for the consent page.
func WithConsentTemplate(consentTpl string) ServerOption {
	return func(s *server) { s.consentTpl = s.tmpl.Lookup(consentTpl) }
}

func WithSetPasswordTemplate(setPasswordTpl string) ServerOption {
	return func(s *server) { s.setPasswordTpl = s.tmpl.Lookup(setPasswordTpl) }
}

func WithSetPasswordSuccessTemplate(setPasswordSuccessTpl string) ServerOption {
	return func(s *server) { s.setPasswordSuccessTpl = s.tmpl.Lookup(setPasswordSuccessTpl) }
}

func WithInviteSetPasswordTemplate(inviteSetPasswordTpl string) ServerOption {
	return func(s *server) { s.inviteSetPasswordTpl = s.tmpl.Lookup(inviteSetPasswordTpl) }
}

func WithInviteSetPasswordSuccessTemplate(inviteSetPasswordSuccessTpl string) ServerOption {
	return func(s *server) { s.inviteSetPasswordSuccessTpl = s.tmpl.Lookup(inviteSetPasswordSuccessTpl) }
}

func WithForgotPasswordTemplate(forgotPasswordTpl string) ServerOption {
	return func(s *server) { s.forgotPasswordTpl = s.tmpl.Lookup(forgotPasswordTpl) }
}

func WithConfirmEmailTemplate(confirmEmailTpl string) ServerOption {
	return func(s *server) { s.confirmEmailTpl = s.tmpl.Lookup(confirmEmailTpl) }
}

func WithErrorTemplate(errTpl string) ServerOption {
	return func(s *server) { s.errTpl = s.tmpl.Lookup(errTpl) }
}

func WithForgotPwdConfirmTemplate(forgotPwdConfirmTpl string) ServerOption {
	return func(s *server) { s.forgotPwdConfirmTpl = s.tmpl.Lookup(forgotPwdConfirmTpl) }
}

// WithLogoutTemplate sets the template to be used for the logout page.
func WithLogoutTemplate(logoutTpl string) ServerOption {
	return func(s *server) { s.logoutTpl = s.tmpl.Lookup(logoutTpl) }
}

// WithRegisterPersonalInfoTemplate sets the template to be used for the
// proxtera register page.
func WithRegisterPersonalInfoTemplate(registerTpl string) ServerOption {
	return func(s *server) { s.registerPersInfoTpl = s.tmpl.Lookup(registerTpl) }
}

// WithSignupTemplate sets the template to be used for the
// user signup page.
func WithSignupTemplate(signupTpl string) ServerOption {
	return func(s *server) { s.signupTpl = s.tmpl.Lookup(signupTpl) }
}

// WithOTPTemplate sets the template to be used for the
// user OTP page.
func WithOTPTemplate(otpTpl string) ServerOption {
	return func(s *server) { s.otpTpl = s.tmpl.Lookup(otpTpl) }
}

// WithHydraClientForURL creates a Hydra client for the given URL.
func WithHydraClientForURL(hydraURL *url.URL) ServerOption {
	return func(s *server) { s.hydra = newHydra(hydraURL) }
}

// WithAuthURL sets the auth host URL to build links against
func WithAuthURL(u string) ServerOption { return func(s *server) { s.u.AuthURL = u } }

// WithSignupURL sets the signup URL.
func WithSignupURL(u string) ServerOption { return func(s *server) { s.u.SignupURL = u } }

// WithSiteURL sets the base project site URL.
func WithSiteURL(u string) ServerOption { return func(s *server) { s.u.SiteURL = u } }

// WithLoginURL sets the project site login URL.  URL should immediately redirect to login flow.
func WithLoginURL(u string) ServerOption { return func(s *server) { s.u.LoginURL = u } }

// WithUsermgmURL sets usermgm base URL.
func WithUsermgmURL(u string) ServerOption { return func(s *server) { s.u.UserMgmURL = u } }

// WithErrorRedirectURL configures IDP to redirects to given URL if an error
// happens in login/consent.
func WithErrorRedirectURL(url string) ServerOption { return func(s *server) { s.ErrRedirURL = url } }

// WithCSRFAuthKey configures IDP to protect POST forms with CSRF token.
func WithCSRFAuthKey(key []byte) ServerOption { return func(s *server) { s.csrfAuthKey = key } }

// WithCSRFSecureDisabled configures IDP to disable secure cookie, intended to use for development environment only.
func WithCSRFSecureDisabled() ServerOption { return func(s *server) { s.csrfSecureDisabled = true } }

// WithEnvironment sets the environment.
func WithEnvironment(env string) ServerOption { return func(s *server) { s.environment = env } }

// WithLogoutConfirmation configures server to show logout confirmation page. By default, it is skipped.
func WithLogoutConfirmation() ServerOption { return func(s *server) { s.logoutConfirmation = true } }

// WithCookieStore adds a cookiestore.
func WithCookieStore(cs *cookieStore) ServerOption { return func(s *server) { s.cookieStore = cs } }

// WithIdentityClients sets up connection to identity.
func WithIdentityClients(conn *grpc.ClientConn) ServerOption {
	return func(s *server) {
		s.suCl = upb.NewSignupClient(conn)
		s.invCl = ipb.NewInviteServiceClient(conn)
	}
}

func WithOpenIDFields(field ...string) ServerOption {
	return func(s *server) {
		for _, f := range field {
			s.openidfields[f] = true
		}
	}
}

// WithProjectName sets the name of the project and uses it in e.g templates.
func WithProjectName(pn string) ServerOption { return func(s *server) { s.projectName = pn } }

// WithLoginRetries sets the number of retries before redirecting to error page.
func WithLoginRetries(lr int) ServerOption { return func(s *server) { s.loginRetries = lr } }

func WithMiddleware(mw func(http.Handler) http.Handler) ServerOption {
	return func(s *server) { s.use = append(s.use, mw) }
}

// WithDisableSignup if disabled makes the page /signup not visitable.
func WithDisableSignup(disable bool) ServerOption {
	return func(s *server) { s.disableSignup = disable }
}

// WithServiceName sets the name of the service.
func WithServiceName(sn string) ServerOption {
	return func(s *server) { s.svcName = sn }
}

type server struct {
	loginTpl                    *tmpl.Template
	consentTpl                  *tmpl.Template
	logoutTpl                   *tmpl.Template
	adminLoginTpl               *tmpl.Template
	registerPersInfoTpl         *tmpl.Template
	signupTpl                   *tmpl.Template
	setPasswordTpl              *tmpl.Template
	setPasswordSuccessTpl       *tmpl.Template
	inviteSetPasswordTpl        *tmpl.Template
	inviteSetPasswordSuccessTpl *tmpl.Template
	forgotPasswordTpl           *tmpl.Template
	forgotPwdConfirmTpl         *tmpl.Template
	confirmEmailTpl             *tmpl.Template
	errTpl                      *tmpl.Template
	otpTpl                      *tmpl.Template

	u urls

	tmpl    *tmpl.Template
	assetFS *hashfs.FS

	authClient         map[string]auth.Authenticator
	hydra              admin.ClientService
	ErrRedirURL        string
	logoutConfirmation bool
	environment        string
	projectName        string
	openidfields       map[string]bool
	loginRetries       int
	disableSignup      bool

	csrfAuthKey        []byte
	csrfSecureDisabled bool
	rememberFor        int64
	dcd                *schema.Decoder
	cookieStore        *cookieStore

	suCl  upb.SignupClient
	invCl ipb.InviteServiceClient

	svcName string

	use []func(http.Handler) http.Handler
}

type urls struct {
	AuthURL    string
	SignupURL  string
	SiteURL    string
	LoginURL   string
	UserMgmURL string
}

// New creates a new HTTP handler for idp service.
func New(tFS, assetFS fs.FS, opts ...ServerOption) (http.Handler, error) {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	gob.Register(identity{})
	s := &server{
		dcd:          decoder,
		openidfields: map[string]bool{},
		authClient:   map[string]auth.Authenticator{},
		use:          []func(http.Handler) http.Handler{},
		assetFS:      hashfs.NewFS(assetFS),
	}
	t, err := tmpl.New("").Funcs(tmpl.FuncMap{
		"assetHash": func(n string) string {
			return path.Join("/assets/",
				s.assetFS.HashName(strings.TrimPrefix(path.Clean(n), "/assets/")),
			)
		},
	}).Funcs(sprig.FuncMap()).ParseFS(tFS, "*.html")
	if err != nil {
		return nil, err
	}
	s.tmpl = t
	for _, o := range opts {
		o(s)
	}
	switch {
	case s.loginTpl == nil:
		return nil, errors.New("login template is required")
	case s.logoutTpl == nil:
		return nil, errors.New("logout template is required")
	case s.authClient == nil:
		return nil, errors.New("authenticator is required")
	case s.hydra == nil:
		return nil, errors.New("hydra client is required")
	case len(s.authClient) == 0:
		return nil, errors.New("auth client is required")
	}

	r := mux.NewRouter()

	r.Use(otelmux.Middleware(s.svcName))
	if s.csrfAuthKey != nil {
		var csrfOptions []csrf.Option
		if s.csrfSecureDisabled {
			csrfOptions = append(csrfOptions, csrf.Secure(false))
		}
		csrfOptions = append(csrfOptions, csrf.Path("/"))
		r.Use(csrf.Protect(s.csrfAuthKey, csrfOptions...))
	}
	for _, u := range s.use {
		r.Use(u)
	}

	r.PathPrefix("/assets/").Handler(
		http.StripPrefix("/assets/",
			http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				log := logging.FromContext(r.Context()).
					WithField("method", "assets").WithField("path", r.URL.Path)
				_, err := fs.Stat(s.assetFS, strings.TrimPrefix(r.URL.Path, "/"))
				if err != nil {
					logging.WithError(err, log).Debug("stat error")
					http.NotFound(w, r)
					return
				}
				w.Header().Set("Cache-Control", "max-age=86400")
				if _, h := hashfs.ParseName(r.URL.Path); h != "" {
					// if asset is hashed extend cache to 180 days
					w.Header().Set("Cache-Control", "max-age=15552000")
				}
				http.FileServer(http.FS(s.assetFS)).ServeHTTP(w, r)
			}),
		),
	)

	r.Path("/login").Methods(http.MethodGet).HandlerFunc(s.getLoginHandler)
	r.Path("/login").Methods(http.MethodPost).HandlerFunc(s.postLoginHandler)
	r.Path("/login/otp").Methods(http.MethodPost).HandlerFunc(s.postOTPLogin)

	if !s.disableSignup {
		r.Path("/signup").Methods(http.MethodGet).HandlerFunc(s.getSignupHandler)
		r.Path("/signup").Methods(http.MethodPost).HandlerFunc(s.postSignupHandler)
	}

	r.Path("/consent").Methods(http.MethodGet).HandlerFunc(s.getConsent)
	r.Path("/consent").Methods(http.MethodPost).HandlerFunc(s.postConsent)

	r.Path("/logout").Methods(http.MethodGet).HandlerFunc(s.getLogoutHandler)
	r.Path("/logout").Methods(http.MethodPost).HandlerFunc(s.postLogoutHandler)

	r.Path("/user/set-password").Methods(http.MethodGet).HandlerFunc(s.getSetPasswordHandler)
	r.Path("/user/set-password").Methods(http.MethodPost).HandlerFunc(s.postSetPasswordHandler)
	r.Path("/user/set-password-success").Methods(http.MethodGet).HandlerFunc(s.getSetPasswordSuccessHandler)

	r.Path("/invite/set-password").Methods(http.MethodGet).HandlerFunc(s.getInviteSetPassword)
	r.Path("/invite/set-password").Methods(http.MethodPost).HandlerFunc(s.postInviteSetPassword)
	r.Path("/invite/set-password-success").Methods(http.MethodGet).HandlerFunc(s.getInviteSetPasswordSuccess)

	r.Path("/forgot-password").Methods(http.MethodGet).HandlerFunc(s.getForgotPasswordHandler)
	r.Path("/forgot-password").Methods(http.MethodPost).HandlerFunc(s.postForgotPasswordHandler)
	r.Path("/user/confirm").Methods(http.MethodGet).HandlerFunc(s.getConfirmEmailHandler)

	r.Path("/otp").Methods(http.MethodGet).HandlerFunc(s.getOTPHandler)

	r.Path("/").Methods(http.MethodGet).
		Handler(http.RedirectHandler(s.u.AuthURL+"/login", http.StatusPermanentRedirect))

	r.Path("/urls.json").Methods(http.MethodGet).HandlerFunc(s.urlsJSONHandler)

	return r, nil
}

func (s *server) lookup(ctx context.Context, name string, def *tmpl.Template) *tmpl.Template {
	if tmpl := s.tmpl.Lookup(name); tmpl != nil {
		return tmpl
	}
	return def
}

func (s *server) urlsJSONHandler(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "application/json")
	urls := URLs{UserMmg: s.u.UserMgmURL}
	json.NewEncoder(res).Encode(urls)
}

type cookieStore struct {
	*sessions.CookieStore
}

func NewCookieStore(config *viper.Viper) (*cookieStore, error) {
	cookieSecret := config.GetString("server.cookieStoreSecret")
	if cookieSecret == "" {
		return nil, errors.New("missing cookie secret")
	}

	cookies := sessions.NewCookieStore([]byte(cookieSecret))
	cookies.Options.HttpOnly = true
	cookies.Options.MaxAge = config.GetInt("server.cookieMaxAge")
	cookies.Options.Secure = config.GetString("runtime.environment") != "development"

	return &cookieStore{
		CookieStore: cookies,
	}, nil
}

func (s *server) authCl(w http.ResponseWriter, r *http.Request, ctx context.Context, hcl *models.OAuth2Client) auth.Authenticator {
	log := logging.FromContext(ctx).WithField("method", "authcl")
	cfg, err := auth.ParseClientConfig(hcl)
	if err != nil || cfg.Authenticator == "" {
		log.WithError(err).Debug("using default authenticator")
		return s.authClient[""]
	}
	if b := s.authClient[cfg.Authenticator]; b != nil {
		return b
	}
	if b := s.authClient[""]; b != nil {
		log.WithField("authenticator", cfg.Authenticator).Debug("falling back to default")
		return b
	}
	log.WithField("client meta", cfg).WithField("client list", s.authClient).Error("not found")
	http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
	return nil
}
