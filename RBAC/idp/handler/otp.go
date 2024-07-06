package handler

import (
	"html/template"
	"net/http"

	"github.com/gorilla/csrf"

	"brank.as/rbac/serviceutil/logging"
)

type OTPForm struct {
	Code            string
	Email           string
	CSRFField       template.HTML
	CodeError       string
	ProcessEndpoint string
	Retry           string
}

type otpTemplateData struct {
	OTPForm OTPForm
	urls
	ProjectName string
}

func (s *server) getOTPHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("method", "getOTPHandler")

	data := otpTemplateData{
		OTPForm: OTPForm{
			CSRFField:       csrf.TemplateField(r),
			ProcessEndpoint: "/otp",
		},
		urls:        s.u,
		ProjectName: s.projectName,
	}
	if err := s.otpTpl.Execute(w, data); err != nil {
		errMsg := "failed to render OTP form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusInternalServerError)
		return
	}
}
