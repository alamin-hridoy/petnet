package handler

import (
	"net/http"

	"brank.as/petnet/serviceutil/logging"
)

func (s *Server) getRegisterSuccess(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	_, ok := sess.Values[sessionUserID].(string)
	if !ok {
		log.Error("missing user id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	_, ok = sess.Values[sessionOrgID].(string)
	if !ok {
		log.Error("missing org id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("success-page.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := template.Execute(w, nil); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
