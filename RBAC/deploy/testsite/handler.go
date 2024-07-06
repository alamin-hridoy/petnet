package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/fs"
	"net/http"
	"strings"

	"github.com/Masterminds/sprig"
	"github.com/gorilla/csrf"
	"github.com/gorilla/schema"
	"github.com/kenshaw/goji"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/mw"
)

const sessionCookieState = "state"

const (
	loginPath  = "/login"
	logoutPath = "/logout"
)

type Server struct {
	templates *template.Template

	*goji.Mux
	log     *logrus.Entry
	assets  fs.FS
	decoder *schema.Decoder

	sess        *mw.Hydra
	authSession string
}

func NewServer(config *viper.Viper,
	log *logrus.Entry,
	hmw *mw.Hydra,
	assets fs.FS,
) (*Server, error) {
	s := &Server{
		Mux:         goji.New(),
		log:         log,
		assets:      assets,
		sess:        hmw,
		authSession: config.GetString("auth.cookiename"),
	}
	if err := s.parseTemplates(); err != nil {
		return nil, err
	}
	log.Logger.SetFormatter(&logrus.JSONFormatter{})

	s.Mux.Use(logging.LoggerMiddleware(log))
	s.Mux.Use(hmw.Middleware)
	s.Mux.Use(csrf.Protect([]byte(config.GetString("csrf.secret")),
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logging.FromContext(r.Context())
			log.WithFields(logrus.Fields{
				"csrf_error": csrf.FailureReason(r).Error(),
				"token":      csrf.Token(r),
				"template":   csrf.TemplateField(r),
			}).Error("csrf error")
			fmt.Fprintln(w, csrf.FailureReason(r))
		})),
		csrf.Path("/")))

	s.HandleFunc(goji.Get(loginPath), s.handleLogin)
	s.HandleFunc(goji.Get(logoutPath), s.handleLogout)
	s.HandleFunc(goji.NewPathSpec("/oauth2/callback"), s.handleCallback)

	s.HandleFunc(goji.NewPathSpec("/*"), s.handleIndex)

	return s, nil
}

func (s *Server) lookupTemplate(name string) *template.Template {
	return s.templates.Lookup(name)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).
		WithField("method", "handleIndex").WithField("path", r.URL.Path)
	log.Info("received")
	if _, err := fs.Stat(s.assets, r.URL.Path); err == nil {
		http.FileServer(http.FS(s.assets)).ServeHTTP(w, r)
		return
	}
	sess, err := s.sess.Get(r, s.authSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tok := sess.Values["token"]
	if err := s.lookupTemplate("home.html").Execute(w, map[string]interface{}{
		"Token": tok,
	}); err != nil {
		http.Error(w, "template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

func (s *Server) templateData(r *http.Request) TemplateData {
	return TemplateData{
		Env:       "testsite",
		CSRFField: csrf.TemplateField(r),
	}
}

func (s *Server) doTemplate(w http.ResponseWriter, r *http.Request, name string, status int) error {
	template := s.lookupTemplate(name)
	if template == nil || isPartialTemplate(name) {
		template, status = s.templates.Lookup("error.html"), http.StatusNotFound
	}

	w.WriteHeader(status)
	return template.Execute(w, s.templateData(r))
}

type TemplateData struct {
	Env       string
	CSRFField template.HTML
}

func isPartialTemplate(name string) bool {
	return strings.HasSuffix(name, ".part.html")
}

func (s *Server) parseTemplates() error {
	templates := template.New("cms-templates").Funcs(template.FuncMap{
		"jsonStringify": func(data interface{}) string {
			mapB, _ := json.Marshal(data)
			return string(mapB)
		},
	}).Funcs(sprig.FuncMap())

	tmpl, err := templates.ParseFS(s.assets, "assets/templates/*.html")
	if err != nil {
		return err
	}
	s.templates = tmpl
	return nil
}
