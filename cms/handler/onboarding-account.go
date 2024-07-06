package handler

import (
	"errors"
	"html/template"
	"net/http"
	"regexp"
	"strings"

	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
)

type (
	AccountInfoForm struct {
		CSRFField               template.HTML
		Bank                    string
		BankAccountNumber       string
		BankAccountHolder       string
		OnlineSupplierCondition int
		TermsCondition          int
		Errors                  map[string]error
		PresetPermission        map[string]map[string]bool
		ServiceRequest          bool
		SaveDraft               string
		SaveContinue            string
	}
)

func (s *Server) getAccountInfo(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	uid := sess.Values[sessionUserID]
	if uid == nil {
		log.Error("missing user id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := sess.Values[sessionOrgID]
	if oid == nil {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("onboarding-account.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid.(string)})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	form := aiProtoToForm(pf.GetProfile().GetAccountInfo())
	form.CSRFField = csrf.TemplateField(r)
	etd := s.getEnforceTemplateData(r.Context())
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	if err := template.Execute(w, form); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postAccountInfo(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	uid, ok := sess.Values[sessionUserID].(string)
	if !ok {
		log.Error("missing user id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid, ok := sess.Values[sessionOrgID].(string)
	if !ok {
		log.Error("missing org id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logging.WithError(err, log).Error("parsing form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form AccountInfoForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Bank, validation.Required),
		validation.Field(&form.BankAccountNumber, is.Digit, validation.Required),
		validation.Field(&form.BankAccountHolder, validation.Required.Error("Bank Account holder is required"), validation.Match(regexp.MustCompile("^[-a-zA-Z0-9-()]+(\\s+[-a-zA-Z0-9-()]+)*$")).Error("Space is not allowed before and after the name.")),
		validation.Field(&form.OnlineSupplierCondition, validation.Required),
		validation.Field(&form.TermsCondition, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Bank"] != nil {
				formErr["Bank"] = errors.New("Bank is required")
			}
			if err["BankAccountNumber"] != nil {
				formErr["BankAccountNumber"] = errors.New("Please put valid bank account number")
			}
			if err["BankAccountHolder"] != nil {
				formErr["BankAccountHolder"] = err["BankAccountHolder"]
			}
			if err["OnlineSupplierCondition"] != nil {
				formErr["OnlineSupplierCondition"] = errors.New("")
			}
			if err["TermsCondition"] != nil {
				formErr["TermsCondition"] = errors.New("")
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}

	template := s.templates.Lookup("onboarding-account.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	etd := s.getEnforceTemplateData(ctx)
	form.PresetPermission = etd.PresetPermission
	form.ServiceRequest = etd.ServiceRequests
	if len(formErr) > 0 {
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	if form.SaveDraft == "SaveDraft" {
		if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
			Profile: &ppb.OrgProfile{
				UserID: uid,
				OrgID:  oid,
				AccountInfo: &ppb.AccountInfo{
					Bank:                    strings.TrimSpace(form.Bank),
					BankAccountNumber:       strings.TrimSpace(form.BankAccountNumber),
					BankAccountHolder:       strings.TrimSpace(form.BankAccountHolder),
					AgreeTermsConditions:    ppb.Boolean(form.TermsCondition),
					AgreeOnlineSupplierForm: ppb.Boolean(form.OnlineSupplierCondition),
				},
			},
		}); err != nil {
			logging.WithError(err, log).Error("creating profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, loginPath, http.StatusSeeOther)
	}

	if form.SaveContinue == "SaveContinue" {
		if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
			Profile: &ppb.OrgProfile{
				UserID: uid,
				OrgID:  oid,
				AccountInfo: &ppb.AccountInfo{
					Bank:                    strings.TrimSpace(form.Bank),
					BankAccountNumber:       strings.TrimSpace(form.BankAccountNumber),
					BankAccountHolder:       strings.TrimSpace(form.BankAccountHolder),
					AgreeTermsConditions:    ppb.Boolean(form.TermsCondition),
					AgreeOnlineSupplierForm: ppb.Boolean(form.OnlineSupplierCondition),
				},
			},
		}); err != nil {
			logging.WithError(err, log).Error("creating profile")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, drpInfoPath, http.StatusSeeOther)
	}
	http.Redirect(w, r, drpInfoPath, http.StatusSeeOther)
}

func aiProtoToForm(bip *ppb.AccountInfo) AccountInfoForm {
	return AccountInfoForm{
		Bank:                    bip.GetBank(),
		BankAccountNumber:       bip.GetBankAccountNumber(),
		BankAccountHolder:       bip.GetBankAccountHolder(),
		TermsCondition:          int(bip.GetAgreeTermsConditions()),
		OnlineSupplierCondition: int(bip.GetAgreeOnlineSupplierForm()),
	}
}
