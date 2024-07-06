package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	atc "brank.as/petnet/profile/core/apitransactiontype"
	brc "brank.as/petnet/profile/core/branch"
	cpnrcl "brank.as/petnet/profile/core/cicopartnerlist"
	emlc "brank.as/petnet/profile/core/email"
	evc "brank.as/petnet/profile/core/event"
	fec "brank.as/petnet/profile/core/fees"
	fic "brank.as/petnet/profile/core/file"
	mc "brank.as/petnet/profile/core/mfa"
	pnrc "brank.as/petnet/profile/core/partner"
	rcc "brank.as/petnet/profile/core/partnercommission"
	pnrcl "brank.as/petnet/profile/core/partnerlist"
	pfc "brank.as/petnet/profile/core/profile"
	rbsuc "brank.as/petnet/profile/core/rbsignup"
	rsh "brank.as/petnet/profile/core/revenuesharing"
	rsp "brank.as/petnet/profile/core/revenuesharingreport"
	ric "brank.as/petnet/profile/core/riskassesment"
	sec "brank.as/petnet/profile/core/session"
	sic "brank.as/petnet/profile/core/signup"
	"brank.as/petnet/profile/integrations/email"
	"brank.as/petnet/profile/partners"
	"brank.as/petnet/profile/permission"
	ats "brank.as/petnet/profile/services/apitransactiontype"
	brs "brank.as/petnet/profile/services/branch"
	cpnrsl "brank.as/petnet/profile/services/cicopartnerlist"
	ems "brank.as/petnet/profile/services/email"
	evs "brank.as/petnet/profile/services/event"
	fes "brank.as/petnet/profile/services/fees"
	fis "brank.as/petnet/profile/services/file"
	ms "brank.as/petnet/profile/services/mfa"
	ops "brank.as/petnet/profile/services/orgprofile"
	pnrs "brank.as/petnet/profile/services/partner"
	rcs "brank.as/petnet/profile/services/partnercommission"
	pnrsl "brank.as/petnet/profile/services/partnerlist"
	rbses "brank.as/petnet/profile/services/rbsession"
	rbsu "brank.as/petnet/profile/services/rbsignup"
	rsc "brank.as/petnet/profile/services/revenuesharing"
	rsr "brank.as/petnet/profile/services/revenuesharingreport"
	ris "brank.as/petnet/profile/services/riskassesment"
	ss "brank.as/petnet/profile/services/service"
	ses "brank.as/petnet/profile/services/session"
	sis "brank.as/petnet/profile/services/signup"
	ups "brank.as/petnet/profile/services/userprofile"
	"brank.as/petnet/profile/storage/postgres"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/serviceutil/mainpkg"
	"brank.as/petnet/serviceutil/middleware"
	"brank.as/petnet/svcutil/metrics"
	"brank.as/petnet/svcutil/mw/meta"
	"brank.as/petnet/svcutil/mw/meta/md"
	rbapb "brank.as/rbac/gunk/v1/authenticate"
	rbipb "brank.as/rbac/gunk/v1/invite"
	rbmpb "brank.as/rbac/gunk/v1/mfa"
	rbopb "brank.as/rbac/gunk/v1/organization"
	rbppb "brank.as/rbac/gunk/v1/permissions"
	rbupb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/svcutil/otelb"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	version = "devel"
	svcName = "profile"
)

func main() {
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	log.Println("dialing trace collector...")
	shutdown := otelb.InitOTELProvider(
		context.Background(),
		svcName,
		config.GetString("trace.collectorHost"),
	)
	defer shutdown()
	log := logging.NewLogger(config).WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})
	log.Info("starting profile service")
	met, err := metrics.NewInfluxDBClient(config)
	if err != nil {
		log.Fatal(err)
	}
	defer met.Close()
	go met.ErrorsFunc(func(e error) { log.WithError(err).Error("influxdb write") })
	store, err := newDBFromConfig(config)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure storage")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	svc, err := setupServices(ctx, log, config, store)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure grpc services")
	}
	cancel()
	grpcEx, gwEx, err := setupGRPCServer(config, log, svc.external, store, met)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to start grpc server")
	}
	grpcInt, gwInt, err := setupGRPCInternal(config, log, svc.internal, met)
	if err != nil {
		log.Fatal(err)
	}
	// serve internal services on admin port
	adm, err := mainpkg.Setup(config, log,
		mainpkg.WithPort(config.GetInt("server.adminPort")),
		mainpkg.WithGRPCServer(grpcInt),
		mainpkg.WithVersion(svcName, version),
		mainpkg.AddRegisterGatewayFunc(gwInt...),
	)
	if err != nil {
		log.Fatal(err)
	}
	svr, err := mainpkg.Setup(config, log,
		mainpkg.WithGRPCServer(grpcEx),
		mainpkg.WithVersion(svcName, version),
		mainpkg.AddRegisterGatewayFunc(gwEx...),
		mainpkg.AdditionalServers(adm),
	)
	if err != nil {
		log.Fatal(err)
	}
	svr.Run()
	log.Info("server exited")
}

