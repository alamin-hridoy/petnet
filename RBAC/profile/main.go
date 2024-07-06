package main

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"

	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"go.opencensus.io/plugin/ocgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	"brank.as/rbac/serviceutil/logging"

	"brank.as/rbac/profile/storage/postgres"
	"brank.as/rbac/svcutil/mainpkg"

	"brank.as/rbac/profile/services/session"
	"brank.as/rbac/profile/services/user"

	ipupb "brank.as/rbac/gunk/v1/user"
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
	config := viper.NewWithOptions(
		viper.EnvKeyReplacer(strings.NewReplacer(".", "_")),
	)
	config.SetConfigFile("env/config")
	config.SetConfigType("ini")
	config.AutomaticEnv()
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error loading configuration: %v", err)
	}
	switch config.GetString("RUNTIME_LOGLEVEL") {
	case "debug":
		log.Logger.SetLevel(logrus.DebugLevel)
	case "trace":
		log.Logger.SetLevel(logrus.TraceLevel)
	}
	log.Info("starting usermgm service", log.Level)

	log.Println("dialing identity-provider api...")
	u := config.GetString("usermgm.api")
	opts := getGRPCOpts(config)
	conn, err := grpc.Dial(u, opts...)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to connect to identity-provider")
	}
	defer conn.Close()

	init := config.GetBool("mainpkg.init")
	cleanup := config.GetBool("mainpkg.cleanup")
	// leader := config.GetBool("mainpkg.leader")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	svc, err := setupServices(ctx, log, config, conn)
	if err != nil {
		logging.WithError(err, log).Fatal("unable to configure grpc services")
	}
	cancel()

	// serve internal services on admin port
	intl, err := mainpkg.Setup(config, log,
		mainpkg.WithPort(config.GetInt("server.adminPort")),
		mainpkg.WithVersion(svcName, version),
		mainpkg.WithDualService(mainpkg.Internal, svc.internal...),
		mainpkg.WithServerOpts(grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: true})),
		mainpkg.WithInit("testinit", func(context.Context) error {
			if init {
				fmt.Println("initializing 1")
				time.Sleep(3 * time.Second)
				fmt.Println("initialized 1")
			}
			return nil
		}),
		mainpkg.WithCleanup("testclean", func(context.Context) error {
			if cleanup {
				fmt.Println("cleaning up 1")
				time.Sleep(5 * time.Second)
				fmt.Println("cleaned up 1")
			}
			return nil
		}),
		// For testing serviceutil leader election client.
		// client.WithElector(),
		// mainpkg.WithLeaderFunc(time.Second, func(ctx context.Context) error { return nil }),
		// mainpkg.WithLeaderFunc(0, func(ctx context.Context) error {
		// 	// if !leader {
		// 	// 	return nil
		// 	// }
		// 	t := time.Tick(time.Second)
		// 	for {
		// 		select {
		// 		case <-ctx.Done():
		// 			logging.FromContext(ctx).Trace("leader exit")
		// 			return ctx.Err()
		// 		case <-t:
		// 			logging.FromContext(ctx).Trace("leading")
		// 		}
		// 	}
		// }),
	)
	if err != nil {
		log.Fatal(err)
	}

	svr, err := mainpkg.Setup(config, log,
		mainpkg.WithDualService(mainpkg.External, svc.external...),
		mainpkg.WithServices(mainpkg.External, mainpkg.SvcRegisterFunc(func(s *grpc.Server) error {
			healthv1.RegisterHealthServer(s, health.NewServer()) // health check service
			return nil
		})),
		mainpkg.WithReflection(),
		mainpkg.WithServerOpts(grpc.StatsHandler(&ocgrpc.ServerHandler{IsPublicEndpoint: true})),
		mainpkg.WithLogOpts(grpc_logrus.WithDecider(func(nm string, err error) bool {
			return nm != "/grpc.health.v1.Health/Check" || err != nil
		})),
		mainpkg.WithVersion(svcName, version),
		mainpkg.WithInit("testinit", func(context.Context) error {
			if init {
				fmt.Println("initializing 2")
				time.Sleep(3 * time.Second)
				fmt.Println("initialized 2")
			}
			return nil
		}),
		mainpkg.WithCleanup("testclean", func(context.Context) error {
			if cleanup {
				fmt.Println("cleaning up 2")
				time.Sleep(5 * time.Second)
				fmt.Println("cleaned up 2")
			}
			return nil
		}),
		mainpkg.AdditionalServers(intl),
	)
	if err != nil {
		log.Fatal(err)
	}
	svr.Run()
}

type Svcs struct {
	external []mainpkg.GWGRPC
	internal []mainpkg.GWGRPC
}

func setupServices(ctx context.Context, log *logrus.Entry, config *viper.Viper, conn *grpc.ClientConn) (*Svcs, error) {
	store, err := newDBFromConfig(config)
	if err != nil {
		return nil, err
	}
	if err := store.RunMigration(config.GetString("database.migrationDir")); err != nil {
		return nil, err
	}

	usr := user.New(store, ipupb.NewUserServiceClient(conn))
	sess := session.New(nil)

	return &Svcs{
		external: []mainpkg.GWGRPC{},
		internal: []mainpkg.GWGRPC{usr, sess},
	}, nil
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

func getGRPCOpts(cnf *viper.Viper) []grpc.DialOption {
	if cnf.GetString("runtime.environment") == "development" {
		return []grpc.DialOption{
			grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(10 * time.Second),
			grpc.WithUnaryInterceptor(injectMDInterceptor()),
		}
	}
	return []grpc.DialOption{
		grpc.WithBlock(), grpc.WithTimeout(10 * time.Second),
		grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "")),
	}
}

func injectMDInterceptor() grpc.UnaryClientInterceptor {
	return func(c context.Context, m string, q, r interface{}, n *grpc.ClientConn,
		i grpc.UnaryInvoker, o ...grpc.CallOption,
	) error {
		c = metautils.ExtractIncoming(c).Clone().ToOutgoing(c)
		return i(c, m, q, r, n, o...)
	}
}
