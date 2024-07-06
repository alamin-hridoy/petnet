package handler

import (
	"errors"
	"html/template"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"

	ipb "brank.as/rbac/gunk/v1/invite"
	upb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"
)

type inviteSetPasswordForm struct {
	CSRFField template.HTML
	Errors    map[string]error
	ResetCode string
	Password  string
	Password2 string
	UserID    string
	urls
}

func (s *server) getInviteSetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logrus.WithContext(ctx).WithField("method", "getInviteSetPassword")

	var f inviteSetPasswordForm
	rc := r.URL.Query().Get("reset_code")
	if rc == "" {
		errMsg := "missing reset code"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return
	}
	ic := r.URL.Query().Get("invite_code")
	if ic == "" {
		errMsg := "missing invite code"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return
	}

	usr, err := s.invCl.RetrieveInvite(ctx, &ipb.RetrieveInviteRequest{
		Code: ic,
	})
	if err != nil {
		logging.WithError(err, log).Error("retrieving invite")
		http.Error(w, "unable to retrieve invite", http.StatusUnauthorized)
		return
	}

	if usr.InviteStatus != storage.InviteSent {
		logging.WithError(err, log).Error("wrong invite status")
		http.Error(w, "invite has expired", http.StatusUnauthorized)
		return
	}

	f.UserID = usr.ID
	f.ResetCode = rc
	f.CSRFField = csrf.TemplateField(r)
	f.urls = s.u
	if err := s.inviteSetPasswordTpl.Execute(w, f); err != nil {
		errMsg := "failed to render invite set password form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (s *server) postInviteSetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "postInviteSetPassword")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	var f inviteSetPasswordForm
	if err := s.dcd.Decode(&f, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	fErr := validation.Errors{}
	if err := validation.ValidateStruct(&f,
		validation.Field(&f.Password, validation.Required, validation.Length(8, 50)),
		validation.Field(&f.Password2, validation.Required, validation.Length(8, 50)),
		validation.Field(&f.ResetCode, validation.Required),
		validation.Field(&f.UserID, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Password"] != nil {
				fErr["Password"] = errors.New("Password should be at least 8 characters")
			}
			if f.Password != f.Password2 {
				fErr["Password"] = errors.New("Passwords should match")
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}
	f.Errors = fErr
	f.CSRFField = csrf.TemplateField(r)
	if len(fErr) > 0 {
		if err := s.inviteSetPasswordTpl.Execute(w, f); err != nil {
			logging.WithError(err, log).Error("template execution")
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
	}

	if _, err := s.suCl.ResetPassword(ctx, &upb.ResetPasswordRequest{
		Code:     f.ResetCode,
		Password: f.Password,
	}); err != nil {
		logging.WithError(err, log).Error("resetting password")
		http.Error(w, "failed resetting password", http.StatusUnprocessableEntity)
		return
	}

	if _, err := s.invCl.Approve(ctx, &ipb.ApproveRequest{
		ID: f.UserID,
	}); err != nil {
		logging.WithError(err, log).Error("approving invite")
		http.Error(w, "unable to retrieve invite", http.StatusUnauthorized)
		return
	}

	http.Redirect(w, r, "/invite/set-password-success", http.StatusSeeOther)
}

func (s *server) getInviteSetPasswordSuccess(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logrus.WithContext(ctx).WithField("method", "getInviteSetPasswordSuccess")

	if err := s.inviteSetPasswordSuccessTpl.Execute(w, setPasswordFormParams{
		urls: s.u,
	}); err != nil {
		errMsg := "failed to render invite set password success form"
		logging.WithError(err, log).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}
