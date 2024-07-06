package mainpkg

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/reflection"

	cors_checker "brank.as/rbac/serviceutil/cors"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/serviceutil/middleware"

	hpb "google.golang.org/grpc/health/grpc_health_v1"
)

const (
	handlerIdleTO             = 120 * time.Second
	serverIdleTO              = 120 * time.Second
	readTO                    = 5 * time.Second
	defaultWriteTO            = 10 * time.Second
	defaultInitTO             = 10 * time.Second
	defaultCleanupTO          = 20 * time.Second
	defaultGracefulShutdownTO = 10 * time.Second
)

type versionMiddleware struct {
	name    string
	version string
}

func (e *versionMiddleware) WrapHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Ver", e.version)
		if r.Method == http.MethodGet {
			switch r.URL.Path {
			case "/":
				fmt.Fprintf(w, "%s server version %s", e.name, e.version)
				return
			case "/version", "/status":
				http.Redirect(w, r, "/", 301)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// mixedHandler returns a new HTTP handler that configured to serve cors handler, gRPC,
// gRPC JSON gateway and grpc-web in a single port.
func (c *Config) mixedHandler(ch *cors.Cors, h http.Handler, gweb *grpcweb.WrappedGrpcServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodOptions:
			ch.HandlerFunc(w, r)
		case c.webHandler != nil && c.webHandler.IsWebRequest(r):
			ch.ServeHTTP(w, r, c.webHandler.ServeHTTP)
		case r.ProtoMajor == 2 && r.Header.Get("Content-Type") == "application/grpc":
			c.grpcServer.ServeHTTP(w, r)
		case gweb.IsGrpcWebRequest(r):
			gweb.ServeHTTP(w, r)
		default:
			ch.ServeHTTP(w, r, h.ServeHTTP)
		}
	})
}

func (c *Config) corsHandler(m *cors_checker.Matcher) *cors.Cors {
	return cors.New(cors.Options{
		AllowOriginFunc:  m.IsAllowedOrigin,
		AllowCredentials: c.corsAllowCredentials,
		AllowedMethods:   c.corsMethods,
		AllowedHeaders:   c.corsHeaders,
		// Adds additional output to debug server side CORS issues
		Debug: c.corsDebug,
	})
}

func (c *Config) grpcwebServer(m *cors_checker.Matcher) *grpcweb.WrappedGrpcServer {
	return grpcweb.WrapServer(c.grpcServer,
		grpcweb.WithOriginFunc(m.IsAllowedOrigin),
		grpcweb.WithAllowedRequestHeaders(c.corsHeaders),
	)
}

func (c *Config) listener(log logrus.FieldLogger) (net.Listener, error) {
	if c.socketPath != "" {
		if !path.IsAbs(c.socketPath) {
			return nil, fmt.Errorf("%q unix socket path must be absolute", c.socketPath)
		}
		if s, err := os.Stat(c.socketPath); err == nil {
			if s.IsDir() {
				return nil, fmt.Errorf("%q unix socket path is a directory", c.socketPath)
			}
			// path exists, remove to avoid collision
			if err := os.Remove(c.socketPath); err != nil {
				return nil, fmt.Errorf("socket cleanup: %w", err)
			}
		}
		// remove socket on shutdown
		c.cleanupFuncs = append(c.cleanupFuncs, opFunc{
			f: func(context.Context) error {
				if err := os.Remove(c.socketPath); err != nil && !errors.Is(err, os.ErrNotExist) {
					return err
				}
				return nil
			},
		})
		// Clear port, use unix socket
		c.port = 0
		l, err := net.Listen("unix", c.socketPath)
		if err != nil {
			return nil, err
		}
		log.WithField("socket", l.Addr()).Info("listening")
		return l, nil
	}

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", c.port))
	if err != nil {
		// do not wrap here, wrap in server()
		return nil, err
	}
	log.WithField("addr", l.Addr()).Info("listening")
	if c.tls != nil {
		tl := tls.NewListener(l, c.tls)
		return tl, nil
	}
	return l, nil
}

func (c *Config) grpcserver(log *logrus.Entry) {
	opt := append(c.grpcOpts,
		grpc.ChainUnaryInterceptor(middleware.New(c.env, log, middleware.Config{
			Internal:          !c.externalService,
			LogOpts:           c.logOpts,
			Interceptors:      c.unaryInt,
			SlackPanicHookURL: c.slackPanicHook,
		})),
		grpc.ChainStreamInterceptor(middleware.NewStream(c.env, log, middleware.StreamConfig{
			Internal:          !c.externalService,
			LogOpts:           c.logOpts,
			Interceptors:      c.streamInt,
			SlackPanicHookURL: c.slackPanicHook,
		})),
	)
	c.grpcServer = grpc.NewServer(opt...)
}

