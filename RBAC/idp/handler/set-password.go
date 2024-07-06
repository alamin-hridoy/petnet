package handler

import (
	"html/template"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"

	upb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/idp/auth"
	"brank.as/rbac/serviceutil/logging"
)

type setPasswordFormParams struct {
	CSRFField  template.HTML
	FormErrors map[string]string
	Code       string
	Password   string
	Password2  string
	urls
}

func (s *server) getSetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logrus.WithContext(ctx).WithField("method", "getSetPasswordHandler")

	var fm setPasswordFormParams
	code := r.URL.Query().Get("reset_code")
	if code == "" {
		errMsg := "missing code"
		log.Error(errMsg)
		http.Error(w, errMsg, http.StatusUnprocessableEntity)
		return
	}
	fm.Code = code
	fm.CSRFField = csrf.TemplateField(r)
	fm.urls = s.u
	if err := s.setPasswordTpl.Execute(w, fm); err != nil {
		errMsg := "failed to render set password form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}

func (s *server) postSetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logrus.WithContext(ctx).WithField("method", "postSetPasswordHandler")
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form setPasswordFormParams
	if err := s.dcd.Decode(&form, r.PostForm); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fErr := map[string]string{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Password, validation.Required, validation.Length(8, 64)),
		validation.Field(&form.Password2, validation.Required, validation.Length(8, 64)),
		validation.Field(&form.Code, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Password"] != nil {
				fErr["Password"] = "Password should be 8 to 64 characters"
			}

			if form.Password != form.Password2 {
				fErr["Password2"] = "Password & Confirm password didn't match."
			}
		}
		logging.WithError(err, log).Error("invalid reset password request")
	}
	if len(fErr) > 0 {
		form.FormErrors = fErr
		form.CSRFField = csrf.TemplateField(r)
		form.urls = s.u
		if err := s.setPasswordTpl.Execute(w, form); err != nil {
			errMsg := "failed to render set password form"
			log.WithError(err).Error(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
			return
		}
		return
	}

	_, err := s.suCl.ResetPassword(ctx, &upb.ResetPasswordRequest{
		Code:     form.Code,
		Password: form.Password,
	})
	if err != nil {
		e := auth.FromStatus(err)
		e.MergeDetails(fErr)
		form.FormErrors = e.Errors
		form.CSRFField = csrf.TemplateField(r)
		form.urls = s.u
		if err := s.setPasswordTpl.Execute(w, form); err != nil {
			errMsg := "failed to render set password form"
			log.WithError(err).Error(errMsg)
			http.Error(w, errMsg, http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/user/set-password-success", http.StatusSeeOther)
}

func (s *server) getSetPasswordSuccessHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logrus.WithContext(ctx).WithField("method", "getSetPasswordSuccessHandler")
	if err := s.setPasswordSuccessTpl.Execute(w, setPasswordFormParams{
		urls: s.u,
	}); err != nil {
		errMsg := "failed to render set password success form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}
