package handler

import (
	"net/http"

	spb "brank.as/petnet/gunk/v1/session"
	"brank.as/petnet/serviceutil/logging"
	"github.com/gorilla/sessions"
)

const (
	sessionCookieExpiry = "expiry"
	sessionCookieToken  = "token"
	sessionCookieUserID = "user-id"
)

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	url := s.sess.AuthCodeURL(authCodeURL)
	promptURL := url + "&prompt=login"
	var sess *sessions.Session
	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		if sess.Options == nil {
			sess.Options = &sessions.Options{}
		}
		sess.Options.MaxAge = -1
		sess.Values[sessionCookieToken] = ""
		if err := s.sess.Save(r, w, sess); err != nil {
			logging.WithError(err, log).Error("saving session")
		}
		logging.WithError(err, log).Error("fetching session")
		http.Redirect(w, r, url, http.StatusTemporaryRedirect)
		return
	}
	uid, ok := sess.Values[sessionCookieUserID].(string)
	if !ok {
		log.Error("missing session user id")
		http.Redirect(w, r, promptURL, http.StatusTemporaryRedirect)
		return
	}
	if res, err := s.pf.GetSession(ctx, &spb.GetSessionRequest{
		IDType: spb.IDType_USERID,
		ID:     uid,
	}); err != nil || res.Expired {
		logging.WithError(err, log).Debug("session expired")
		http.Redirect(w, r, promptURL, http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		logging.WithError(err, log).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusTemporaryRedirect)
		return
	}

	sess.Values[sessionCookieToken] = ""
	if err := s.sess.Save(r, w, sess); err != nil {
		logging.WithError(err, log).Error("deleting session")
		http.Redirect(w, r, errorPath, http.StatusTemporaryRedirect)
		return
	}

	uid, ok := sess.Values[sessionCookieUserID].(string)
	if !ok {
		log.Error("missing session user id")
	}

	if _, err = s.pf.ExpireSession(ctx, &spb.ExpireSessionRequest{
		IDType: spb.IDType_USERID,
		ID:     uid,
	}); err != nil {
		logging.WithError(err, log).Error("expiring session")
	}

	url := s.urls.SSO + "/oauth2/sessions/logout"
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func (s *Server) handleCallback(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	oauthErr := r.FormValue("error")
	if oauthErr != "" {
		log.WithField("oauth_error", oauthErr).Error(r.FormValue("error_description"))
		http.Error(w, oauthErr, http.StatusBadRequest)
		return
	}

	if r.URL.Query().Get(sessionCookieState) != authCodeURL {
		log.Error("session cookie state")
		http.Error(w, "not same sessionCookieState", http.StatusInternalServerError)
		return
	}

	if err := s.sess.Exchange(w, r); err != nil {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	s.GetUserInfoFromCookie(w, r, true)
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
