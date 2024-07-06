package handler

import (
	"html/template"
	"net/http"

	"brank.as/petnet/serviceutil/logging"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
)

type VerifyAccountForm struct {
	CSRFField template.HTML
	UserID    string
	Email     string
}

func (s *Server) getVerifyAccount(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid, ok := sess.Values[sessionUserID].(string)
	if !ok {
		log.Error("missing user id")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	email, ok := sess.Values[sessionEmail].(string)
	if !ok {
		log.Error("missing email")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("verify-account.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := template.Execute(w, VerifyAccountForm{
		CSRFField: csrf.TemplateField(r),
		UserID:    uid,
		Email:     email,
	}); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postVerifyAccount(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var f VerifyAccountForm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&f,
		validation.Field(&f.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if _, err := s.rbac.ResendConfirmEmail(ctx, &rbupb.ResendConfirmEmailRequest{
		UserID: f.UserID,
	}); err != nil {
		logging.WithError(err, log).Error("resending confirmation email")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, verifyAccountPath, http.StatusSeeOther)
}
