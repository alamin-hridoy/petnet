package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
)

type (
	currencyTemplateData struct {
		CSRFField                  template.HTML
		OrgID                      string
		Status                     string
		AllDocsSubmitted           bool
		BusinessInfo               *ppb.BusinessInfo
		AccountInfo                *ppb.AccountInfo
		SelectedCurrency           int32
		CurrencyList               map[int32]string
		Errors                     map[string]error
		PresetPermission           map[string]map[string]bool
		User                       *User
		ServiceRequest             bool
		DsaCode                    string
		TerminalIdOtc              string
		TerminalIdDigital          string
		TransactionTypesForDSACode TransactionTypesForDSACode
	}

	CurrencyForm struct {
		Bank              string
		BankAccountNumber string
		BankAccountHolder string
		Currency          string
	}
)

func (s *Server) getCurrencyConfig(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	template := s.templates.Lookup("currency-configuration.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid := goji.Param(r, "id")
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: string(oid)})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	loadProfile := pf.GetProfile()
	status := loadProfile.GetStatus().String()
	if loadProfile.GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	if loadProfile.GetStatus() == ppb.Status_PendingDocuments {
		status = "Pending Documents"
	}

	fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: string(oid)})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	transactionTypesForDSACode := checkTrnTypeForDSACode(loadProfile.GetTransactionTypes())
	etd := s.getEnforceTemplateData(r.Context())
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := currencyTemplateData{
		OrgID:                      oid,
		Status:                     status,
		AllDocsSubmitted:           allDocsSubmitted(fs.GetFileUploads()),
		CSRFField:                  csrf.TemplateField(r),
		BusinessInfo:               loadProfile.GetBusinessInfo(),
		AccountInfo:                loadProfile.GetAccountInfo(),
		CurrencyList:               ppb.Currency_name,
		PresetPermission:           etd.PresetPermission,
		ServiceRequest:             etd.ServiceRequests,
		User:                       &usrInfo.UserInfo,
		DsaCode:                    loadProfile.DsaCode,
		TerminalIdOtc:              loadProfile.TerminalIdOtc,
		TerminalIdDigital:          loadProfile.TerminalIdDigital,
		TransactionTypesForDSACode: transactionTypesForDSACode,
	}
	if data.AccountInfo.Currency != ppb.Currency_UnknownCurrency {
		data.SelectedCurrency = ppb.Currency_value[data.AccountInfo.Currency.String()]
	}

	data.User.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (cf CurrencyForm) validate() error {
	return validation.ValidateStruct(&cf,
		validation.Field(&cf.Currency,
			validation.Required.Error("This field is required"),
		),
	)
}

func (s *Server) changeCurrencyConfig(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	orgId := goji.Param(r, "id")

	var form CurrencyForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	vErr := form.validate()
	if vErr != nil {
		if vErrs, ok := vErr.(validation.Errors); ok {
			template := s.templates.Lookup("currency-configuration.html")
			if template == nil {
				log.Error("unable to load template")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			oid := goji.Param(r, "id")
			pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: string(oid)})
			if err != nil {
				logging.WithError(err, log).Info("getting profile")
			}
			etd := s.getEnforceTemplateData(r.Context())
			loadProfile := pf.GetProfile()
			data := currencyTemplateData{
				OrgID:            orgId,
				CSRFField:        csrf.TemplateField(r),
				AccountInfo:      loadProfile.GetAccountInfo(),
				Errors:           vErrs,
				PresetPermission: etd.PresetPermission,
				ServiceRequest:   etd.ServiceRequests,
			}
			if err := template.Execute(w, data); err != nil {
				log.Infof("error with template execution: %+v", err)
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			return
		}
	}
	c, err := strconv.Atoi(form.Currency)
	if err != nil {
		fmt.Println(err)
	}
	if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
		Profile: &ppb.OrgProfile{
			OrgID: orgId,
			AccountInfo: &ppb.AccountInfo{
				Currency: ppb.Currency(c),
			},
		},
	}); err != nil {
		logging.WithError(err, log).Error("creating profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/dashboard/currency/%s", orgId), http.StatusSeeOther)
}
