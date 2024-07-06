package main

import (
	"net/http"
	"net/url"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/svcutil/random"
)

const pkceSess = "auth"

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context()).WithField("method", "login.handlelogin")
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

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	log := s.log.WithField("method", "handleLogout")

	sess, err := s.sess.Get(r, s.authSession)
	if err != nil {
		errMsg := "fetching session"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
	if sess.Options == nil {
		sess.Options = &sessions.Options{}
	}
	sess.Options.MaxAge = -1

	sessionCookieToken := "token"
	sess.Values[sessionCookieToken] = ""
	if err := s.sess.Save(r, w, sess); err != nil {
		errMsg := "deleting session"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
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
	log := logging.FromContext(r.Context()).WithField("method", "handleCallback")

	oauthErr := r.FormValue("error")
	if oauthErr != "" {
		log.WithField("oauth_error", oauthErr).Error(r.FormValue("error_description"))
		http.Error(w, oauthErr, http.StatusBadRequest)
		return
	}

	ss, err := s.sess.Get(r, pkceSess)
	if err != nil {
		logging.WithError(err, log).Error("getting session")
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	str := ss.Values["pkce"]
	if ss.Options == nil {
		ss.Options = &sessions.Options{}
	}
	ss.Options.MaxAge = -1
	ss.Save(r, w)

	if r.URL.Query().Get(sessionCookieState) != str {
		log.Error("session cookie state")
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}

	if err := s.sess.Exchange(w, r); err != nil {
		log.WithError(err).Error("exchange")
		http.Redirect(w, r, "/login", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
