package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"brank.as/petnet/cms/paginator"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	rbupb "brank.as/rbac/gunk/v1/user"
)

type (
	TransactionApplicant struct {
		OrgID       string
		CompanyName string
		User        User
	}

	transactionTemplateData struct {
		TransactionApplicants []TransactionApplicant
		PaginationData        paginator.Paginator
		UserInfo              *User
		SearchTerm            string
		PresetPermission      map[string]map[string]bool
		ServiceRequest        bool
	}
)

const transactionCompanyPerPage int32 = 10

func (s *Server) getTransactionCompanies(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	template := s.templates.Lookup("transaction-list-of-companies.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	pageNumber, err := url.PathUnescape(queryParams.Get("page"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	var offset int32 = 0
	convertedPageNumber, _ := strconv.Atoi(pageNumber)

	if convertedPageNumber <= 0 {
		convertedPageNumber = 1
	} else {
		offset = perPage*int32(convertedPageNumber) - perPage
	}

	searchTerm, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	newReq := &ppb.ListProfilesRequest{
		Limit:       transactionCompanyPerPage,
		Offset:      offset,
		CompanyName: searchTerm,
	}

	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	if sb == "asc" {
		newReq.SortBy = ppb.SortBy_ASC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	if sbc == "company" {
		newReq.SortByColumn = ppb.SortByColumn_CompanyName
	}
	newReq.OrgType = ppb.OrgType_DSA
	pf, err := s.pf.ListProfiles(r.Context(), newReq)
	if err != nil {
		log.Error("getting profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	currentPage, _ := strconv.Atoi(queryParams.Get("page"))
	if currentPage == 0 {
		currentPage = 1
	}

	var transactionAplicantList []TransactionApplicant
	for _, applicant := range pf.Profiles {
		if applicant.UserID != "" {
			u, err := s.rbac.GetUser(r.Context(), &rbupb.GetUserRequest{ID: applicant.UserID})
			if err != nil {
				logging.WithError(err, log).Error("getting user")
			} else {
				transactionAplicantList = append(transactionAplicantList, TransactionApplicant{
					OrgID:       applicant.OrgID,
					CompanyName: applicant.BusinessInfo.CompanyName,
					User:        userData(u.GetUser()),
				})
			}
		}
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(r.Context())
	templateData := transactionTemplateData{
		TransactionApplicants: transactionAplicantList,
		PaginationData:        paginator.Paginator{},
		UserInfo:              &usrInfo.UserInfo,
		SearchTerm:            searchTerm,
		PresetPermission:      etd.PresetPermission,
		ServiceRequest:        etd.ServiceRequests,
	}
	if len(transactionAplicantList) > 0 {
		templateData.PaginationData = paginator.NewPaginator(int32(currentPage), perPage, pf.Total, r)
	}

	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
