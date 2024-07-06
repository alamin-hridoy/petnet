package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"brank.as/petnet/cms/paginator"
	"brank.as/petnet/cms/storage"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
)

type (
	DSAInviteApplicant struct {
		CSRFField            template.HTML
		ID                   string
		OrgID                string
		CompanyName          string
		DateApplied          time.Time
		RiskScore            string
		Status               string
		Partner              string
		ReminderSent         bool
		IsDocumentsSubmitted bool
		User                 User
	}

	ProvidersListTempData struct {
		CSRFField           template.HTML
		CSRFFieldValue      string
		SearchTerm          string
		InviteUsrErr        string
		PresetPermission    map[string]map[string]bool
		ServiceRequest      bool
		LoginUserInfo       *User
		UserInfo            *User
		PaginationData      paginator.Paginator
		DSAApplicants       []DSAInviteApplicant
		PartnerList         []*spbl.PartnerList
		FirstName           string
		LastName            string
		Email               string
		Provider            string
		TransactionTypes    []string
		HasProviderFrmError bool
		InviteProviderErr   string
		Errors              map[string]error
	}

	InviteProviderForm struct {
		FirstName        string
		LastName         string
		Email            string
		Provider         string
		TransactionTypes []string
	}
)

func (s *Server) getProvidersList(w http.ResponseWriter, r *http.Request) {
	s.manageProvidersListForm(w, r, nil)
}