func (c *Config) httpServer(handler http.Handler) *http.Server {
	return &http.Server{
		Handler:      handler,
		ReadTimeout:  readTO,
		WriteTimeout: c.writeTimeout,
		IdleTimeout:  serverIdleTO,
	}
}

func (c *Config) grpcDialOptions() []grpc.DialOption {
	creds := insecure.NewCredentials()
	if c.tls != nil {
		creds = credentials.NewTLS(&tls.Config{ServerName: c.host})
	}
	return []grpc.DialOption{grpc.WithTransportCredentials(creds)}
}

func (c *Config) grpcGateway() *runtime.ServeMux {
	return runtime.NewServeMux(append([]runtime.ServeMuxOption{
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{}),
	}, c.gwMuxOptions...)...)
}

func (c *Config) registerGrpcGateway(gmux *runtime.ServeMux, opts []grpc.DialOption) error {
	grpcAddr := fmt.Sprintf("localhost:%d", c.port)
	ctx, cancel := context.WithCancel(context.Background())
	c.cancel = cancel
	for _, r := range c.registerGateway {
		if err := r.RegisterGateway(ctx, gmux, grpcAddr, opts); err != nil {
			return fmt.Errorf("unable to register gateway: %w", err)
		}
	}
	return nil
}

func (c *Config) handler() (http.Handler, error) {
	corsMatch, err := cors_checker.NewOriginMatcher(c.corsOrigins)
	if err != nil {
		return nil, fmt.Errorf("cannot create cors matcher: %w", err)
	}
	grpcGW := c.grpcGateway()
	if err := c.registerGrpcGateway(grpcGW, c.grpcDialOptions()); err != nil {
		return nil, fmt.Errorf("failed to setup grpc gateway and webserver: %w", err)
	}
	for _, svc := range c.grpcServices {
		if err := svc.RegisterSvc(c.grpcServer); err != nil {
			return nil, fmt.Errorf("failed to register grpc services: %w", err)
		}
	}
	handler := c.mixedHandler(c.corsHandler(corsMatch), grpcGW, c.grpcwebServer(corsMatch))
	for _, middleware := range c.middlewares {
		handler = middleware.WrapHandler(handler)
	}
	return h2c.NewHandler(handler, &http2.Server{IdleTimeout: handlerIdleTO}), nil
}

func (c *Config) server(logr *logrus.Logger) (*server, error) {
	handler, err := c.handler()
	if err != nil {
		return nil, fmt.Errorf("unable to create handler: %w", err)
	}

	listener, err := c.listener(logr)
	if err != nil {
		return nil, fmt.Errorf("unable to listen on port %d: %w", c.port, err)
	}
	health := health.NewServer()
	if c.waitReady {
		health.SetServingStatus("", hpb.HealthCheckResponse_NOT_SERVING)
	}

	ifunc := make(map[string][]func(context.Context) error)
	for _, f := range c.initFuncs {
		ifunc[f.svcName] = append(ifunc[f.svcName], f.f)
	}
	cfunc := make(map[string][]func(context.Context) error)
	for _, f := range c.cleanupFuncs {
		cfunc[f.svcName] = append(cfunc[f.svcName], f.f)
	}

	if (len(c.leadWorkers) != 0 || len(c.leadCron) != 0) && c.leader == nil {
		c.leader, err = newLead(c.electorSock, c.LeadElector, c.mockElector, c.leadWorkers, c.leadCron)
		if err != nil {
			return nil, err
		}
	}

	if c.reflection {
		reflection.Register(c.grpcServer)
	}
	return &server{
		debugPort:       c.debugPort,
		httpServer:      c.httpServer(handler),
		listener:        listener,
		grpcServer:      c.grpcServer,
		logr:            logr,
		cancel:          c.cancel,
		health:          health,
		leader:          c.leader,
		initFunc:        ifunc,
		cleanupFunc:     cfunc,
		cron:            &cronJobs{crontab: c.cronFuncs},
		initTimeout:     c.initTimeout,
		cleanupTimeout:  c.cleanupTimeout,
		shutdownTimeout: c.gracefulShutdownTimeout,
		sub:             c.subs,
	}, nil
}

