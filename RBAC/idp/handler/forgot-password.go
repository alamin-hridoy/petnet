package handler

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"

	upb "brank.as/rbac/gunk/v1/user"
)

var resetTokenExpiry time.Duration = 1 * time.Hour

type forgotPasswordFormParams struct {
	CSRFField           template.HTML
	ForgotPasswordError string
	EmailSent           bool
	urls
}

func (s *server) getForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	p := &forgotPasswordFormParams{urls: s.u}
	s.showForgotPasswordForm(p, w, r)
}

func (s *server) postForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithContext(r.Context()).WithField("method", "postForgotPasswordHandler")
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	usrEmail := r.FormValue("email")
	if usrEmail == "" {
		log.Error(errors.New("no email in form"))
		s.showForgotPasswordForm(&forgotPasswordFormParams{
			ForgotPasswordError: "Please, provide a valid email address",
			urls:                s.u,
		}, w, r)
		return
	}

	if !isEmailValidFormat(usrEmail) {
		log.Error(errors.New("wrong email format"))
		s.showForgotPasswordForm(&forgotPasswordFormParams{
			ForgotPasswordError: "Please, provide a valid email address",
			urls:                s.u,
		}, w, r)
		return
	}

	_, err := s.suCl.ForgotPassword(ctx, &upb.ForgotPasswordRequest{
		Email: usrEmail,
	})
	if err != nil {
		log.WithError(err).Error("resetting password")
		s.showForgotPasswordForm(&forgotPasswordFormParams{
			ForgotPasswordError: "Unknown Email Address",
			urls:                s.u,
		}, w, r)
		return
	}

	s.showForgotPasswordConfirmationForm(&forgotPasswordFormParams{
		EmailSent: true,
		urls:      s.u,
	}, w, r)
}

func (s *server) showForgotPasswordForm(p *forgotPasswordFormParams, w http.ResponseWriter, r *http.Request) {
	log := logrus.WithContext(r.Context()).WithField("method", "showForgotPasswordForm")

	p.CSRFField = csrf.TemplateField(r)
	if err := s.forgotPasswordTpl.Execute(w, p); err != nil {
		log.Printf("failed to render forgot password form: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}

func (s *server) showForgotPasswordConfirmationForm(p *forgotPasswordFormParams, w http.ResponseWriter, r *http.Request) {
	log := logrus.WithContext(r.Context()).WithField("method", "showForgotPasswordConfirmationForm")

	if err := s.forgotPwdConfirmTpl.Execute(w, p); err != nil {
		log.Printf("failed to render forgot password confirmation form: %s", err.Error())
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

// isEmailValidFormat checks if the email provided passes the required structure and length.
func isEmailValidFormat(e string) bool {
	if len(e) < 3 && len(e) > 254 {
		return false
	}
	return emailRegex.MatchString(e)
}
