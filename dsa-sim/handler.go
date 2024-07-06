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

	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
)

const (
	loginOAuthPath      = "/login-oauth"
	loginBasicPath      = "/login-basic"
	logoutPath          = "/logout"
	authOptionPath      = "/auth-option"
	choosePartnerPath   = "/choose-partner"
	remitPath           = "/remit"
	listRemitPath       = "/list-remit"
	createRemitPath     = "/create-remit"
	disburseRemitPath   = "/disburse-remit"
	registerUserPath    = "/register-user"
	createProfilePath   = "/create-profile"
	createRecipientPath = "/create-recipient"
	createQuotePath     = "/create-quote"
)

type Server struct {
	templates *template.Template

	*goji.Mux
	log     *logrus.Entry
	assets  fs.FS
	decoder *schema.Decoder

	sess        *mw.Hydra
	authSession string
	cl          cl
}

func NewServer(config *viper.Viper,
	log *logrus.Entry,
	hmw *mw.Hydra,
	assets fs.FS,
	cl cl,
) (*Server, error) {
	s := &Server{
		Mux:         goji.New(),
		log:         log,
		assets:      assets,
		sess:        hmw,
		authSession: config.GetString("auth.cookiename"),
		decoder:     schema.NewDecoder(),
		cl:          cl,
	}
	s.decoder.IgnoreUnknownKeys(true)
	if err := s.parseTemplates(); err != nil {
		return nil, err
	}
	log.Logger.SetFormatter(&logrus.JSONFormatter{})

	s.Mux.Use(logging.LoggerMiddleware(log))
	s.Mux.Use(hmw.DSASimMiddleware)
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

	s.HandleFunc(goji.Get(loginOAuthPath), s.handleOAuthLogin)
	s.HandleFunc(goji.Get(loginBasicPath), s.getBasicLogin)
	s.HandleFunc(goji.Post(loginBasicPath), s.postBasicLogin)
	s.HandleFunc(goji.Get(logoutPath), s.handleLogout)
	s.HandleFunc(goji.NewPathSpec("/oauth2/callback"), s.handleCallback)

	s.HandleFunc(goji.Get(authOptionPath), s.getAuthOption)

	s.HandleFunc(goji.Post(choosePartnerPath), s.postChoosePartner)
	s.HandleFunc(goji.Get(remitPath), s.getRemit)
	s.HandleFunc(goji.Post(createRemitPath), s.postCreateRemit)
	s.HandleFunc(goji.Post(disburseRemitPath), s.postDisburseRemit)
	s.HandleFunc(goji.Get(listRemitPath), s.getListRemit)
	s.HandleFunc(goji.Post(registerUserPath), s.postRegisterUser)
	s.HandleFunc(goji.Post(createProfilePath), s.postCreateProfile)
	s.HandleFunc(goji.Post(createRecipientPath), s.postCreateRecipient)
	s.HandleFunc(goji.Post(createQuotePath), s.postCreateQuote)

	s.HandleFunc(goji.NewPathSpec("/*"), s.handleIndex)

	return s, nil
}

func (s *Server) lookupTemplate(name string) *template.Template {
	return s.templates.Lookup(name)
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	if _, err := fs.Stat(s.assets, r.URL.Path); err == nil {
		http.FileServer(http.FS(s.assets)).ServeHTTP(w, r)
		return
	}
	if r.URL.Path != "" && r.URL.Path != "/" {
		return
	}
	http.Redirect(w, r, remitPath, http.StatusSeeOther)
	return
}

func (s *Server) templateData(r *http.Request) TemplateData {
	return TemplateData{
		Env:       "dsa-sim",
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
