package handler

import (
	"net/http"

	"github.com/sirupsen/logrus"

	upb "brank.as/rbac/gunk/v1/user"
)

type confirmTmpl struct {
	Username  string
	FirstName string
	LastName  string
	Email     string
	urls
}

func (s *server) getConfirmEmailHandler(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithContext(r.Context()).WithField("method", "getConfirmEmailHandler")
	ctx := r.Context()

	code := r.URL.Query().Get("confirm_code")
	if code == "" {
		log.Error("confirm code can't be empty")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	u, err := s.suCl.EmailConfirmation(ctx, &upb.EmailConfirmationRequest{Code: code})
	if err != nil {
		log.WithError(err).Error("code invalid")
		http.Redirect(w, r, s.ErrRedirURL, http.StatusTemporaryRedirect)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	if err := s.confirmEmailTpl.Execute(w, confirmTmpl{
		Username:  u.GetUsername(),
		FirstName: u.GetFirstName(),
		LastName:  u.GetLastName(),
		Email:     u.GetEmail(),
		urls:      s.u,
	}); err != nil {
		log.WithError(err).Error("failed to render")
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}
