package handler

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/url"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	phmw "brank.as/petnet/api/perahub-middleware"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	ttpb "brank.as/petnet/gunk/dsa/v2/transactiontype"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbmpb "brank.as/rbac/gunk/v1/mfa"
	rbsapb "brank.as/rbac/gunk/v1/serviceaccount"
	rbusr "brank.as/rbac/gunk/v1/user"
)

const (
	apiEnv  = "apienv"
	prodEnv = "production"
	sandEnv = "sandbox"
)

type (
	ApiKeyPasswordForm struct {
		CSRFField               template.HTML
		Environment             string
		TransactionType         string
		SelectedEnvironment     string
		SelectedTransactionType string
		Name                    string
		Status                  string
		Password                string
		Email                   string
		UserInfo                *User
		InvalidPassword         bool
		InvlidTransactionType   bool
		InvalidEnvironment      bool
		InvalidEmail            bool
		ServiceRequest          bool
		HasLiveAccess           bool
		CreateApiError          bool
		CompanyName             string
		PresetPermission        map[string]map[string]bool
	}

	apikeyTempData struct {
		AccountApiKeys             []*ttpb.ApiKeyTransactionType
		UserInfo                   *User
		Environment                string
		CompanyName                string
		OnboardingIncompleteStatus bool
		HasLiveAccess              bool
		ServiceRequest             bool
		PresetPermission           map[string]map[string]bool
	}

	apikeySuccessTempData struct {
		ClientID         string
		ClientSecret     string
		Environment      string
		UserInfo         *User
		HasLiveAccess    bool
		ServiceRequest   bool
		CompanyName      string
		PresetPermission map[string]map[string]bool
	}
)

