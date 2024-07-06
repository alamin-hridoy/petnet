package main

import (
	"html/template"
	"net/http"
	"net/url"

	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"brank.as/petnet/svcutil/random"
	"github.com/gorilla/csrf"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

const pkceSess = "auth"

type loginForm struct {
	CSRFField template.HTML
	Email     string
	Password  string
}

func (s *Server) handleOAuthLogin(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ss, err := s.sess.Get(r, pkceSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
	}
	str, err := random.String(16)
	if err != nil {
		logging.WithError(err, log).Error("random")
		http.Error(w, "random "+err.Error(), http.StatusInternalServerError)
		return
	}
	if ss != nil {
		ss.Values["pkce"] = str
		ss.Save(r, w)
	}

	aopts := []oauth2.AuthCodeOption{}
	for k, v := range r.URL.Query() {
		for _, o := range v {
			aopts = append(aopts, oauth2.SetAuthURLParam(k, o))
		}
	}
	u := s.sess.AuthCodeURL(str, aopts...)
	logging.FromContext(r.Context()).WithField("authredir", u).Info("redirecting")
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (s *Server) getBasicLogin(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	template := s.templates.Lookup("login.html")
	if template == nil {
		log.Error("unable to load template")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	if err := template.Execute(w, loginForm{
		CSRFField: csrf.TemplateField(r),
	}); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
}

func (s *Server) postBasicLogin(w http.ResponseWriter, r *http.Request) {
	var form loginForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if form.Email != "teller@mail.com" || form.Password != "secret" {
		http.Error(w, "permission denied, wrong email or password", http.StatusForbidden)
		return
	}

	sess, err := s.sess.Get(r, s.authSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	sess.Values[mw.SessionAuthType] = mw.BasicAuthType
	if err := s.sess.Save(r, w, sess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	sess, err := s.sess.Get(r, s.authSession)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if authType, ok := sess.Values[mw.SessionAuthType]; ok && authType == mw.BasicAuthType {
		sess.Values[mw.SessionAuthType] = ""
		if err := s.sess.Save(r, w, sess); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	if sess.Options == nil {
		sess.Options = &sessions.Options{}
	}
	sess.Options.MaxAge = -1

	sessionCookieToken := "token"
	sess.Values[sessionCookieToken] = ""
	if err := s.sess.Save(r, w, sess); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	u, err := url.Parse(s.sess.Endpoint.AuthURL)
	if err != nil {
		u = &url.URL{Scheme: "http", Host: "127.0.0.1:3000"}
	}
	u.Path = "/oauth2/sessions/logout"
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	oauthErr := r.FormValue("error")
	if oauthErr != "" {
		log.WithField("oauth_error", oauthErr).Error(r.FormValue("error_description"))
		http.Error(w, oauthErr, http.StatusBadRequest)
		return
	}

	ss, err := s.sess.Get(r, pkceSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Redirect(w, r, "/auth-option", http.StatusInternalServerError)
		return
	}
	str := ss.Values["pkce"]
	if ss.Options == nil {
		ss.Options = &sessions.Options{}
	}
	ss.Options.MaxAge = -1
	ss.Save(r, w)

	sessionCookieState := "state"
	if r.URL.Query().Get(sessionCookieState) != str {
		log.Error("session cookie state")
		http.Redirect(w, r, "/auth-option", http.StatusInternalServerError)
		return
	}

	if err := s.sess.Exchange(w, r); err != nil {
		log.WithError(err).Error("exchange")
		http.Redirect(w, r, "/auth-option", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
