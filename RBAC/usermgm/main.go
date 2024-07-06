package main

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/viper"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"brank.as/rbac/serviceutil/auth/hydra"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/serviceutil/middleware"
	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/svcutil/mainpkg"
	"brank.as/rbac/svcutil/metrics"
	"brank.as/rbac/svcutil/mw"
	"brank.as/rbac/svcutil/otelb"

	"brank.as/rbac/usermgm/core/auth"
	"brank.as/rbac/usermgm/core/challenge"
	"brank.as/rbac/usermgm/core/mfauth"
	"brank.as/rbac/usermgm/core/oauthclient"
	"brank.as/rbac/usermgm/core/org"
	perm "brank.as/rbac/usermgm/core/permissions"
	"brank.as/rbac/usermgm/core/scopes"
	"brank.as/rbac/usermgm/core/svcacct"
	"brank.as/rbac/usermgm/core/user"

	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/integrations/keto"
	localVal "brank.as/rbac/usermgm/integrations/validation"

	"brank.as/rbac/usermgm/services/invite"
	"brank.as/rbac/usermgm/services/mfa"
	"brank.as/rbac/usermgm/services/oauth2"
	"brank.as/rbac/usermgm/services/organization"
	"brank.as/rbac/usermgm/services/permissions"
	"brank.as/rbac/usermgm/services/product"
	role "brank.as/rbac/usermgm/services/roles"
	scpSvc "brank.as/rbac/usermgm/services/scopes"
	"brank.as/rbac/usermgm/services/signup"
	"brank.as/rbac/usermgm/services/svcaccount"
	userSvc "brank.as/rbac/usermgm/services/user"
	"brank.as/rbac/usermgm/services/userauth"
	"brank.as/rbac/usermgm/services/validation"
	"brank.as/rbac/usermgm/storage/postgres"
)

var (
	version = "devel"
	svcName = "usermgm"
)

func main() {
	log := logging.NewLogger().WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})
	config := viper.NewWithOptions(viper.EnvKeyReplacer(strings.NewReplacer(".", "_")))
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	shutdown := otelb.InitOTELProvider(
		context.Background(),
		svcName,
		config.GetString("trace.collectorHost"))
	defer shutdown()

	met, err := metrics.NewInfluxDBClient(config)
	if err != nil {
		log.Fatal(err)
	}
	defer met.Close()
	go met.ErrorsFunc(func(e error) { log.WithError(err).Debug("influxdb") })

	store, err := newDBFromConfig(config)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure database")
	}
	setupTimeout := 90 * time.Second
	if config.GetDuration("server.setup_timeout") > setupTimeout {
		setupTimeout = config.GetDuration("server.setup_timeout")
	}
	hy, err := initHydra(config)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to connect to hydra")
	}

	ctx, cancel := context.WithTimeout(context.Background(), setupTimeout)
	svc, err := setupServices(ctx, log, config, store)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure grpc services")
	}
	cancel()

	imeta, err := middleware.NewMetadata(log,
		mw.MetaFunc(func(ctx context.Context) (context.Context, error) {
			// Normalize REST headers and gRPC metadata to gRPC metadata.
			nmd := metadata.MD{}
			md := metadata.MD(metautils.ExtractIncoming(ctx))
			for k, v := range md {
				nmd.Append(k, v...)
			}
			return metautils.NiceMD(nmd).ToIncoming(ctx), nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	// serve internal services on admin port
	intl, err := mainpkg.Setup(config, log,
		mainpkg.WithPort(config.GetInt("server.adminPort")),
		mainpkg.WithVersion(svcName, version),
		mainpkg.WithReflection(),
		mainpkg.WithServerOpts(grpc.ChainUnaryInterceptor(met.UnaryServerInterceptor("", nil))),
		mainpkg.WithUnaryInterceptors(
			otelgrpc.UnaryServerInterceptor(),
			imeta.UnaryServerInterceptor(),
		),
		mainpkg.WithServerOpts(grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: false})),
		mainpkg.WithDualService(mainpkg.Internal, svc.internal...),
	)
	if err != nil {
		log.Fatal(err)
	}

	env := config.GetString("permissions.environment")
	vld := mw.NewPermissionValidatorFromClient(svc.local, env)
	meta, err := middleware.NewMetadata(log,
		mw.ResetMD(), // clear keys
		hy, mw.NewServiceAccount(svc.local),
		NewOrg(store),
		mw.ValidateMD(hydra.OrgIDKey), // enforce org in metadata
	)
	if err != nil {
		log.Fatal(err)
	}

	initTO := 15 * time.Second
	if to := config.GetDuration("init.timeout"); to != 0 {
		initTO = to
	}

	svr, err := mainpkg.Setup(config, log,
		mainpkg.WithVersion(svcName, version),
		mainpkg.WithReflection(),
		mainpkg.WithServerOpts(grpc.ChainUnaryInterceptor(
			met.UnaryServerInterceptor("", nil),
			otelgrpc.UnaryServerInterceptor(),
		)),
		mainpkg.WithUnaryInterceptors(
			meta.UnaryServerInterceptor(),
			vld.UnaryServerInterceptor(),
		),
		mainpkg.WithServerOpts(grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: true})),
		mainpkg.WithDualService(mainpkg.External, svc.external...),
		mainpkg.AdditionalServers(intl),
		mainpkg.WithInitTimeout(initTO),
		mainpkg.OptionList(svc.opts),
	)
	if err != nil {
		log.Fatal(err)
	}
	svr.Run()
}