type server struct {
	lead            bool
	debugPort       int
	httpServer      *http.Server
	listener        net.Listener
	grpcServer      *grpc.Server
	logr            *logrus.Logger
	cancel          context.CancelFunc
	stopCh          chan struct{}
	health          *health.Server
	leader          *leader
	initFunc        map[string][]func(context.Context) error
	cleanupFunc     map[string][]func(context.Context) error
	cron            *cronJobs
	initTimeout     time.Duration
	cleanupTimeout  time.Duration
	shutdownTimeout time.Duration
	sub             []*server
}

type readiness hpb.HealthCheckResponse_ServingStatus

const (
	Ready    readiness = readiness(hpb.HealthCheckResponse_SERVING)
	NotReady readiness = readiness(hpb.HealthCheckResponse_NOT_SERVING)
)

// Ready will set the service name status for use with health and readiness probes.
// Defalt registeres service is the empty string ("").
// All other services are registered when their initial readiness status is set.
func (s *server) Ready(service string, ready readiness) error {
	if s.health == nil {
		return fmt.Errorf("incorrect server, not serving healthz")
	}
	s.health.SetServingStatus(service, hpb.HealthCheckResponse_ServingStatus(ready))
	return nil
}

func (s *server) initFuncs() map[string][]func(context.Context) error {
	fn := s.initFunc
	if fn == nil {
		fn = make(map[string][]func(context.Context) error)
	}
	for _, sb := range s.sub {
		for k, f := range sb.initFuncs() {
			if fn[k] == nil {
				s.Ready(k, NotReady)
			}
			fn[k] = append(fn[k], f...)
		}
	}
	return fn
}

func (s *server) cleanupFuncs() map[string][]func(context.Context) error {
	fn := s.cleanupFunc
	if fn == nil {
		fn = make(map[string][]func(context.Context) error)
	}
	for _, sb := range s.sub {
		for k, f := range sb.cleanupFuncs() {
			if fn[k] == nil {
				s.Ready(k, NotReady)
			}
			fn[k] = append(fn[k], f...)
		}
	}
	return fn
}

