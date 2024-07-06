package main

import (
	"context"
	"embed"
	"encoding/gob"
	"errors"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gorilla/schema"
	"github.com/kenshaw/goji"
	"github.com/kenshaw/sentinel"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/afero"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
	"github.com/yookoala/realpath"

	"brank.as/petnet/cms/handler"
	"brank.as/petnet/cms/internal/core"
	"brank.as/petnet/cms/storage"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/metrics"
	"brank.as/petnet/svcutil/mw"
	"brank.as/rbac/svcutil/otelb"
)

const (
	svcName    = "cms"
	version    = "development"
	timeLayout = "Jan. 2, 2006 15:04:05 MST"
)

//go:embed assets
var assets embed.FS

type SessionUser struct {
	UserID    string `json:"sub"`
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func main() {
	// For using the struct as value in gorilla/sessions
	gob.Register(SessionUser{})

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

	cs := handler.NewConns(log, c)
	cl := handler.NewSvcClients(cs)

	s, err := newServer(log, c, cl, cs)
	if err != nil {
		log.Fatal(err)
	}

	met, err := metrics.NewInfluxDBClient(c)
	if err != nil {
		log.Fatalf("metrics init failed %v", err)
	}
	defer met.Close()
	go met.ErrorsFunc(func(e error) { log.WithError(e).Error("influxdb") })
	s.Use(met.HTTPMiddleware(""))
	s.Use(func(h http.Handler) http.Handler {
		recov := negroni.NewRecovery()
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			recov.ServeHTTP(w, r, h.ServeHTTP)
		})
	})

	l, err := net.Listen("tcp", ":"+c.GetString("server.port"))
	if err != nil {
		log.Fatal(err)
	}
	ss, _ := sentinel.WithContext(context.Background(), os.Interrupt)
	if err := ss.ManageHTTP(l, s); err != nil {
		log.Fatal(err)
	}
	log.Infof("starting server on port :%s", c.GetString("server.port"))
	if err := ss.Run(log, 10*time.Second); err != nil {
		log.Fatal(err)
	}
}

func newServer(logger *logrus.Entry, c *viper.Viper, cl handler.Cl, cs *handler.Conns) (*handler.Server, error) {
	env := c.GetString("runtime.environment")
	logger.WithField("environment", env).Info("configuring service")

	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)

	urls, err := urlsFromConfig(c)
	if err != nil {
		return nil, err
	}

	hmw, err := mw.NewHydra(c, mw.Config{
		IgnorePaths:  []string{"/login", "/oauth2/callback", "/root", "/favicon.ico", "/error"},
		IgnorePrefix: []string{"/fonts", "/images", "/css", "/js", "/templates", "/register", "/registration", "/u/files"},
	}, cs.GetPfInt())
	if err != nil {
		return nil, err
	}
	asst, err := fs.Sub(assets, "assets")
	if err != nil {
		return nil, err
	}
	if env == "localdev" {
		wd, err := os.Getwd()
		if err != nil {
			return nil, err
		}
		assetPath, err := realpath.Realpath(filepath.Join(wd, "assets"))
		if err != nil {
			return nil, err
		}
		asst = afero.NewIOFS(afero.NewBasePathFs(afero.NewOsFs(), assetPath))
	}
	gcs, err := storage.NewOnboardingGCSStorage(logger, c)
	if err != nil {
		logging.WithError(err, logger).Fatal("unable to initialize google cloud storage")
	}

	remcoCommSvc := core.NewRemcoCommissionSvc(cl.GetProfileCL(), cl.GetDRPSandboxCL(), cl.GetDRPLiveCL(), logger)
	srv, err := handler.NewServer(goji.New(), env, logger,
		asst, decoder, urls, hmw, c, cl, cs, gcs, svcName,
		handler.WithRemcoCommissionSvc(remcoCommSvc),
	)
	return srv, err
}

func urlsFromConfig(config *viper.Viper) (handler.URLMap, error) {
	base := config.GetString("server.baseURL")
	if base == "" {
		return handler.URLMap{}, errors.New("missing base URL")
	}
	sso := config.GetString("templates.loginURL")
	if sso == "" {
		return handler.URLMap{}, errors.New("missing login URL")
	}
	return handler.URLMap{
		SSO:  sso,
		Base: base,
	}, nil
}
