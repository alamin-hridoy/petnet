package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"github.com/laher/mergefs"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mainpkg"
	"brank.as/rbac/svcutil/metrics"
	"brank.as/rbac/svcutil/mw"
	"brank.as/rbac/svcutil/otelb"

	"brank.as/rbac/idp/assets"
	"brank.as/rbac/idp/auth"
	"brank.as/rbac/idp/auth/dummy"
	"brank.as/rbac/idp/auth/usermgm"
	"brank.as/rbac/idp/handler"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"

	authpb "brank.as/rbac/gunk/v1/authenticate"
	upb "brank.as/rbac/gunk/v1/user"
)

var (
	version = "devel"
	service = "idp"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log := logging.NewLogger().WithFields(logrus.Fields{"service": service, "version": version})
	config, err := mainpkg.DefaultConfig()
	if err != nil {
		log.WithError(err).Fatalf("configuration loading")
	}

	switch config.GetString("LOG_LEVEL") {
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	case "trace":
		log.Logger.SetLevel(logrus.TraceLevel)
	}

	ctx := logging.WithLogger(context.Background(), log)
	shutdown := otelb.InitOTELProvider(ctx, service, config.GetString("trace.collectorHost"))
	defer shutdown()

	met, err := metrics.NewInfluxDBClient(config)
	if err != nil {
		return err
	}
	defer met.Close()
	go met.ErrorsFunc(func(e error) { log.WithError(e).Error("influxdb") })

	hydraURL, err := url.Parse(config.GetString("hydra.adminURL"))
	if err != nil {
		return err
	}

	be, err := initBackends(config)
	if err != nil {
		return err
	}
	authOpts := []handler.ServerOption{}
	for k := range be {
		if k == "" {
			authOpts = append(authOpts, handler.WithDefaultAuthenticator(be[k]))
			continue
		}
		authOpts = append(authOpts, handler.WithAuthenticator(k, be[k]))
	}
	cookieStore, err := handler.NewCookieStore(config)
	if err != nil {
		return err
	}

	retryInterval := config.GetDuration("grpc.retryInterval")
	retryMax := config.GetUint("grpc.retryMax")
	retryopts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(retryInterval * time.Millisecond)),
		grpc_retry.WithMax(retryMax),
		grpc_retry.WithCodes(codes.Unavailable),
	}

	ctx, cancel := context.WithTimeout(context.Background(),
		config.GetDuration("authenticator.dialtimeout"))
	idConn, err := grpc.DialContext(ctx, config.GetString("usermgm.api"),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryopts...)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithUnaryInterceptor(injectMDInterceptor(config)),
	)
	if err != nil {
		cancel()
		return err
	}
	cancel()

	cfg := config.GetString
	def := func(o, d string) string {
		if o != "" {
			return o
		}
		return d
	}
	siteURL, err := url.Parse(cfg("site.url"))
	if err != nil {
		siteURL = &url.URL{Path: "/"}
	}
	sitePath := func(p string) string {
		u, err := siteURL.Parse(p)
		if err != nil {
			return siteURL.String()
		}
		return u.String()
	}

	opts := []handler.ServerOption{
		handler.WithDisableSignup(config.GetBool("handler.disableSignup")),
		handler.WithHydraClientForURL(hydraURL),
		handler.WithEnvironment(cfg("runtime.environment")),
		handler.WithCookieStore(cookieStore),
		handler.WithIdentityClients(idConn),
		handler.WithOpenIDFields("userid", "orgid"),
		handler.WithServiceName(service),

		// Templates
		handler.WithLoginTemplate(cfg("tmpl.login")),
		handler.WithConsentTemplate(cfg("tmpl.consent")),
		handler.WithLogoutTemplate(cfg("tmpl.logout")),
		handler.WithSetPasswordTemplate(cfg("tmpl.setPassword")),
		handler.WithSetPasswordSuccessTemplate(cfg("tmpl.setPasswordSuccess")),
		handler.WithInviteSetPasswordTemplate(cfg("tmpl.inviteSetPassword")),
		handler.WithInviteSetPasswordSuccessTemplate(cfg("tmpl.inviteSetPasswordSuccess")),
		handler.WithForgotPasswordTemplate(cfg("tmpl.forgotPassword")),
		handler.WithForgotPwdConfirmTemplate(cfg("tmpl.forgotPasswordConfirmation")),
		handler.WithConfirmEmailTemplate(cfg("tmpl.confirmEmail")),
		handler.WithErrorTemplate(cfg("tmpl.error")),
		handler.WithRegisterPersonalInfoTemplate(cfg("tmpl.registerPersonalInfo")),
		handler.WithSignupTemplate(cfg("tmpl.signup")),
		handler.WithOTPTemplate(cfg("tmpl.otp")),

		// URLs
		handler.WithAuthURL(def(cfg("base.url"), cfg("site.url"))),
		handler.WithSiteURL(cfg("site.url")),
		handler.WithSignupURL(def(cfg("signup.url"), sitePath("/signup"))),
		handler.WithLoginURL(def(cfg("login.url"), sitePath("/login"))),
		handler.WithUsermgmURL(cfg("usermgm.url")),
		handler.WithErrorRedirectURL(cfg("server.errorurl")),

		// Middleware
		handler.WithMiddleware(met.HTTPMiddleware("")),
		handler.WithMiddleware(mw.Logger(log)),
		handler.WithMiddleware(mw.Gzip),
		handler.WithMiddleware(mw.ContentType("text/html")),
		handler.WithMiddleware(mw.Recovery()),

		// TODO: Deprecate
		handler.WithProjectName(config.GetString("project.name")),
		handler.WithLoginRetries(config.GetInt("login.retries")),
	}
	if csrfAuthKey := config.GetString("server.csrfAuthKey"); csrfAuthKey != "" {
		opts = append(opts, handler.WithCSRFAuthKey([]byte(csrfAuthKey)))
		if config.GetString("runtime.environment") == "development" {
			opts = append(opts, handler.WithCSRFSecureDisabled())
		}
	}
	opts = append(opts, authOpts...)

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	aFS := mergefs.Merge(afero.NewIOFS(
		afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(wd, "assets")),
	), assets.Static)
	tFS := mergefs.Merge(afero.NewIOFS(
		afero.NewBasePathFs(afero.NewOsFs(), filepath.Join(wd, "assets", "templates")),
	), assets.Templates)
	h, err := handler.New(tFS, aFS, opts...)
	if err != nil {
		return err
	}

	cl, err := HydraCleanup(config)
	if err != nil {
		return err
	}
	sched := mainpkg.NewCronTicker(func() time.Duration {
		if d := config.GetDuration("hydra.cleanup"); d != 0 {
			return d
		}
		return time.Second
	}())
	cleanOpt := mainpkg.WithCron("hydra", sched, cl.Clean)
	if config.GetString("server.path") != "" {
		cleanOpt = mainpkg.WithLeaderCron("hydra", sched, cl.Clean)
	}

	wto := 5 * time.Second
	if to := config.GetDuration("server.writeTimeout"); to != 0 {
		wto = to
	}

	srv, err := mainpkg.Setup(config, log,
		mainpkg.WithTimeout(wto),
		mainpkg.WithWebHandler(mainpkg.HTTPOnly(h)),
		cleanOpt,
	)
	srv.Run()
	return nil
}

