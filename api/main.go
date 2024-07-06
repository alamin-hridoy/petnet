package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/pariz/gountries"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	// common serviceutil
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/serviceutil/mainpkg"
	"brank.as/petnet/serviceutil/middleware"

	// svcutil
	"brank.as/petnet/svcutil/metrics"
	"brank.as/petnet/svcutil/mw"
	"brank.as/petnet/svcutil/mw/meta"
	"brank.as/rbac/svcutil/otelb"

	// middleware
	phmw "brank.as/petnet/api/perahub-middleware"

	// integrations
	bpi "brank.as/petnet/api/integration/bills-payment"
	micins_int "brank.as/petnet/api/integration/microinsurance"
	"brank.as/petnet/api/integration/perahub"
	rtai "brank.as/petnet/api/integration/remittoaccount"
	revcom_int "brank.as/petnet/api/integration/revenue-commission"
	"brank.as/petnet/api/storage/postgres"

	// core logic
	"brank.as/petnet/api/core/auth"
	bpac "brank.as/petnet/api/core/bills-payment"
	bpacBc "brank.as/petnet/api/core/bills-payment/bayadcenter"
	bpacEp "brank.as/petnet/api/core/bills-payment/ecpay"
	bpacMp "brank.as/petnet/api/core/bills-payment/multipay"
	cicoc "brank.as/petnet/api/core/cashincashout"
	fc "brank.as/petnet/api/core/fee"
	miCore "brank.as/petnet/api/core/microinsurance"
	pc "brank.as/petnet/api/core/partner"
	qc "brank.as/petnet/api/core/quote"
	"brank.as/petnet/api/core/remit"
	aya "brank.as/petnet/api/core/remit/ayannah"
	bp "brank.as/petnet/api/core/remit/bpi"
	ceb "brank.as/petnet/api/core/remit/cebuana"
	cebint "brank.as/petnet/api/core/remit/cebuanaint"
	ic "brank.as/petnet/api/core/remit/instacash"
	ie "brank.as/petnet/api/core/remit/intelexpress"
	ir "brank.as/petnet/api/core/remit/iremit"
	jr "brank.as/petnet/api/core/remit/japanremit"
	mb "brank.as/petnet/api/core/remit/metrobank"
	ph "brank.as/petnet/api/core/remit/perahubremit"
	rm "brank.as/petnet/api/core/remit/remitly"
	"brank.as/petnet/api/core/remit/ria"
	tf "brank.as/petnet/api/core/remit/transfast"
	unt "brank.as/petnet/api/core/remit/uniteller"
	usc "brank.as/petnet/api/core/remit/ussc"
	ws "brank.as/petnet/api/core/remit/wise"
	"brank.as/petnet/api/core/remit/wu"
	remitc "brank.as/petnet/api/core/remittance"
	rtac "brank.as/petnet/api/core/remittoaccount"
	rtaBpi "brank.as/petnet/api/core/remittoaccount/rtabpi"
	rtaMb "brank.as/petnet/api/core/remittoaccount/rtamb"
	rtaUb "brank.as/petnet/api/core/remittoaccount/rtaub"
	"brank.as/petnet/api/core/static"
	uc "brank.as/petnet/api/core/user"
	apiutil "brank.as/petnet/api/util"

	// grpc services
	authSvc "brank.as/petnet/api/services/auth"
	bpas "brank.as/petnet/api/services/bills-payment"
	cicos "brank.as/petnet/api/services/cashincashout"
	fSvc "brank.as/petnet/api/services/fee"
	"brank.as/petnet/api/services/microinsurance"
	rpSvc "brank.as/petnet/api/services/partner"
	qteSvc "brank.as/petnet/api/services/quote"
	remits "brank.as/petnet/api/services/remittance"
	rtas "brank.as/petnet/api/services/remittoaccount"
	revcom "brank.as/petnet/api/services/revenue-commission"
	"brank.as/petnet/api/services/terminal"
	usrSvc "brank.as/petnet/api/services/user"

	// proto
	pfppb "brank.as/petnet/gunk/dsa/v2/partner"
	ptnrLst "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	rsr "brank.as/petnet/gunk/dsa/v2/revenuesharingreport"
	pfSvc "brank.as/petnet/gunk/dsa/v2/service"
	trxtp "brank.as/petnet/gunk/dsa/v2/transactiontype"
	osapb "brank.as/rbac/gunk/v1/oauth2"
	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

