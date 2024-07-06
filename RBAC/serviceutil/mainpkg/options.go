package mainpkg

import (
	"context"
	"crypto/tls"
	"strings"
	"time"

	glogrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"github.com/spf13/viper"
)

type Config struct {
	// parsed from env
	corsOrigins          []string
	corsMethods          []string
	corsHeaders          []string
	corsAllowCredentials bool
	corsDebug            bool
	externalService      bool
	env                  string
	socketPath           string
	port                 int
	host                 string
	debugPort            int
	electorSock          string
	mockElector          string
	slackPanicHook       string
	tls                  *tls.Config
	writeTimeout         time.Duration

	// initiated with Option
	middlewares             []Middleware
	registerGateway         []GatewayRegisterer
	grpcServices            []ServiceRegisterer
	unaryInt                []grpc.UnaryServerInterceptor
	streamInt               []grpc.StreamServerInterceptor
	logOpts                 []glogrus.Option
	grpcServer              *grpc.Server
	reflection              bool
	grpcOpts                []grpc.ServerOption
	gwMuxOptions            []runtime.ServeMuxOption
	webHandler              WebHandler
	initFuncs               []opFunc
	initTimeout             time.Duration
	LeadElector             func(string) (Leader, error)
	leader                  *leader
	leadWorkers             map[time.Duration][]LeaderFunc
	leadCron                []cronFunc
	cronFuncs               []cronFunc
	cleanupFuncs            []opFunc
	cleanupTimeout          time.Duration
	gracefulShutdownTimeout time.Duration

	waitReady bool
	cancel    context.CancelFunc

	subs []*server
}

type cronFunc struct {
	name string
	cron Schedule
	next time.Time
	f    LeaderFunc
}

type opFunc struct {
	svcName string
	f       func(context.Context) error
}

type Option func(conf *Config)

// WithTimeout option sets the write timeout for the http server
func WithTimeout(timeout time.Duration) Option {
	return func(conf *Config) { conf.writeTimeout = timeout }
}

// WithPort option override the port in config.
// Key: "server.port"
func WithPort(port int) Option { return func(conf *Config) { conf.port = port } }

// WithPath option override the path in config.
func WithPath(path string) Option { return func(conf *Config) { conf.socketPath = path } }

// WithDebugPort option sets the debug server port (defalt 12000).
//
// Sidecar containers: use a default to avoid port collision.
//
// Local development: set using local config or environment variable,
// Avoid deploying non-standard debug ports, as this hinders debug/alert response.
//
// Key: "debug.port"
func WithDebugPort(port int) Option { return func(conf *Config) { conf.debugPort = port } }

func AddMiddleware(m ...Middleware) Option {
	return func(conf *Config) { conf.middlewares = append(conf.middlewares, m...) }
}

// AdditionalServers registers child servers that will be started and stopped with the parent.
//
// When running multiple grpc/http servers on different ports, set up each port service separately.
// Pass all configured servers to consolidate to a single call to `Run()` in the main function.
func AdditionalServers(srv ...*server) Option {
	return func(conf *Config) { conf.subs = append(conf.subs, srv...) }
}

// OptionList allows configuring a set of options for a given environment,
// avoiding the need to append all standard service options
// to an environment-specific set of options.
func OptionList(o []Option) Option {
	return func(c *Config) {
		for _, opt := range o {
			opt(c)
		}
	}
}

func AddRegisterGatewayFunc(rs ...RegisterGatewayFunc) Option {
	return func(conf *Config) {
		for _, g := range rs {
			conf.registerGateway = append(conf.registerGateway, g)
		}
	}
}

func WithGateway(gw ...GatewayRegisterer) Option {
	return func(conf *Config) { conf.registerGateway = append(conf.registerGateway, gw...) }
}

func AwaitReady() Option { return func(conf *Config) { conf.waitReady = true } }

// Services to be dual-registered as grpc and grpc-gateway.
func WithDualService(external exposedService, gwg ...GWGRPC) Option {
	return func(conf *Config) {
		conf.externalService = bool(external)
		for i := range gwg {
			conf.registerGateway = append(conf.registerGateway, gwg[i])
			conf.grpcServices = append(conf.grpcServices, gwg[i])
		}
	}
}

type exposedService bool

const (
	External exposedService = true
	Internal exposedService = false
)

func WithServices(external exposedService, svc ...ServiceRegisterer) Option {
	return func(conf *Config) {
		conf.externalService = bool(external)
		conf.grpcServices = append(conf.grpcServices, svc...)
	}
}

func WithUnaryInterceptors(in ...grpc.UnaryServerInterceptor) Option {
	return func(conf *Config) { conf.unaryInt = append(conf.unaryInt, in...) }
}

func WithStreamInterceptors(in ...grpc.StreamServerInterceptor) Option {
	return func(conf *Config) { conf.streamInt = append(conf.streamInt, in...) }
}

func WithSlackHook(url string) Option { return func(conf *Config) { conf.slackPanicHook = url } }
func WithLogOpts(opt ...glogrus.Option) Option {
	return func(conf *Config) { conf.logOpts = append(conf.logOpts, opt...) }
}

func WithServerOpts(opt ...grpc.ServerOption) Option {
	return func(conf *Config) { conf.grpcOpts = append(conf.grpcOpts, opt...) }
}

func WithKeepAlive(pingTime, timeout time.Duration) Option {
	return func(c *Config) {
		c.grpcOpts = append(c.grpcOpts,
			grpc.KeepaliveParams(keepalive.ServerParameters{
				Time:    pingTime,
				Timeout: timeout,
			}),
			grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
				// MinTime is the minimum amount of time a client should wait before sending
				// a keepalive ping.
				MinTime:             pingTime,
				PermitWithoutStream: true,
			}),
		)
	}
}