type Svcs struct {
	external []mainpkg.GWGRPC
	internal []mainpkg.GWGRPC
	local    *localVal.Svc
	opts     []mainpkg.Option
}

func setupServices(ctx context.Context, log *logrus.Entry, config *viper.Viper, st *postgres.Storage) (*Svcs, error) {
	env := config.GetString("runtime.environment")
	envLst := config.GetStringSlice("bootstrap.envlist")
	cl, err := client.NewAdminClient(config.GetString("hydra.adminUrl"))
	if err != nil {
		return nil, err
	}
	// bootstrap if necessary
	k := keto.New(net.JoinHostPort(config.GetString("keto.host"), config.GetString("keto.port")))
	p := perm.New(st, k)
	sa := svcacct.New(config, cl, st, p)
	bs, err := Bootstrap(ctx, config, log, st, p, sa)
	if err != nil {
		return nil, err
	}

	mailer, err := initMailer(config)
	if err != nil {
		return nil, err
	}

	ma, err := mfauth.New(config, st, mailer)
	if err != nil {
		return nil, err
	}

	ocl := oauthclient.New(st, cl)
	chlg := challenge.New(st, k)
	og := org.New(st, p, st)
	usr := user.New(user.Config{
		Env:               env,
		PublicSignup:      config.GetBool("org.autocreate"),
		AutoApprove:       config.GetBool("org.autoapprove"),
		UserExistError:    config.GetBool("user.existError"),
		NotifyDisableUser: config.GetBool("user.notifydisable"),
		NotifyEnableUser:  config.GetBool("user.notifyenable"),
		ResetDuration:     config.GetDuration("user.resetdurationsec") * time.Second,
	}, st, st, mailer, og, ma)
	ua := auth.New(config, st, ma, st)
	sp := scopes.New(st, bs)

	sg := signup.New(usr, mailer)
	inv := invite.New(st, st, usr, mailer, invite.WithEnv(env))
	us := userSvc.New(st, usr)
	oa := oauth2.New(ocl, envLst)
	perm := permissions.New(p, chlg)
	roles := role.New(p)
	prd := product.New(p, st, config.GetStringSlice("bootstrap.envlist"))
	org := organization.New(og)

	val := validation.New(chlg)
	ssa := svcaccount.New(sa, st, usr, envLst,
		localVal.NewLocal(val, nil), // Only needs local permission validation.
	)
	auth := userauth.New(st, ua)
	m := mfa.New(ma)
	sc := scpSvc.New(sp, sp)

	lock := make(chan struct{})
	return &Svcs{
		external: []mainpkg.GWGRPC{perm, roles, sg, inv, ssa, us, oa, prd, m, org},
		internal: []mainpkg.GWGRPC{
			val, oa, auth, sc,
			intlSvc{GWGRPC: ssa, intlRegr: ssa},
			intlSvc{GWGRPC: m, intlRegr: m},
		},
		local: localVal.NewLocal(val, ssa),
		opts: []mainpkg.Option{
			mainpkg.WithInit("", func(ctx context.Context) error {
				defer close(lock)
				return st.RunMigration(config.GetString("database.migrationDir"))
			}),
			mainpkg.WithInit("rbac", func(ctx context.Context) error {
				<-lock
				if !config.GetBool("bootstrap.startup") {
					return nil
				}
				return bs.Init(ctx)
			}),
			mainpkg.WithInit("rbac", func(ctx context.Context) error {
				<-lock
				return st.ExpireMFAEvents(ctx)
			}),
		},
	}, nil
}

// intlSvc registers using the `RegisterInternal` method.
type intlSvc struct {
	mainpkg.GWGRPC
	intlRegr
}