const (
	svcName = "drp"
	version = "development"
)

func main() {
	c := viper.NewWithOptions(
		viper.EnvKeyReplacer(
			strings.NewReplacer(".", "_"),
		),
	)
	c.SetConfigFile("env/config")
	c.SetConfigType("ini")
	c.AutomaticEnv()
	if err := c.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}

	log.Println("dialing trace collector...")
	shutdown := otelb.InitOTELProvider(
		context.Background(),
		svcName,
		c.GetString("trace.collectorHost"),
	)
	defer shutdown()

	log := logging.NewLogger(c).WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})

	switch c.GetString("runtime.loglevel") {
	case "trace":
		log.Logger.SetLevel(logrus.TraceLevel)
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	default:
		log.Logger.SetLevel(logrus.InfoLevel)
	}
	log.WithField("log level", log.Logger.Level).Info("starting service")

	u := util{}
	var err error
	u.met, err = metrics.NewInfluxDBClient(c)
	if err != nil {
		log.Fatal(err)
	}
	u.met.DefaultTags(map[string]string{
		"env": c.GetString("runtime.environment"),
	})
	go u.met.ErrorsFunc(func(e error) { logging.WithError(e, log).Warn("influxdb") })

	u.cs = newConns(log, c)
	defer u.cs.close()
	u.st = newDB(log, c)
	u.hs = newHydra(log, c)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	svcs, err := u.newSvcs(ctx, log, c)
	if err != nil {
		log.Fatal(err)
	}
	intSvr := u.newIntlServer(log, c, svcs)
	extSvr := u.newExtlServer(log, c, intSvr, svcs)

	extSvr.Run()
	log.Info("exiting")
}

type Services struct {
	Internal []mainpkg.GWGRPC
	External []mainpkg.GWGRPC
	Option   []mainpkg.Option
}

func newHydra(log *logrus.Entry, config *viper.Viper) *hydra.Service {
	cl := &http.Client{
		Transport: http.DefaultTransport,
		Timeout:   1 * time.Second,
	}
	aud := config.GetString("hydra.audience")
	hAdm := config.GetString("hydra.adminURL")
	switch "" {
	case aud, hAdm:
		log.Fatal("hydra configuration missing entries")
	}
	hs, err := hydra.NewService(cl, hAdm, hydra.WithOptional())
	if err != nil {
		log.Fatal(err)
	}
	return hs
}

func newDB(log *logrus.Entry, c *viper.Viper) *postgres.Storage {
	st, err := postgres.New(c)
	if err != nil {
		log.Fatal(err)
	}
	if err := st.RunMigration(c.GetString("database.migrationDir")); err != nil {
		log.Fatal(err)
	}
	return st
}

func LoggerInterceptor(log logrus.FieldLogger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo,
		h grpc.UnaryHandler,
	) (interface{}, error) {
		ctx = logging.WithLogger(ctx, log)
		res, err := h(ctx, req)
		logging.FromContext(ctx).WithField("response", res).
			WithField("resp_error", err).Trace("handled")
		return res, err
	}
}

func UserDetailsInternal() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if info.FullMethod == "/authenticate.SessionService/Login" {
			return handler(ctx, req)
		}
		ot, err := strconv.Atoi(hydra.OrgType(ctx))
		if err != nil {
			return nil, fmt.Errorf("org type is empty")
		}
		md := metautils.ExtractIncoming(ctx)
		ctx = md.Set(phmw.OrgType, ppb.OrgType_name[int32(ot)]).ToIncoming(ctx)
		switch ot {
		case int(ppb.OrgType_PetNet), int(ppb.OrgType_DSA):
			ctx = md.Set(phmw.UserID, md.Get("x-forward-clientid")).Set(phmw.DSAOrgID, md.Get("x-forward-dsaorgid")).ToIncoming(ctx)
		}

		return handler(ctx, req)
	}
}