type Service interface {
	Register(*grpc.Server)
	RegisterGateway(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
type Svcs struct {
	external []Service
	internal []Service
}

func setupServices(ctx context.Context, log *logrus.Entry, c *viper.Viper, store *postgres.Storage) (*Svcs, error) {
	if err := store.RunMigration(c.GetString("database.migrationDir")); err != nil {
		return nil, err
	}
	var ucl rbupb.UserServiceClient = ops.Mock{}
	var acl rbupb.UserAuthServiceClient = ops.Mock{}
	var pebcl rbppb.PermissionServiceClient = ops.Mock{}
	var prcl rbppb.ProductServiceClient = ops.Mock{}
	var sicl rbupb.SignupClient = ops.Mock{}
	var secl rbapb.SessionServiceClient = ops.Mock{}
	var mcl rbmpb.MFAServiceClient
	var icl rbipb.InviteServiceClient
	var ocl rbopb.OrganizationServiceClient
	if !c.GetBool("usermgm.mock") {
		log.WithField("host", c.GetString("usermgm.admin")).Println("dialing identity admin")
		admConn, err := grpc.Dial(
			c.GetString("usermgm.admin"),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithChainUnaryInterceptor(
				otelgrpc.UnaryClientInterceptor(),
				injectMDInterceptor(c),
			),
		)
		if err != nil {
			logging.WithError(err, log).Fatal("unable to connect to identity admin")
			return nil, err
		}
		log.WithField("host", c.GetString("usermgm.api")).Println("dialing identity-provider with bootstrap user")
		apiBootstrapConn, err := grpc.Dial(
			c.GetString("usermgm.api"),
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithChainUnaryInterceptor(
				otelgrpc.UnaryClientInterceptor(),
				injectMDInterceptor(c),
			),
		)
		if err != nil {
			logging.WithError(err, log).Fatal("unable to connect to identity api with bootstrap user")
			return nil, err
		}
		sicl = rbupb.NewSignupClient(apiBootstrapConn)
		ucl = rbupb.NewUserServiceClient(apiBootstrapConn)
		acl = rbupb.NewUserAuthServiceClient(admConn)
		prcl = rbppb.NewProductServiceClient(apiBootstrapConn)
		secl = rbapb.NewSessionServiceClient(admConn)
		mcl = rbmpb.NewMFAServiceClient(apiBootstrapConn)
		icl = rbipb.NewInviteServiceClient(apiBootstrapConn)
		ocl = rbopb.NewOrganizationServiceClient(apiBootstrapConn)
		prcl = rbppb.NewProductServiceClient(apiBootstrapConn)
		pebcl = rbppb.NewPermissionServiceClient(apiBootstrapConn)
	}
	mailer := email.New(
		c.GetString("smtp.host"),
		c.GetInt("smtp.port"),
		c.GetString("smtp.username"),
		c.GetString("smtp.password"),
		c.GetString("smtp.fromAddr"),
		c.GetString("smtp.fromName"),
		c.GetString("cms.url"),
	)
	su := rbsuc.New(sicl)
	st := pfc.New(store)
	op := ops.New(store, ucl, st)
	up := ups.New(st)
	tt := ats.New(atc.New(store))
	rc := ric.New(store)
	rs := ris.New(rc)
	br := brs.New(brc.New(store))
	emli := emlc.New(mailer)
	em := ems.New(emli)
	fi := fis.New(fic.New(store))
	fe := fes.New(fec.New(store))
	rcsc := rcs.New(rcc.New(store))
	rsrp := rsr.New(rsp.New(store))
	rvsh := rsc.New(rsh.New(store))
	sv := pnrs.New(pnrc.New(store))
	svl := pnrsl.New(pnrcl.New(store))
	svlc := cpnrsl.New(cpnrcl.New(store))
	rbse := rbses.New(st, acl, ucl, secl)
	rbsus := rbsu.New(su)
	se := ses.New(sec.New(store))
	ev := evs.New(evc.New(store))
	si := sis.New(sic.New(store, sicl, ocl, icl, c.GetBool("local.disableLoginMFA")))
	m := ms.New(mc.New(store, mcl))
	s := ss.New(store, emli)
	if !c.GetBool("local.disablePermissionBootstrap") {
		if err := permission.BootstrapAdminPermissions(ctx, log, pebcl, prcl, store); err != nil {
			logging.WithError(err, log).Fatal("unable to bootstrap petnet admin permissions")
			return nil, err
		}
	}
	if err := partners.BootstrapAdminPartners(ctx, log, store); err != nil {
		logging.WithError(err, log).Fatal("unable to bootstrap petnet partner")
		return nil, err
	}
	if err := partners.BootstrapAdminCicoPartners(ctx, log, store); err != nil {
		logging.WithError(err, log).Fatal("unable to bootstrap petnet cico partner")
		return nil, err
	}
	if err := partners.BootstrapAdminRTAPartners(ctx, log, store); err != nil {
		logging.WithError(err, log).Fatal("unable to bootstrap petnet RTA partner")
		return nil, err
	}
	return &Svcs{
		external: []Service{},
		internal: []Service{op, up, rbse, se, br, em, fi, fe, sv, ev, si, m, rbsus, rs, svl, s, tt, rcsc, rvsh, rsrp, svlc},
	}, nil
}

func setupGRPCServer(config *viper.Viper, log *logrus.Entry, s []Service, st *postgres.Storage,
	met *metrics.Influxdb,
) (*grpc.Server, []mainpkg.RegisterGatewayFunc, error) {
	m, err := meta.NewMetadata(log,
		meta.MetaFunc(func(ctx context.Context) (context.Context, error) {
			nmd := metadata.MD{}
			md := metadata.MD(metautils.ExtractIncoming(ctx))
			for k, v := range md {
				nmd.Append(k, v...)
			}
			return metautils.NiceMD(nmd).ToIncoming(ctx), nil
		}),
		md.LoadOrgProfile(st),
	)
	if err != nil {
		return nil, nil, err
	}
	mw := middleware.New(
		config.GetString("runtime.environment"), log,
		[]grpc_logrus.Option{}, true,
		met.UnaryServerInterceptor("", nil),
		m.UnaryServerInterceptor(),
	)
	srv := grpc.NewServer(grpc.UnaryInterceptor(mw),
		grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: true}))
	reflection.Register(srv)
	gw := []mainpkg.RegisterGatewayFunc{}
	for _, svc := range s {
		svc.Register(srv)
		gw = append(gw, svc.RegisterGateway)
	}
	return srv, gw, nil
}