type intlRegr interface {
	RegisterInternal(*grpc.Server) error
	RegisterGatewayInternal(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}

func (i intlSvc) RegisterSvc(s *grpc.Server) error { i.intlRegr.RegisterInternal(s); return nil }
func (i intlSvc) RegisterGateway(
	c context.Context, m *runtime.ServeMux, a string, o []grpc.DialOption,
) error {
	return i.RegisterGatewayInternal(c, m, a, o)
}

func initHydra(config *viper.Viper) (*hydra.Service, error) {
	cl := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   30 * time.Second,
	}
	aud := config.GetString("hydra.audience")
	hydraAdmin := config.GetString("hydra.adminURL")
	switch "" {
	case aud, hydraAdmin:
		return nil, fmt.Errorf("hydra configuration missing entries aud(%q), url(%q)", aud, hydraAdmin)
	}
	hs, err := hydra.NewService(cl, hydraAdmin, hydra.WithOptional())
	if err != nil {
		return nil, err
	}
	return hs, nil
}

// NewDBFromConfig build database connection from config file.
func newDBFromConfig(config *viper.Viper) (*postgres.Storage, error) {
	cf := func(c string) string { return config.GetString("database." + c) }
	ci := func(c string) string { return strconv.Itoa(config.GetInt("database." + c)) }
	dbParams := " " + "user=" + cf("user")
	dbParams += " " + "host=" + cf("host")
	dbParams += " " + "port=" + cf("port")
	dbParams += " " + "dbname=" + cf("dbname")
	if password := cf("password"); password != "" {
		dbParams += " " + "password=" + password
	}
	dbParams += " " + "sslmode=" + cf("sslMode")
	dbParams += " " + "connect_timeout=" + ci("connectionTimeout")
	dbParams += " " + "statement_timeout=" + ci("statementTimeout")
	dbParams += " " + "idle_in_transaction_session_timeout=" + ci("idleTransacionTimeout")
	return postgres.NewStorage(dbParams)
}

func initMailer(config *viper.Viper) (email.Mailer, error) {
	switch {
	case config.GetBool("mailer.disable"):
		return email.NoopSender{}, nil
	case config.GetBool("mailer.mock"):
		return email.NewMock(
			email.MailerConfig{
				IDPURL:        config.GetString("idp.url"),
				SignupURL:     config.GetString("mailer.signupURL"),
				ECRedirectURL: config.GetString("mailer.emailConfirmURL"),
				SiteURL:       config.GetString("site.url"),
				Subjects: email.Subjects{
					ForgotPW:       config.GetString("mailer.subjects.resetPassword"),
					ConfirmEmail:   config.GetString("mailer.subjects.confirmEmail"),
					EmailMFA:       config.GetString("mailer.subjects.emailMFA"),
					UserInvite:     config.GetString("mailer.subjects.userInvite"),
					InviteApproved: config.GetString("mailer.subjects.emailApproved"),
					UserDisable:    config.GetString("mailer.subjects.userDisable"),
					UserEnable:     config.GetString("mailer.subjects.userEnable"),
				},
			},
			afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), config.GetString("mailer.assetdir"))),
		)
	default:
		return email.New(
			email.MailerConfig{
				Server:        config.GetString("smtp.host"),
				Port:          config.GetInt("smtp.port"),
				Username:      config.GetString("smtp.username"),
				Password:      config.GetString("smtp.password"),
				FromMail:      config.GetString("smtp.fromAddr"),
				FromName:      config.GetString("smtp.fromName"),
				MainURL:       config.GetString("smtp.baseURL"),
				IDPURL:        config.GetString("idp.url"),
				SignupURL:     config.GetString("mailer.signupURL"),
				ECRedirectURL: config.GetString("mailer.emailConfirmURL"),
				SiteURL:       config.GetString("site.url"),
				CACert:        config.GetString("smtp.cacert"),
				CACommonName:  config.GetString("smtp.certcn"),
				Subjects: email.Subjects{
					ForgotPW:       config.GetString("mailer.subjects.resetPassword"),
					ConfirmEmail:   config.GetString("mailer.subjects.confirmEmail"),
					EmailMFA:       config.GetString("mailer.subjects.emailMFA"),
					UserInvite:     config.GetString("mailer.subjects.userInvite"),
					InviteApproved: config.GetString("mailer.subjects.emailApproved"),
					UserDisable:    config.GetString("mailer.subjects.userDisable"),
					UserEnable:     config.GetString("mailer.subjects.userEnable"),
				},
			},
			afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), config.GetString("mailer.assetdir"))),
		)
	}
}