func initBackends(conf *viper.Viper) (map[string]auth.Authenticator, error) {
	auths := conf.GetStringSlice("authenticator.names")
	def := conf.GetString("authenticator.default")
	to := conf.GetDuration("authenticator.dialtimeout")

	cf := func(name, field string) string {
		return conf.GetString(strings.Join([]string{"authenticator", name, field}, "."))
	}

	retryInterval := conf.GetDuration("grpc.retryInterval")
	retryMax := conf.GetUint("grpc.retryMax")
	retryopts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(retryInterval * time.Millisecond)),
		grpc_retry.WithMax(retryMax),
		grpc_retry.WithCodes(codes.Unavailable),
	}
	dialopts := []grpc.DialOption{
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryopts...)),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy":"round_robin"}`),
		grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(to),
		grpc.WithChainUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	}

	be := make(map[string]auth.Authenticator, len(auths))
	for _, a := range auths {
		switch cf(a, "type") {
		case "session", "":
			host := cf(a, "host")
			conn, err := grpc.Dial(host, dialopts...)
			if err != nil {
				return nil, fmt.Errorf("dialing %s: %w", a, err)
			}
			var consent auth.ConsentGrantor
			if chost := cf(a, "consent"); chost != "" {
				if chost == host {
					consent = auth.NewConsentGrantor(conn)
				} else {
					cconn, err := grpc.Dial(chost, dialopts...)
					if err != nil {
						return nil, fmt.Errorf("dialing %s: %w", a, err)
					}
					consent = auth.NewConsentGrantor(cconn)
				}
			}
			be[a] = usermgm.NewSessionAuth(
				authpb.NewSessionServiceClient(conn),
				conf.GetDuration(strings.Join([]string{"authenticator", a, "duration"}, ".")),
				consent,
			)
		case "usermgm":
			inconn, err := grpc.Dial(cf(a, "internal"), dialopts...)
			if err != nil {
				return nil, fmt.Errorf("dialing internal %s: %w", a, err)
			}

			exconn, err := grpc.Dial(cf(a, "api"), dialopts...)
			if err != nil {
				return nil, fmt.Errorf("dialing %s: %w", a, err)
			}
			authdef, err := usermgm.NewAuthenticator(
				upb.NewUserAuthServiceClient(inconn),
				upb.NewUserServiceClient(exconn),
				conf.GetDuration(strings.Join([]string{"authenticator", a, "duration"}, ".")),
			)
			if err != nil {
				return nil, err
			}
			be[a] = authdef
		case "dummy":
			be[a] = dummy.New()
		}
		if a == def {
			be[""] = be[a]
		}
	}
	return be, nil
}

func injectMDInterceptor(cnf *viper.Viper) grpc.UnaryClientInterceptor {
	cred := clientcredentials.Config{
		ClientID:     cnf.GetString("auth.clientID"),
		ClientSecret: cnf.GetString("auth.clientSecret"),
		TokenURL:     cnf.GetString("auth.url") + "/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	ts := cred.TokenSource(context.Background())
	return func(c context.Context, m string, rq, rp interface{}, cc *grpc.ClientConn, inv grpc.UnaryInvoker, o ...grpc.CallOption) error {
		tok, err := ts.Token()
		if err != nil {
			return err
		}
		c = metautils.ExtractOutgoing(c).
			Set("authorization", "Bearer "+tok.AccessToken).ToOutgoing(c)
		return inv(c, m, rq, rp, cc, o...)
	}
}