func setupGRPCInternal(c *viper.Viper, log *logrus.Entry, s []Service,
	met *metrics.Influxdb,
) (*grpc.Server, []mainpkg.RegisterGatewayFunc, error) {
	hs, err := initHydra(c)
	if err != nil {
		log.Fatal(err)
	}
	iMD, err := meta.NewMetadata(log, hs.GRPC())
	if err != nil {
		log.Fatal(err)
	}
	mw := middleware.New(c.GetString("runtime.environment"), log,
		[]grpc_logrus.Option{}, false,
		otelgrpc.UnaryServerInterceptor(),
		met.UnaryServerInterceptor("", nil),
		iMD.UnaryServerInterceptor(),
	)
	srv := grpc.NewServer(grpc.UnaryInterceptor(mw),
		grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: false}))
	reflection.Register(srv)
	gw := []mainpkg.RegisterGatewayFunc{}
	for _, svc := range s {
		svc.Register(srv)
		gw = append(gw, svc.RegisterGateway)
	}
	return srv, gw, nil
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
	db, err := postgres.NewStorage(dbParams)
	if err != nil {
		return nil, err
	}
	return db, db.RunMigration(cf("migrationDir"))
}

func injectMDInterceptor(c *viper.Viper) grpc.UnaryClientInterceptor {
	cred := clientcredentials.Config{
		ClientID:     c.GetString("auth.clientID"),
		ClientSecret: c.GetString("auth.clientSecret"),
		TokenURL:     c.GetString("auth.url") + "/oauth2/token",
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

func initHydra(config *viper.Viper) (*hydra.Service, error) {
	cl := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   1 * time.Second,
	}
	aud := config.GetString("hydra.audience")
	hAdm := config.GetString("hydra.adminURL")
	switch "" {
	case aud, hAdm:
		return nil, fmt.Errorf("hydra configuration missing entries")
	}
	return hydra.NewService(cl, hAdm, hydra.WithOptional())
}