type conns struct {
	idInt *grpc.ClientConn
	idExt *grpc.ClientConn
	pfInt *grpc.ClientConn
}

func newConns(log *logrus.Entry, c *viper.Viper) *conns {
	log.WithField("host", c.GetString("identity.internal")).Info("dialing identity internal")
	idInt, err := grpc.Dial(
		c.GetString("identity.internal"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("host", c.GetString("profile.internal")).Info("dialing profile internal")
	pfInt, err := grpc.Dial(
		c.GetString("profile.internal"),
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatal("unable to connect to profile internal")
	}

	return &conns{
		idInt: idInt,
		pfInt: pfInt,
	}
}

func (cs *conns) close() {
	cs.idInt.Close()
	cs.idExt.Close()
	cs.pfInt.Close()
}

func (u *util) newIntlServer(log *logrus.Entry, c *viper.Viper, svcs *Services) *mainpkg.Server {
	// internal API
	iMD, err := meta.NewMetadata(log, u.hs.GRPC(), u.hs.InternalGRPC())
	if err != nil {
		log.Fatal(err)
	}

	grpcMeasurement := c.GetString("influxdb.grpcMeasure")
	if grpcMeasurement == "" {
		grpcMeasurement = "conex_grpc_latency"
	}
	iMW := middleware.New(c.GetString("runtime.environment"), log, nil, false,
		otelgrpc.UnaryServerInterceptor(),
		LoggerInterceptor(log),
		u.met.UnaryServerInterceptor("internal_"+grpcMeasurement, nil),
		iMD.UnaryServerInterceptor(),
		UserDetailsInternal(),
	)
	svr := grpc.NewServer(grpc.UnaryInterceptor(iMW))

	intSvr, err := mainpkg.Setup(c, log,
		mainpkg.WithGRPCServer(svr),
		mainpkg.WithDualService(mainpkg.Internal, svcs.Internal...),
		mainpkg.WithPort(c.GetInt("internal.port")),
	)
	if err != nil {
		log.Fatal(err)
	}
	return intSvr
}

func (u *util) newExtlServer(log *logrus.Entry, c *viper.Viper, intSvr *mainpkg.Server, svcs *Services) *mainpkg.Server {
	md, err := meta.NewMetadata(log,
		phmw.Reset(),            // Clear metadata
		u.hs.GRPC(), phmw.New(), // User auth
		// service account auth
		phmw.NewServiceAccount(sapb.NewValidationServiceClient(u.cs.idInt), trxtp.NewTransactionTypeServiceClient(u.cs.pfInt), ppb.NewOrgProfileServiceClient(u.cs.pfInt), log),
		phmw.ConfirmDSA(c.GetString("runtime.environment")), // Enforcement
		meta.MetaFunc(func(ctx context.Context) (context.Context, error) {
			// annotate metrics
			metrics.SetTag(ctx, "dsa_id", phmw.GetDSA(ctx))
			return ctx, nil
		}),
	)
	if err != nil {
		log.Fatal(err)
	}

	grpcMeasurement := c.GetString("influxdb.grpcMeasure")
	if grpcMeasurement == "" {
		grpcMeasurement = "conex_grpc_latency"
	}
	serviceClient := pfSvc.NewServiceServiceClient(u.cs.pfInt)
	ptnrClient := ptnrLst.NewPartnerListServiceClient(u.cs.pfInt)
	mwr := middleware.New(c.GetString("runtime.environment"), log, nil, true,
		otelgrpc.UnaryServerInterceptor(),
		u.met.UnaryServerInterceptor(grpcMeasurement, nil),
		LoggerInterceptor(log),
		md.UnaryServerInterceptor(),
		mw.ValidateAccess(serviceClient, ptnrClient, u.trmSvc),
	)
	svr := grpc.NewServer(grpc.UnaryInterceptor(mwr))

	eh := apiutil.NewErrorHandler(log)
	extSvr, err := mainpkg.Setup(c, log,
		mainpkg.WithGRPCServer(svr),
		mainpkg.WithDualService(mainpkg.External, svcs.External...),
		mainpkg.AdditionalServers(intSvr),
		mainpkg.WithGatewayProtoError(eh.HTTPErrorHandler),
		mainpkg.OptionList(svcs.Option),
	)
	if err != nil {
		log.Fatal(err)
	}

	return extSvr
}

// newSvcs configure services.
func (u *util) newSvcs(ctx context.Context, log *logrus.Entry, c *viper.Viper) (*Services, error) {
	st := u.st

	sched := "0 0 1 * *"
	newSched := "0 0 * * *" // every day at 12.00AM
	cl := &http.Client{
		Timeout:   c.GetDuration("perahub.timeout"),
		Transport: u.met.NewTransport("http_perahub_gateway", mw.JSONType(nil)),
	}

	env := c.GetString("runtime.environment")
	var phKey, nonexAPIKey string
	if env == "live" {
		phKey = c.GetString("perahub.liveclientkey")
		nonexAPIKey = c.GetString("perahub.nonexApiKey")
	} else {
		phKey = c.GetString("perahub.sandboxclientkey")
		nonexAPIKey = c.GetString("perahub.defaultAPIKey")
	}

	phintg, err := perahub.New(cl,
		c.GetString("runtime.environment"),
		c.GetString("perahub.baseurl"),
		c.GetString("perahub.nonexurl"),
		c.GetString("perahub.billerurl"),
		c.GetString("perahub.transacturl"),
		c.GetString("perahub.partnerid"),
		c.GetString("perahub.billsurl"),
		phKey,
		nonexAPIKey,
		c.GetString("perahub.serverip"),
		c.GetString("perahub.signingkey"),
		map[string]perahub.OAuthCreds{
			static.WISECode: {
				ClientID:     c.GetString("perahub.wiseClientID"),
				ClientSecret: c.GetString("perahub.wiseClientSecret"),
			},
		},
		perahub.WithLogger(log),
		perahub.WithCiCoURL(c.GetString("perahub.cicourl")),
		perahub.WithPHRemittanceURL(c.GetString("perahub.remittanceUrl")),
		perahub.WithPerahubDefaultAPIKey(c.GetString("perahub.defaultAPIKey")),
	)
	if err != nil {
		return nil, err
	}
	var phm bool
	if c.GetBool("perahub.httpmock") {
		phintg = phintg.SetMock(perahub.NewHTTPMock(st))
		phm = true
	}

	stccore := static.New(phintg, st)

	// remitters
	rmts := []remit.Remitter{
		wu.New(st, phintg, stccore, phm),
		ir.New(st, phintg, stccore),
		tf.New(st, phintg, stccore),
		ria.New(st, phintg, stccore),
		rm.New(st, phintg, stccore),
		mb.New(st, phintg, stccore),
		bp.New(st, phintg, stccore),
		usc.New(st, phintg, stccore),
		ic.New(st, phintg, stccore),
		jr.New(st, phintg, stccore),
		ws.New(st, phintg, stccore),
		unt.New(st, phintg, stccore),
		ceb.New(st, phintg, stccore),
		cebint.New(st, phintg, stccore),
		aya.New(st, phintg, stccore),
		ie.New(st, phintg, stccore),
		ph.New(st, phintg, stccore),
	}

	// core logic
	rmtcore, err := remit.New(st, phintg, stccore, rmts)
	if err != nil {
		return nil, err
	}

	// validators
	q := gountries.New()
	trmval := terminal.NewValidators(q)
	trmsvc, err := terminal.New(rmtcore, stccore, trmval)
	if err != nil {
		return nil, err
	}
	u.trmSvc = trmsvc

	// start remit to account
	rtaUrl, err := url.Parse(c.GetString("perahub.rtaurl"))
	if err != nil {
		return nil, err
	}

	if c.GetBool("perahub.httpmock") {
		phintg = phintg.SetMock(perahub.NewHTTPMock(st))
	}

	rtaClnt := rtai.NewRTAClient(phintg, rtaUrl)
	rtaPtnrs := []rtac.Remitter{
		rtaBpi.New(st, rtaClnt),
		rtaUb.New(st, rtaClnt),
		rtaMb.New(st, rtaClnt),
	}
	rtaval := rtas.NewValidators(q)
	rtaCore, err := rtac.New(rtaPtnrs, rtaClnt)
	if err != nil {
		return nil, err
	}
	rtaSvc, err := rtas.New(rtaCore, rtaval)
	if err != nil {
		return nil, err
	}
	// end remit to account

	// start bills payment
	bpUrl, err := url.Parse(c.GetString("perahub.billsurl"))
	if err != nil {
		return nil, err
	}
	if c.GetBool("perahub.httpmock") {
		phintg = phintg.SetMock(perahub.NewHTTPMock(st))
	}
	bpClnt := bpi.NewBillsClient(phintg, bpUrl)
	bpPtnrs := []bpac.Biller{
		bpacBc.New(st, bpClnt),
		bpacEp.New(st, bpClnt),
		bpacMp.New(st, bpClnt),
	}
	bpval := bpas.NewValidators(q)
	bpCore, err := bpac.New(bpPtnrs, bpClnt, st)
	if err != nil {
		return nil, err
	}
	bpSvc, err := bpas.New(bpCore, bpval)
	if err != nil {
		return nil, err
	}
	// end billspayment

	ptnrval := rpSvc.NewValidators()
	ptnrsvc := rpSvc.New(stccore, pc.New(st, phintg), pfppb.NewPartnerServiceClient(u.cs.pfInt), pfSvc.NewServiceServiceClient(u.cs.pfInt), ptnrLst.NewPartnerListServiceClient(u.cs.pfInt), ptnrval)

	feeval := fSvc.NewValidators()
	feesvc := fSvc.New(stccore, fc.New(st, phintg), feeval)

	usrval := usrSvc.NewValidators()
	usrsvc := usrSvc.New(uc.New(st, phintg), usrval)

	qteval := qteSvc.NewValidators()
	qtesvc := qteSvc.New(qc.New(st, phintg), qteval)

	cicovc := cicos.New(cicoc.New(phintg, st))
	remitsvc := remits.New(remitc.New(st, phintg))

	// internal services
	if c.GetBool("perahub.httpmock") {
		phintg = phintg.SetMock(perahub.NewHTTPMock(st))
	}

	asvc := authSvc.New(auth.New(phintg, st), osapb.NewAuthClientServiceClient(u.cs.idInt))

	revComBaseUrl, err := url.Parse(c.GetString("perahub.revcommurl"))
	if err != nil {
		return nil, err
	}

	revComClient := revcom_int.NewRevCommClient(phintg, revComBaseUrl)
	revComSvc := revcom.NewRevenueCommissionService(
		revcom.WithDSAStore(revComClient),
		revcom.WithCommissionFeeStore(revComClient),
		revcom.WithDSACommissionStore(revComClient),
		revcom.WithStorage(st),
		revcom.WithRevenueSharingReport(rsr.NewRevenueSharingReportServiceClient(u.cs.pfInt)),
		revcom.WithProfileDSA(ppb.NewOrgProfileServiceClient(u.cs.pfInt)),
	)

	micInsBaseUrl, err := url.Parse(c.GetString("perahub.micInsUrl"))
	if err != nil {
		return nil, err
	}

	miClient := micins_int.NewMicroInsuranceClient(phintg, micInsBaseUrl)
	miSvc := microinsurance.NewMicroInsuranceSvc(miCore.NewMicroInsuranceCoreSvc(st, miClient))

	return &Services{
		External: []mainpkg.GWGRPC{ptnrsvc, trmsvc, feesvc, usrsvc, qtesvc, remitsvc, cicovc, rtaSvc, miSvc, bpSvc, revComSvc},
		Internal: []mainpkg.GWGRPC{asvc, trmsvc, remitsvc, revComSvc, cicovc, rtaSvc, miSvc, bpSvc},
		Option: []mainpkg.Option{
			mainpkg.WithCron("Create Trannsaction Report", mainpkg.NewCrontab(sched), revComSvc.SyncTransactionReport),
			mainpkg.WithCron("update remco id", mainpkg.NewCrontab(newSched), ptnrsvc.UpdateRemcoId),
		},
	}, nil
}

type util struct {
	hs     *hydra.Service
	st     *postgres.Storage
	cs     *conns
	met    *metrics.Influxdb
	trmSvc *terminal.Svc
}
