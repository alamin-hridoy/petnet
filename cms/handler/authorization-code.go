package handler

import (
	"html/template"
	"net/http"
	"strings"

	ppu "brank.as/petnet/gunk/dsa/v1/user"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"brank.as/rbac/gunk/v1/oauth2"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
)

type (
	authCodeTempData struct {
		Clients                    []*oauth2.Oauth2Client
		Environment                string
		UserInfo                   *User
		CSRFField                  template.HTML
		CompanyName                string
		OnboardingIncompleteStatus bool
		HasLiveAccess              bool
		ErrorMsg                   string
	}

	authCodeForm struct {
		Name              string
		ClientType        int
		CORS              string
		LogoURL           string
		Scopes            string
		RedirectURL       string
		LogoutRedirectURL string
	}

	authCodeSuccessTempData struct {
		ClientID     string
		ClientSecret string
		Environment  string
		UserInfo     *User
	}
)

func (s *Server) getAuthorizationCode(w http.ResponseWriter, r *http.Request) {
	s.getTemplateData(w, r, "")
}

func (s *Server) getTemplateData(w http.ResponseWriter, r *http.Request, errType string) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	oid := mw.GetOrgID(ctx)
	var idApiEnvType string
	var apiEnvType string
	switch goji.Param(r, apiEnv) {
	case prodEnv:
		if !s.hasLiveAccess(ctx, oid) {
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		apiEnvType = "production"
		idApiEnvType = "live"
	case sandEnv:
		apiEnvType = "sandbox"
		idApiEnvType = "sandbox"
	default:
		log.Error("unknown api environment")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	template := s.templates.Lookup("authorization-code.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	res, err := s.rbac.ListClients(ctx, &oauth2.ListClientsRequest{
		Env: idApiEnvType,
	})
	if err != nil {
		log.Error("unable to list account")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	ErrorMsg := ""
	if errType == "datainvalid" {
		ErrorMsg = "Invalid Url. Please input valid URL"
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := authCodeTempData{
		Clients:                    res.GetClients(),
		Environment:                apiEnvType,
		ErrorMsg:                   ErrorMsg,
		UserInfo:                   &usrInfo.UserInfo,
		HasLiveAccess:              s.hasLiveAccess(ctx, mw.GetOrgID(ctx)),
		CSRFField:                  csrf.TemplateField(r),
		CompanyName:                pf.GetProfile().GetBusinessInfo().CompanyName,
		OnboardingIncompleteStatus: pf.GetProfile().Status == ppb.Status_UnknownStatus,
	}

	uidd := mw.GetUserID(r.Context())
	gp, err := s.pf.GetUserProfile(r.Context(), &ppu.GetUserProfileRequest{
		UserID: uidd,
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	data.UserInfo.ProfileImage = gp.GetProfile().ProfilePicture

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postAuthorizationCode(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	var idApiEnvType string
	var apiEnvType string
	switch goji.Param(r, apiEnv) {
	case prodEnv:
		if !s.hasLiveAccess(ctx, mw.GetOrgID(ctx)) {
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		apiEnvType = "production"
		idApiEnvType = "live"
	case sandEnv:
		apiEnvType = "sandbox"
		idApiEnvType = "sandbox"
	default:
		log.Error("unknown api environment")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form authCodeForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Name, validation.Required),
		validation.Field(&form.ClientType),
		validation.Field(&form.CORS, validation.Required),
		validation.Field(&form.LogoURL, validation.Required),
		validation.Field(&form.RedirectURL, validation.Required),
		validation.Field(&form.LogoutRedirectURL, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		s.getTemplateData(w, r, "datainvalid")
		return
	}

	cors := strings.Split(strings.TrimSpace(form.CORS), ",")
	scopes := strings.Split(strings.TrimSpace(form.Scopes), ",")
	redURL := strings.Split(strings.TrimSpace(form.RedirectURL), ",")

	createClient, err := s.rbac.CreateClient(r.Context(), &oauth2.CreateClientRequest{
		Name:              form.Name,
		Audience:          idApiEnvType,
		Env:               idApiEnvType,
		ClientType:        oauth2.ClientType(form.ClientType),
		CORS:              cors,
		LogoURL:           form.LogoURL,
		Scopes:            scopes,
		RedirectURL:       redURL,
		LogoutRedirectURL: strings.TrimSpace(form.LogoutRedirectURL),
		IdentitySource:    "perahub",
		Config:            &oauth2.ClientConfig{},
	})
	if err != nil {
		log.Error("unable create client")
		s.getTemplateData(w, r, "datainvalid")
		return
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := authCodeSuccessTempData{
		ClientID:     createClient.ClientID,
		ClientSecret: createClient.Secret,
		Environment:  apiEnvType,
		UserInfo:     &usrInfo.UserInfo,
	}
	template := s.templates.Lookup("authorization-code-success.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
