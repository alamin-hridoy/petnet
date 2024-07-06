package handler

import (
	"net/http"

	"brank.as/petnet/serviceutil/logging"

	mpb "brank.as/petnet/gunk/v1/mfa"
	rbupb "brank.as/rbac/gunk/v1/user"
)

func (s *Server) getPreliminaryScreen(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	code := r.URL.Query().Get("confirm_code")
	if code == "" {
		sess, err := s.sess.Get(r, sessionCookieName)
		if err != nil {
			log.WithError(err).Error("fetching session")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		_, ok := sess.Values[sessionOrgID].(string)
		if !ok {
			log.Error("missing org id")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		_, ok = sess.Values[sessionUserID].(string)
		if !ok {
			log.Error("missing user id")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		template := s.templates.Lookup("preliminary-screen.html")
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
		return
	}

	res, err := s.rbac.EmailConfirmation(ctx, &rbupb.EmailConfirmationRequest{
		Code: code,
	})
	if err != nil {
		logging.WithError(err, log).Error("confirmation code doesn't exist")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	sess.Values[sessionUserID] = res.GetUserID()
	sess.Values[sessionOrgID] = res.GetOrgID()
	if err := s.sess.Save(r, w, sess); err != nil {
		log.WithError(err).Error("saving session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if !s.disableActionMFA {
		if _, err := s.pf.EnableMFA(ctx, &mpb.EnableMFARequest{
			UserID: res.UserID,
			Type:   mpb.MFA_EMAIL,
			Source: res.Email,
		}); err != nil {
			logging.WithError(err, log).Error("enable mfa")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	template := s.templates.Lookup("preliminary-screen.html")
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
