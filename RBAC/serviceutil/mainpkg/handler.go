package mainpkg

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/pprof"
	rprof "runtime/pprof"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/health"
	hpb "google.golang.org/grpc/health/grpc_health_v1"
)

var (
	mu sync.Mutex
	pp bool = true
)

func (s *server) debug(tggr <-chan struct{}) error {
	mux := http.NewServeMux()
	mux.Handle("/debug/", http.StripPrefix("/debug/", http.HandlerFunc(s.pprof)))
	mux.HandleFunc("/log/trigger", trigger(s.logr))
	mux.Handle("/healthz/", http.StripPrefix("/healthz/", healthz(s.logr, s.health)))
	port := strconv.Itoa(s.debugPort)
	switch s.debugPort {
	default:
		if s.debugPort > 0 {
			break
		}
		fallthrough
	case 80, 443, 0:
		port = "12000"
	}
	srv := http.Server{
		Addr:        net.JoinHostPort("", port),
		Handler:     mux,
		IdleTimeout: 10 * time.Second,
		ErrorLog:    log.New(log.Writer(), "PPROF", log.Default().Flags()),
	}
	go func() {
		<-tggr
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
	}()
	s.logr.WithField("server", "debug").Infof("listening on %q", srv.Addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		s.logr.WithError(err).Error("HTTP")
		return err
	}
	return nil
}

func (s *server) pprof(w http.ResponseWriter, r *http.Request) {
	switch strings.TrimPrefix(r.URL.Path, "/") {
	case "", "pprof":
		pprof.Index(w, r)
	case "trace":
		pprof.Trace(w, r)
	case "profile":
		pprof.Profile(w, r)
	case "cmdline":
		pprof.Cmdline(w, r)
	case "symbol":
		pprof.Symbol(w, r)
	default:
		name := strings.ToLower(strings.TrimPrefix(r.URL.Path, "/"))
		if rprof.Lookup(name) == nil {
			pprof.Index(w, r)
			return
		}
		fmt.Println("name:", name)
		pprof.Handler(name).ServeHTTP(w, r)
	}
}

func healthz(logr *logrus.Logger, h *health.Server) http.HandlerFunc {
	log := logr.WithField("service", "healthz")
	return func(w http.ResponseWriter, r *http.Request) {
		log := log.WithField("request", r.URL.Path)
		if strings.Count(r.URL.Path, "/") > 0 {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte("invalid service"))
			log.Debug("invalid service")
			return
		}

		st, err := h.Check(r.Context(), &hpb.HealthCheckRequest{Service: r.URL.Path})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(hpb.HealthCheckResponse_SERVICE_UNKNOWN.String()))
			log.WithError(err).Error("health service")
			return
		}
		log.WithField("status", st.Status.String()).Trace("handled")
		switch st.Status {
		case hpb.HealthCheckResponse_NOT_SERVING:
			w.WriteHeader(http.StatusServiceUnavailable)
		case hpb.HealthCheckResponse_SERVING:
			w.WriteHeader(http.StatusOK)
		case hpb.HealthCheckResponse_SERVICE_UNKNOWN:
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(st.Status.String()))
	}
}

func trigger(logr *logrus.Logger) http.HandlerFunc {
	logr.WithField("loglevel", logr.Level.String()).Info("logger init")
	log := logr.WithField("service", "log debug trigger")
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			log.WithField("method", http.MethodGet).Info("handled")
		case http.MethodPost, http.MethodPut:
			switch r.URL.Query().Get("level") {
			case "trace":
				logr.SetLevel(logrus.TraceLevel)
			case "debug":
				logr.SetLevel(logrus.DebugLevel)
			case "info":
				logr.SetLevel(logrus.InfoLevel)
			case "warn":
				logr.SetLevel(logrus.WarnLevel)
			case "error":
				logr.SetLevel(logrus.ErrorLevel)
			case "fatal":
				logr.SetLevel(logrus.FatalLevel)
			default:
				switch logr.Level {
				case logrus.DebugLevel:
					logr.SetLevel(logrus.TraceLevel)
				case logrus.TraceLevel:
					logr.SetLevel(logrus.InfoLevel)
				case logrus.InfoLevel:
					logr.SetLevel(logrus.DebugLevel)
				default:
					logr.SetLevel(logrus.InfoLevel)
				}
			}
			log.WithField("log level", logr.Level.String()).Info("set")
		}
		w.Write([]byte(logr.Level.String() + "\n"))
	}
}