func WithReflection() Option                 { return func(c *Config) { c.reflection = true } }
func WithGRPCServer(srv *grpc.Server) Option { return func(c *Config) { c.grpcServer = srv } }

func WithGatewayProtoError(gwProtoError runtime.ErrorHandlerFunc) Option {
	return AddGatewayMuxOption(runtime.WithErrorHandler(gwProtoError))
}

func AddGatewayMuxOption(option ...runtime.ServeMuxOption) Option {
	return func(conf *Config) { conf.gwMuxOptions = append(conf.gwMuxOptions, option...) }
}

func WithWebHandler(webHandler WebHandler) Option {
	return func(conf *Config) { conf.webHandler = webHandler }
}

// WithInit sets an initializer function for the given service.
// This will handle setting readiness during server initialization.
func WithInit(serviceName string, f func(context.Context) error) Option {
	return func(conf *Config) {
		conf.initFuncs = append(conf.initFuncs, opFunc{svcName: serviceName, f: f})
	}
}

//
func WithInitTimeout(dur time.Duration) Option {
	return func(conf *Config) { conf.initTimeout = dur }
}

// WithCleanup sets an cleanup function for the given service.
// This will handle setting readiness duiring server cleanup.
func WithCleanup(serviceName string, f func(context.Context) error) Option {
	return func(conf *Config) {
		conf.cleanupFuncs = append(conf.cleanupFuncs, opFunc{svcName: serviceName, f: f})
	}
}

func WithCleanupTimeout(dur time.Duration) Option {
	return func(conf *Config) { conf.cleanupTimeout = dur }
}

func WithGracefulShutdownTimeout(dur time.Duration) Option {
	return func(conf *Config) { conf.gracefulShutdownTimeout = dur }
}

// WithCron schedules a background task.
// This will run on every replica, if only one pod needs to run, use WithLeaderCron.
// Especially useful for cache updates or other in-memory tasks.
func WithCron(name string, cron Schedule, f func(context.Context) error) Option {
	return func(conf *Config) {
		conf.cronFuncs = append(conf.cronFuncs, cronFunc{name: name, cron: cron, f: f})
	}
}

// WithLeaderCron schedules a background task on the leader pod only.
// This will only run on one replica, if each pod needs the task, use WithCron.
// Especially useful for database jobs, scheduled reports, or other persistent state tasks.
//
// Leader functions must respect context cancellation, to avoid running after leader lease is lost.
func WithLeaderCron(name string, cron Schedule, f LeaderFunc) Option {
	return func(conf *Config) {
		conf.leadCron = append(conf.leadCron, cronFunc{name: name, cron: cron, f: f})
	}
}

// WithLeaderFunc will execute the function when the pod becomes the leader
// and at every interval while the pod is still leader.
//
// Interval of '0' indicates a long-running task.  It will be started when the pod becomes leader
// and context will only be cancelled when leader lease is lost.
//
// Leader functions must respect context cancellation, to avoid running after leader lease is lost.
func WithLeaderFunc(ival time.Duration, f LeaderFunc) Option {
	return func(conf *Config) { conf.leadWorkers[ival] = append(conf.leadWorkers[ival], f) }
}

// GWGRPC registers both gRPC and gRPC-gateway for a gRPC service.
type GWGRPC interface {
	ServiceRegisterer
	GatewayRegisterer
}

type RegisterGatewayFunc func(ctx context.Context, mux *runtime.ServeMux, gwAddr string, opts []grpc.DialOption) error

func (f RegisterGatewayFunc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, gwAddr string, opts []grpc.DialOption) error {
	return f(ctx, mux, gwAddr, opts)
}

// GatewayRegisterer sets up the grpc-gateway connection on a gRPC service.
type GatewayRegisterer interface {
	RegisterGateway(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error
}
type GWRegisterFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

func (sf GWRegisterFunc) RegisterGateway(c context.Context, m *runtime.ServeMux, a string, o []grpc.DialOption) error {
	return sf(c, m, a, o)
}

// ServiceRegisterer adds the service to the given gRPC server.
type ServiceRegisterer interface {
	RegisterSvc(*grpc.Server) error
}
type SvcRegisterFunc func(*grpc.Server) error

func (sf SvcRegisterFunc) RegisterSvc(s *grpc.Server) error { return sf(s) }

func parseConfig(c *viper.Viper) *Config {
	return &Config{
		corsOrigins:          strings.Split(c.GetString("cors.origins"), ","),
		corsMethods:          strings.Split(c.GetString("cors.methods"), ","),
		corsHeaders:          strings.Split(c.GetString("cors.headers"), ","),
		corsAllowCredentials: c.GetBool("cors.allowCredentials"),
		corsDebug:            c.GetBool("cors.debug"),
		externalService:      true,
		env:                  c.GetString("runtime.environment"),
		port:                 c.GetInt("server.port"),
		debugPort: func() int {
			if p := c.GetInt("debug.port"); p != 0 {
				return p
			}
			return 12000
		}(),
		host:                    c.GetString("server.host"),
		electorSock:             c.GetString("elector.sock"),
		mockElector:             c.GetString("elector.mock_response"),
		writeTimeout:            defaultWriteTO,
		leadWorkers:             map[time.Duration][]LeaderFunc{},
		subs:                    []*server{},
		initTimeout:             defaultInitTO,
		cleanupTimeout:          defaultCleanupTO,
		gracefulShutdownTimeout: defaultGracefulShutdownTO,
	}
}

func WithVersion(name, version string) Option {
	return AddMiddleware(&versionMiddleware{name: name, version: version})
}