func (s *server) start(ctx context.Context) error {
	mu.Lock()
	s.lead = pp
	if pp {
		pp = false
	}
	mu.Unlock()
	if ctx == nil {
		ctx = context.Background()
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ictx := errgroup.WithContext(ctx)
	s.init(ictx, eg)
	// Serve http traffic.
	eg.Go(func() error { return s.httpServer.Serve(s.listener) })
	eg.Go(func() error {
		<-ictx.Done()
		s.stopCh <- struct{}{}
		return nil
	})
	s.logr.WithField("addr", s.listener.Addr().String()).Info("starting")
	for i := range s.sub {
		srv := s.sub[i]
		if s.lead {
			// Ready() can be called on any of the servers to set readiness.
			// Ensure they all point to the leader instance, since that is serving healthz.
			srv.health = s.health
		}
		eg.Go(func() error { return srv.start(ctx) })
	}

	err := eg.Wait()
	if err != http.ErrServerClosed {
		s.logr.WithField("method", "mainpkg.start").Error(err)
	}
	return err
}

func (s *server) init(ctx context.Context, eg *errgroup.Group) {
	if !s.lead {
		// initialize on lead server only
		return
	}
	// Only serve one debug server.
	eg.Go(func() error { return s.debug(ctx.Done()) })
	// Only one crontab.
	eg.Go(func() error {
		return s.cron.startCron(s.withServer(ctx), s.logr.WithField("worker", "crontab"))
	})

	init := s.initFuncs()
	if len(init) == 0 {
		eg.Go(func() error {
			err := s.leader.Start(s.withServer(ctx), s.logr.WithField("worker", "leader"))
			if err != nil {
				s.logr.WithError(err).Error("leader pool failed")
				return err
			}
			return nil
		})
		return
	}
	start := time.Now()
	ictx, cancel := context.WithTimeout(ctx, s.initTimeout)
	ifg, ictx := errgroup.WithContext(ictx)
	for k := range init {
		k, fcs := k, init[k]
		ifg.Go(func() error {
			ig, ctx := errgroup.WithContext(ictx)
			for i := range fcs {
				f := fcs[i]
				ig.Go(func() error { return f(ctx) })
			}
			if err := ig.Wait(); err != nil {
				return err
			}
			s.Ready(k, Ready)
			return nil
		})
	}

	eg.Go(func() error {
		defer cancel()
		if err := ifg.Wait(); err != nil {
			s.logr.WithError(err).Error("init")
			return err
		}
		s.logr.WithField("delay", time.Since(start).String()).Info("initialization complete")
		// One leader pool.
		eg.Go(func() error {
			err := s.leader.Start(s.withServer(ctx), s.logr.WithField("worker", "leader"))
			if err != nil {
				s.logr.WithError(err).Error("leader pool failed")
				return err
			}
			return nil
		})
		return nil
	})
}

func (s *server) stop(ctx context.Context) error {
	s.cancel()
	s.health.Shutdown()
	eg, ctx := errgroup.WithContext(ctx)
	// shutdown context for the http server.
	sctx, cancel := context.WithTimeout(ctx, 15*time.Second)
	eg.Go(func() error {
		// block and shutdown the gRPC connections before http
		s.grpcServer.GracefulStop()
		if err := s.httpServer.Shutdown(sctx); err != nil {
			s.logr.Error(err)
		}
		cancel() // signal that connections are closed, cleanup can begin.
		return nil
	})
	for i := range s.sub { // trigger stop on all servers
		srv := s.sub[i]
		eg.Go(func() error { srv.stop(ctx); return nil })
	}
	<-sctx.Done() // wait for active connections to close
	s.grpcServer.Stop()
	if err := s.httpServer.Close(); err != nil {
		s.logr.Error(err)
	}
	s.cleanup(ctx, eg)
	return eg.Wait()
}

// cleanup server resources.
func (s *server) cleanup(ctx context.Context, eg *errgroup.Group) {
	if !s.lead {
		return
	}
	// Run cleanup
	cleanup := s.cleanupFuncs()
	if len(cleanup) == 0 {
		return
	}
	start := time.Now()
	cctx, cancel := context.WithTimeout(ctx, s.cleanupTimeout)
	cg, cctx := errgroup.WithContext(cctx)
	for k := range cleanup {
		k, fcs := k, cleanup[k]
		cg.Go(func() error {
			ig, ctx := errgroup.WithContext(cctx)
			for i := range fcs {
				f := fcs[i]
				ig.Go(func() error { return f(ctx) })
			}
			if err := ig.Wait(); err != nil {
				return fmt.Errorf("cleanup %q: %w", k, err)
			}
			return nil
		})
	}
	eg.Go(func() error {
		defer cancel()
		if err := cg.Wait(); err != nil {
			return err
		}
		s.logr.WithField("delay", time.Since(start).String()).Info("cleanup complete")
		return nil
	})
}

func Setup(config *viper.Viper, log *logrus.Entry, opts ...Option) (*server, error) {
	c := parseConfig(config)
	for _, opt := range opts {
		opt(c)
	}
	if c.grpcServer == nil {
		c.grpcserver(log)
	}
	for i, cr := range c.cronFuncs {
		if cr.cron == nil {
			return nil, fmt.Errorf("invalid cronjob %q: no Schedule", cr.name)
		}
		if err := c.cronFuncs[i].cron.Parse(); err != nil {
			return nil, fmt.Errorf("invalid cronjob schedule %q: %w", cr.name, err)
		}
	}
	for _, s := range c.subs {
		if s.leader == nil {
			continue
		}
		c.leadCron = append(c.leadCron, s.leader.cron...)
		if c.LeadElector == nil {
			c.leader = s.leader
			continue
		}
		s.leader.le.Close()
		c.cronFuncs = append(c.cronFuncs, s.cron.crontab...)
		s.cron.crontab = nil
	}
	if c.mockElector != "" && c.env != "development" {
		return nil, fmt.Errorf("mock elector allowed in development env only (%q)", c.env)
	}
	return c.server(log.Logger)
}

// rootCtx is used for running the server in tests. Cancellation will stop the server.
var rootCtx = context.Background()

func (s *server) Run() {
	log := s.logr.WithField("method", "mainpkg.Run")
	cch := make(chan struct{})
	s.stopCh = make(chan struct{}, 1)
	defer close(s.stopCh)

	ctx, cancel := context.WithCancel(rootCtx)
	go func() {
		// catch interrupt signals
		ch := signals()

		select {
		case sig := <-ch:
			log.WithField("signal", sig.String()).Info("shutting down")
		case <-s.stopCh:
		}
		// perform graceful shutdown
		ctx, shutdownCancel := context.WithTimeout(ctx, s.shutdownTimeout)
		defer shutdownCancel()
		go func() {
			sig := <-ch
			log.WithField("signal", sig.String()).Info("force shutdown")
			cancel()
		}()
		if err := s.stop(ctx); err != nil {
			logging.WithError(err, log).Error("shutdown")
		}
		close(cch)
	}()
	// serve application
	if err := s.start(ctx); err != nil && err != http.ErrServerClosed {
		log.WithError(err).Error("server exited")
	}
	log.Info("server exiting")
	<-cch
	log.Info("server exited")
	return
}
