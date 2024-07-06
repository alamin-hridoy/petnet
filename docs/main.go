package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/kenshaw/goji"
	"github.com/kenshaw/sentinel"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/yookoala/realpath"

	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
)

const (
	svcName = "docs"
	version = "development"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
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

	log := logging.NewLogger(config).WithFields(logrus.Fields{
		"service": svcName,
		"version": version,
	})
	log.Info("starting service")
	s, err := newServer(config, log)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", ":"+config.GetString("server.port"))
	if err != nil {
		return err
	}
	ss, _ := sentinel.WithContext(context.Background(), os.Interrupt)
	if err := ss.ManageHTTP(l, s); err != nil {
		return err
	}
	log.Infof("starting server on port :%s", config.GetString("server.port"))
	return ss.Run(log, 10*time.Second)
}

type Server struct {
	*goji.Mux
	logger    logrus.FieldLogger
	assets    string
	templates *template.Template
}

func newServer(config *viper.Viper, log *logrus.Entry) (*Server, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	assets, err := realpath.Realpath(filepath.Join(wd, "assets"))
	if err != nil {
		return nil, err
	}

	templates := template.Must(template.New("").Delims("^^", "%%").
		ParseGlob(filepath.Join(assets, "templates", "*.html")))
	s := &Server{
		Mux:       goji.New(),
		logger:    log,
		assets:    assets,
		templates: templates,
	}

	s.HandleFunc(goji.Get("/docs"), s.handleDocs)
	s.HandleFunc(goji.NewPathSpec("/*"), s.handleIndex)
	s.Use(mw.Logger(log))
	return s, nil
}

func (s *Server) handleDocs(res http.ResponseWriter, req *http.Request) {
	template := s.templates.Lookup("redoc.html")
	if template == nil {
		http.Error(res, "unable to load template", http.StatusInternalServerError)
		return
	}
	err := template.Execute(res, nil)
	if err != nil {
		s.logger.Infof("error with template execution: %+v", err)
	}
}

func (s *Server) handleIndex(res http.ResponseWriter, req *http.Request) {
	base := s.assets
	switch {
	case req.URL.Path == "" || req.URL.Path == "/":
		template := s.templates.Lookup("redoc.html")
		if template == nil {
			http.Error(res, "unable to load template", http.StatusInternalServerError)
			return
		}
		err := template.Execute(res, nil)
		if err != nil {
			s.logger.Infof("error with template execution: %+v", err)
		}
		return
	case filepath.Ext(req.URL.Path) == ".html":
		base = filepath.Join(s.assets, "templates")
	}
	http.FileServer(http.Dir(base)).ServeHTTP(res, req)
}
