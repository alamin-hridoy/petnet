package handler

import (
	"html/template"
	"net/http"
	"regexp"

	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	epb "brank.as/rbac/gunk/v1/errors"
	ipb "brank.as/rbac/gunk/v1/invite"
	upb "brank.as/rbac/gunk/v1/user"
)

type signupFormParams struct {
	CSRFField  template.HTML
	FormErrors map[string]interface{}
	Username   string
	Email      string
	FirstName  string
	LastName   string
	Password   string
	Password2  string
	InvCode    string
	OrgName    string
	urls
}

func (s *server) getSignupHandler(w http.ResponseWriter, r *http.Request) {
	log := logrus.WithContext(r.Context()).WithField("method", "getSignupHandler")

	var form signupFormParams
	form.FormErrors = make(map[string]interface{})
	form.urls = s.u

	code := r.URL.Query().Get("invite_code")
	form.InvCode = code
	if code != "" {
		inv, err := s.invCl.RetrieveInvite(r.Context(), &ipb.RetrieveInviteRequest{Code: code})
		if err != nil {
			logging.WithError(err, log).Error("retrieve invite")
		}
		if inv.GetActive() { // Prefill invite data
			form.FirstName = inv.GetFirstName()
			form.LastName = inv.GetLastName()
			form.OrgName = inv.GetCompanyName()
		}
	}

	w.Header().Set("Content-Type", "text/html")
	form.CSRFField = csrf.TemplateField(r)
	if err := s.signupTpl.Execute(w, form); err != nil {
		log.WithError(err).Error("failed to render signup form")
		http.Error(w, genericErrMsg, http.StatusInternalServerError)
		return
	}
}

func (s *server) postSignupHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logrus.WithContext(ctx).WithField("method", "postSignupHandler")
	if err := r.ParseForm(); err != nil {
		log.WithError(err).Error("parsing form")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form signupFormParams
	if err := s.dcd.Decode(&form, r.PostForm); err != nil {
		log.WithError(err).Error("decoding form")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	form.FormErrors = map[string]interface{}{}
	pRe := regexp.MustCompile("[0-9]+")
	switch {
	case form.Email == "" && form.InvCode == "":
		form.FormErrors["Email"] = "Email is required and needs to be valid."
	case len(form.Password) < 8 || len(pRe.FindAllString(form.Password, -1)) == 0:
		form.FormErrors["Password"] = "Password should be at least 8 characters long and contain at least 1 number."
	case form.Password2 != form.Password:
		form.FormErrors["Password2"] = "Passwords do not match."
	}

	if len(form.FormErrors) != 0 {
		log.Error("form errors:", form.FormErrors)
		form.CSRFField = csrf.TemplateField(r)
		if err := s.signupTpl.Execute(w, form); err != nil {
			log.WithError(err).Error("failed to render signup form")
			http.Error(w, genericErrMsg, http.StatusInternalServerError)
		}
		return
	}

	_, err := s.suCl.Signup(ctx, &upb.SignupRequest{
		Username:   form.Username,
		FirstName:  form.FirstName,
		LastName:   form.LastName,
		Email:      form.Email,
		Password:   form.Password,
		InviteCode: form.InvCode,
	})
	if err != nil {
		log.WithError(err).Error("signup error")
		switch status.Code(err) {
		case codes.FailedPrecondition:
			form.FormErrors["InvitationCodeUsed"] = true
		case codes.AlreadyExists:
			form.FormErrors["EmailExists"] = true
			form.FormErrors["SignupURL"] = s.u.SignupURL
		case codes.InvalidArgument:
			st := status.Convert(err)
			d := st.Proto().GetDetails()
			if len(d) > 0 {
				e := &epb.Details{}
				if err := d[0].UnmarshalTo(e); err == nil {
					for k, v := range e.Messages {
						form.FormErrors[k] = v
					}
					break
				} else {
					logging.WithError(err, log).Error("invalid details")
				}
			}
			fallthrough
		default:
			form.FormErrors["General"] = "Something went wrong."
		}
		if err := s.signupTpl.Execute(w, form); err != nil {
			log.WithError(err).Error("failed to render signup form")
			http.Error(w, genericErrMsg, http.StatusInternalServerError)
		}
		return
	}

	http.Redirect(w, r, "/user/set-password-success", http.StatusSeeOther)
}