func (s *Server) getApiKeyList(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	oid := mw.GetOrgID(ctx)
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	var apiEnvType string
	switch goji.Param(r, apiEnv) {
	case prodEnv:
		if !s.hasLiveAccess(ctx, mw.GetOrgID(ctx)) {
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		apiEnvType = "production"
	case sandEnv:
		apiEnvType = "sandbox"
	default:
		log.Error("unknown api environment")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("api-key-list.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	res, err := s.pf.ListUserAPIKeyTransactionType(ctx, &ttpb.ListUserAPIKeyTransactionTypeRequest{
		OrgID:  oid,
		UserID: uid,
	})
	if err != nil {
		logging.WithError(err, log).Error("unable to list account")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	etd := s.getEnforceTemplateData(r.Context())
	data := apikeyTempData{
		AccountApiKeys:             res.GetApiTransactionType(),
		UserInfo:                   &usrInfo.UserInfo,
		Environment:                apiEnvType,
		HasLiveAccess:              s.hasLiveAccess(ctx, mw.GetOrgID(ctx)),
		CompanyName:                pf.GetProfile().GetBusinessInfo().CompanyName,
		OnboardingIncompleteStatus: pf.GetProfile().Status == ppb.Status_UnknownStatus,
		ServiceRequest:             etd.ServiceRequests,
		PresetPermission:           etd.PresetPermission,
	}
	data.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) getApiKeyGenerate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	apiEnvType := goji.Param(r, apiEnv)
	if apiEnvType != prodEnv && apiEnvType != sandEnv {
		log.Error("unknown api environment")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if apiEnvType == prodEnv && !s.hasLiveAccess(ctx, mw.GetOrgID(ctx)) {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid := mw.GetOrgID(ctx)
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	transactiontype := pf.Profile.TransactionTypes
	profilestatus := pf.Profile.Status
	etd := s.getEnforceTemplateData(r.Context())
	var form ApiKeyPasswordForm
	CreateApiError := r.URL.Query().Get("CreateApiError")
	if CreateApiError == "true" {
		form.CreateApiError = true
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	form.CSRFField = csrf.TemplateField(r)
	form.Environment = apiEnvType
	form.TransactionType = transactiontype
	form.Status = profilestatus.String()
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	form.HasLiveAccess = s.hasLiveAccess(ctx, mw.GetOrgID(ctx))
	form.UserInfo = &usrInfo.UserInfo
	form.CompanyName = pf.GetProfile().GetBusinessInfo().CompanyName
	form.UserInfo.ProfileImage = usrInfo.ProfileImage
	s.loadAPIKeyGenerateForm(w, r, form)
}

func (s *Server) postApiKeyGenerate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	uid := mw.GetUserID(ctx)
	apiEnvType := goji.Param(r, apiEnv)
	if apiEnvType != prodEnv && apiEnvType != sandEnv {
		log.Error("unknown api environment")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if apiEnvType == prodEnv && !s.hasLiveAccess(ctx, mw.GetOrgID(ctx)) {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := mw.GetOrgID(ctx)
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	transactiontype := pf.Profile.TransactionTypes
	profilestatus := pf.Profile.Status
	var form ApiKeyPasswordForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	txnType := form.TransactionType
	txnEnv := form.Environment
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Email, validation.Required),
		validation.Field(&form.Password, validation.Required),
		validation.Field(&form.TransactionType, validation.Required),
		validation.Field(&form.Environment, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		if err, ok := (err).(validation.Errors); ok {
			if err["Email"] != nil {
				form.InvalidEmail = true
			}
			if err["Password"] != nil {
				form.InvalidPassword = true
			}
			if err["TransactionType"] != nil {
				form.InvlidTransactionType = true
			}
			if err["Environment"] != nil {
				form.InvalidEnvironment = true
			}

		}
		form.SelectedTransactionType = txnType
		form.SelectedEnvironment = txnEnv
		form.TransactionType = transactiontype
		form.Status = profilestatus.String()
		form.HasLiveAccess = s.hasLiveAccess(ctx, mw.GetOrgID(ctx))
		s.loadAPIKeyGenerateForm(w, r, form)
		return
	}
	md := metautils.ExtractIncoming(ctx)
	ctx = md.Set(phmw.TransactionType, form.TransactionType).ToIncoming(ctx)

	hasError := false
	lu, err := s.rbac.GetUser(ctx, &rbusr.GetUserRequest{
		ID: uid,
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if lu.User.GetEmail() != form.Email {
		hasError = true
		logging.WithError(err, log).Error("invalid email")
		form.InvalidEmail = true
	}
	_, err = s.rbac.ValidateMFA(ctx, &rbmpb.ValidateMFARequest{
		UserID: uid,
		Type:   rbmpb.MFA_PASS,
		Token:  form.Password,
	})
	if err != nil {
		hasError = true
		logging.WithError(err, log).Error("incorrect password")
		form.InvalidPassword = true
	}
	if hasError {
		form.HasLiveAccess = s.hasLiveAccess(ctx, mw.GetOrgID(ctx))
		form.TransactionType = transactiontype
		form.SelectedTransactionType = txnType
		form.SelectedEnvironment = txnEnv
		form.Status = profilestatus.String()
		s.loadAPIKeyGenerateForm(w, r, form)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/api-key/generate/%s/success?Email=%s&TransactionType=%s", apiEnvType, url.QueryEscape(form.Email), form.TransactionType), http.StatusSeeOther)
}

func (s *Server) getApiKeySuccess(w http.ResponseWriter, r *http.Request) {
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
	Email := r.URL.Query().Get("Email")
	if Email == "" {
		log.Error("missing api-key Email")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	TransactionType := r.URL.Query().Get("TransactionType")
	if TransactionType == "" {
		log.Error("missing api-key TransactionType")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("api-key-success.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := mw.GetUserID(ctx)
	oid := mw.GetOrgID(ctx)
	clientID := ""
	secret := ""

	getApi, err := s.pf.GetAPITransactionType(ctx, &ttpb.GetAPITransactionTypeRequest{
		UserID:          uid,
		OrgID:           oid,
		Environment:     apiEnvType,
		TransactionType: ttpb.TransactionType(ttpb.TransactionType_value[TransactionType]),
	})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			createAccount, err := s.rbac.CreateAccount(
				ctx, &rbsapb.CreateAccountRequest{
					Name:     Email,
					Env:      idApiEnvType,
					AuthType: rbsapb.AuthType_OAuth2,
				})
			if err != nil {
				log.Error("unable create account")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			clientID = createAccount.ClientID
			secret = createAccount.Secret
			var form ApiKeyPasswordForm
			if err := s.decoder.Decode(&form, r.PostForm); err != nil {
				logging.WithError(err, log).Error("decoding form")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}

			_, err = s.pf.CreateApiKeyTransactionType(ctx, &ttpb.CreateApiKeyTransactionTypeRequest{
				UserID:          uid,
				OrgID:           oid,
				ClientID:        createAccount.ClientID,
				Environment:     apiEnvType,
				TransactionType: ttpb.TransactionType(ttpb.TransactionType_value[TransactionType]),
			})

			if err != nil {
				log.Error("unable create api")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
		}
		if status.Convert(err).Code() != codes.NotFound {
			log.Error("unable get api")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	if getApi != nil {
		log.Error("unable get api")
		http.Redirect(w, r, fmt.Sprintf("/api-key/generate/%s?CreateApiError=%s", apiEnvType, "true"), http.StatusSeeOther)
		return
	}
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	etd := s.getEnforceTemplateData(r.Context())
	data := apikeySuccessTempData{
		ClientID:         clientID,
		ClientSecret:     secret,
		Environment:      apiEnvType,
		HasLiveAccess:    s.hasLiveAccess(ctx, mw.GetOrgID(ctx)),
		UserInfo:         &usrInfo.UserInfo,
		ServiceRequest:   etd.ServiceRequests,
		CompanyName:      pf.GetProfile().GetBusinessInfo().CompanyName,
		PresetPermission: etd.PresetPermission,
	}
	data.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) loadAPIKeyGenerateForm(w http.ResponseWriter, r *http.Request, form ApiKeyPasswordForm) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	template := s.templates.Lookup("api-key-input-password.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	apiEnvType := goji.Param(r, apiEnv)
	if apiEnvType != prodEnv && apiEnvType != sandEnv {
		log.Error("unknown api environment")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if apiEnvType == prodEnv && !s.hasLiveAccess(ctx, mw.GetOrgID(ctx)) {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	etd := s.getEnforceTemplateData(r.Context())
	form.ServiceRequest = etd.ServiceRequests
	form.PresetPermission = etd.PresetPermission
	form.CSRFField = csrf.TemplateField(r)
	form.Environment = apiEnvType
	form.UserInfo = &usrInfo.UserInfo
	form.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, form); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) hasLiveAccess(ctx context.Context, oid string) bool {
	log := logging.FromContext(ctx)
	res, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).WithField("org_id", oid).Info("unable to get profile")
		return false
	}

	if res.GetProfile().GetStatus() == ppb.Status_Accepted {
		return true
	}

	return false
}