func (s *Server) manageProvidersListForm(w http.ResponseWriter, r *http.Request, formErr map[string]error) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	template := s.templates.Lookup("providers-list.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var f InviteProviderForm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
	}

	queryParams := r.URL.Query()
	pageNumber, err := url.PathUnescape(queryParams.Get("page"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	searchTerm, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	var offset int32 = 0
	convertedPageNumber, _ := strconv.Atoi(pageNumber)

	if convertedPageNumber <= 0 {
		convertedPageNumber = 1
	} else {
		offset = perPage*int32(convertedPageNumber) - perPage
	}

	newReq := &ppb.ListProfilesRequest{
		Limit:       perPage,
		Offset:      offset,
		CompanyName: searchTerm,
		IsProvider:  true,
	}

	ptnrs, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		IsProvider: true,
	})
	if err != nil {
		log.Error("unable to get partner list")
	}

	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if sb == "asc" {
		newReq.SortBy = ppb.SortBy_ASC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if sbc == "company" {
		newReq.SortByColumn = ppb.SortByColumn_CompanyName
	}

	sts, err := url.PathUnescape(queryParams.Get("status"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if sts != "" {
		sStrArr := strings.Split(sts, ",")
		sArr := []ppb.Status{}
		for _, s := range sStrArr {
			switch s {
			case "accepted":
				sArr = append(sArr, ppb.Status_Accepted)
			case "completed":
				sArr = append(sArr, ppb.Status_Completed)
			case "pending":
				sArr = append(sArr, ppb.Status_Pending)
			case "rejected":
				sArr = append(sArr, ppb.Status_Rejected)
			case "pending-document":
				sArr = append(sArr, ppb.Status_PendingDocuments)
			case "incomplete":
				sArr = append(sArr, ppb.Status_UnknownStatus)
			}
		}
		newReq.Status = sArr
	}

	inviteProviderErr := ""
	errorMsg, _ := url.PathUnescape(queryParams.Get("errorMsg"))
	switch errorMsg {
	case "AlreadyExists":
		inviteProviderErr = "Email address already exists."
	case "CreateProfile":
		inviteProviderErr = "Failed to create org profile."
	case "ServiceReq":
		inviteProviderErr = "Failed to create service request."
	case "CreateUserProfile":
		inviteProviderErr = "Failed to create user profile."
	case "GetPartnerByStype":
		inviteProviderErr = "Failed to get provider service type."
	case "InviteUser":
		inviteProviderErr = "Failed to send invite."
	case "AddUser":
		inviteProviderErr = "Failed to add user to role."
	case "RoleNotFound":
		inviteProviderErr = "provider role not found. please create a role as name 'provider'."
	}

	ptnrList := ptnrs.GetPartnerList()
	newReq.OrgType = ppb.OrgType_DSA
	pf, err := s.pf.ListProfiles(r.Context(), newReq)
	formErrs := validation.Errors{}
	hasErr := false
	if formErr != nil {
		hasErr = true
		formErrs = formErr
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(ctx)
	tData := ProvidersListTempData{
		CSRFField:           csrf.TemplateField(r),
		SearchTerm:          searchTerm,
		PresetPermission:    etd.PresetPermission,
		ServiceRequest:      etd.ServiceRequests,
		LoginUserInfo:       &usrInfo.UserInfo,
		CSRFFieldValue:      csrf.Token(r),
		UserInfo:            &usrInfo.UserInfo,
		PartnerList:         ptnrList,
		HasProviderFrmError: hasErr,
		Errors:              formErrs,
		FirstName:           f.FirstName,
		LastName:            f.LastName,
		Email:               f.Email,
		Provider:            f.Provider,
		TransactionTypes:    f.TransactionTypes,
		InviteProviderErr:   inviteProviderErr,
	}
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
		if err := template.Execute(w, tData); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	currentPage, _ := strconv.Atoi(queryParams.Get("page"))
	if currentPage == 0 {
		currentPage = 1
	}

	var dsaAplicantList []DSAInviteApplicant
	for _, applicant := range pf.Profiles {
		u, _ := s.rbac.GetUser(r.Context(), &rbupb.GetUserRequest{ID: applicant.UserID})
		var udata User
		if u != nil {
			udata = userData(u.GetUser())
		}
		status := applicant.GetStatus().String()
		if applicant.GetStatus() == ppb.Status_UnknownStatus {
			status = "Incomplete"
		}
		convertedTime := applicant.DateApplied.AsTime()
		dsaAplicantList = append(dsaAplicantList, DSAInviteApplicant{
			CSRFField:    csrf.TemplateField(r),
			ID:           applicant.ID,
			OrgID:        applicant.OrgID,
			CompanyName:  applicant.BusinessInfo.CompanyName,
			DateApplied:  convertedTime,
			RiskScore:    applicant.RiskScore.String(),
			Status:       status,
			ReminderSent: applicant.GetReminderSent() == ppb.Boolean_True,
			User:         udata,
			Partner:      applicant.Partner,
		})
	}

	templateData := ProvidersListTempData{
		CSRFField:           csrf.TemplateField(r),
		SearchTerm:          searchTerm,
		PresetPermission:    etd.PresetPermission,
		ServiceRequest:      etd.ServiceRequests,
		LoginUserInfo:       &usrInfo.UserInfo,
		CSRFFieldValue:      csrf.Token(r),
		DSAApplicants:       dsaAplicantList,
		UserInfo:            &usrInfo.UserInfo,
		PartnerList:         ptnrList,
		HasProviderFrmError: hasErr,
		Errors:              formErrs,
		FirstName:           f.FirstName,
		LastName:            f.LastName,
		Email:               f.Email,
		Provider:            f.Provider,
		TransactionTypes:    f.TransactionTypes,
		InviteProviderErr:   inviteProviderErr,
	}
	if len(dsaAplicantList) > 0 {
		templateData.PaginationData = paginator.NewPaginator(int32(currentPage), perPage, pf.Total, r)
	}

	templateData.LoginUserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, templateData); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postInviteProvides(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var f InviteProviderForm
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&f,
		validation.Field(&f.FirstName, validation.Required),
		validation.Field(&f.LastName, validation.Required),
		validation.Field(&f.Email, validation.Required, is.Email),
		validation.Field(&f.Provider, validation.Required),
		validation.Field(&f.TransactionTypes, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["FirstName"] != nil {
				formErr["FirstName"] = err["FirstName"]
			}
			if err["LastName"] != nil {
				formErr["LastName"] = err["LastName"]
			}
			if err["Email"] != nil {
				formErr["Email"] = err["Email"]
			}
			if err["Provider"] != nil {
				formErr["Provider"] = err["Provider"]
			}
			if err["TransactionTypes"] != nil {
				formErr["TransactionTypes"] = err["TransactionTypes"]
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}

	if len(formErr) > 0 {
		s.manageProvidersListForm(w, r, formErr)
		return
	}

	ipf := &InviteProviderForm{
		FirstName:        f.FirstName,
		LastName:         f.LastName,
		Email:            f.Email,
		Provider:         f.Provider,
		TransactionTypes: f.TransactionTypes,
	}

	d, err := json.Marshal(ipf)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	me := &mfaEvent{
		resource: string(storage.InviteProvider),
		action:   tpb.ActionType_Create,
		data:     d,
	}
	if err := s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, providersListPath, http.StatusSeeOther)
	}
	http.Redirect(w, r, providersListPath+"?show_otp=true", http.StatusSeeOther)
}
