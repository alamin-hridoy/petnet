package handler

import (
	"errors"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/csrf"
	"github.com/ory/hydra-client-go/client/admin"
)

func (s *server) showLogoutForm(challenge string, w http.ResponseWriter, r *http.Request) {
	if err := s.logoutTpl.Execute(w, struct {
		Challenge string
		CSRFField template.HTML
	}{
		Challenge: challenge,
		CSRFField: csrf.TemplateField(r),
	}); err != nil {
		log.Printf("failed to render logout form: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}

func (s *server) getLogoutHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	challengeKey := r.URL.Query().Get("logout_challenge")

	lr, err := s.hydra.GetLogoutRequest(admin.NewGetLogoutRequestParamsWithContext(ctx).
		WithLogoutChallenge(challengeKey))
	if err != nil {
		if errors.Is(err, admin.NewGetLogoutRequestGone()) {
			log.Printf("logout challenge is already processed: challenge=%s", challengeKey)
			http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
			return
		}
		log.Printf("failed to get logout request: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}

	// if the logout request is initiated by relaying party, show the confirmation form.
	if s.logoutConfirmation && lr.Payload.RpInitiated {
		s.showLogoutForm(challengeKey, w, r)
		return
	}
	s.acceptLogout(challengeKey, w, r)
}

type logoutForm struct {
	Challenge string
	Accept    bool
}

func (s *server) rejectLogout(challenge string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	_, err := s.hydra.RejectLogoutRequest(admin.NewRejectLogoutRequestParamsWithContext(ctx).
		WithLogoutChallenge(challenge))
	if err != nil {
		log.Printf("failed to reject logout request: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, s.ErrRedirURL, http.StatusFound)
}

func (s *server) acceptLogout(challenge string, w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	acpt, err := s.hydra.AcceptLogoutRequest(admin.NewAcceptLogoutRequestParamsWithContext(ctx).
		WithLogoutChallenge(challenge))
	if err != nil {
		log.Printf("failed to accept logout request: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, *acpt.Payload.RedirectTo, http.StatusFound)
}

func (s *server) postLogoutHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form logoutForm
	if err := s.dcd.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if form.Accept {
		s.acceptLogout(form.Challenge, w, r)
		return
	}
	s.rejectLogout(form.Challenge, w, r)
}
