package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"brank.as/petnet/cms/paginator"
	"brank.as/petnet/gunk/dsa/v2/partnerlist"
	svr "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type (
	DSASvcApplicant struct {
		CSRFField   template.HTML
		CompanyName string
		ServiceName string
		PartName    string
		Status      string
		UpatedBy    string
		Remark      string
		ID          string
		OrgID       string
		Updated     time.Time
	}

	DSASvcTemplateData struct {
		CSRFFieldValue   string
		OrgID            string
		CSRFField        template.HTML
		UserInfo         *User
		DSASvcApplicants []DSASvcApplicant
		PaginationData   paginator.Paginator
		SearchTerm       string
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
	}

	AddRemarkForm struct {
		OrgID       string
		Remark      string
		ServiceName string
		UpatedBy    string
		CSRFField   template.HTML
	}
)

func (s *Server) getServiceReq(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	template := s.templates.Lookup("service-req.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
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
	newReq := &svr.GetAllServiceRequestRequest{
		Limit:    perPage,
		Offset:   offset,
		Statuses: []svr.ServiceRequestStatus{svr.ServiceRequestStatus_ACCEPTED, svr.ServiceRequestStatus_PENDING, svr.ServiceRequestStatus_REJECTED},
	}

	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if sb == "asc" {
		newReq.SortBy = svr.SortBy_ASC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if sbc == "company" {
		newReq.SortByColumn = svr.ServiceSort_COMPANYNAME
	}

	if sbc == "partner" {
		newReq.SortByColumn = svr.ServiceSort_PARTNER
	}

	if sbc == "date" {
		newReq.SortByColumn = svr.ServiceSort_LASTUPDATED
	}

	svcnm, err := url.PathUnescape(queryParams.Get("service_name"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if svcnm != "" {
		svcnmArr := strings.Split(svcnm, ",")
		svArr := []svr.ServiceType{}
		for _, rs := range svcnmArr {
			switch rs {
			case "rem":
				svArr = append(svArr, svr.ServiceType_REMITTANCE)
			case "bp":
				svArr = append(svArr, svr.ServiceType_BILLSPAYMENT)
			case "cio":
				svArr = append(svArr, svr.ServiceType_CASHINCASHOUT)
			case "mi":
				svArr = append(svArr, svr.ServiceType_MICROINSURANCE)
			}
		}

		newReq.Types = svArr
	}

	if searchTerm != "" {
		newReq.CompanyName = searchTerm
	}

	sts, err := url.PathUnescape(queryParams.Get("status"))
	if err != nil {
		log.Error("unable to decode url type param")
	}
	if sts != "" {
		sStrArr := strings.Split(sts, ",")
		sArr := []svr.ServiceRequestStatus{}
		for _, s := range sStrArr {
			switch s {
			case "accepted":
				sArr = append(sArr, svr.ServiceRequestStatus_ACCEPTED)
			case "rejected":
				sArr = append(sArr, svr.ServiceRequestStatus_REJECTED)
			case "pending":
				sArr = append(sArr, svr.ServiceRequestStatus_PENDING)
			}
		}
		newReq.Statuses = sArr
	}
	etd := s.getEnforceTemplateData(ctx)
	lSvcR, _ := s.pf.GetAllServiceRequest(ctx, newReq)
	userInfoLists := map[string]*rbupb.User{}
	if len(lSvcR.GetServiceRequst()) > 0 {
		userLists := []string{}
		for _, lv := range lSvcR.GetServiceRequst() {
			if lv.UpdatedBy != "" {
				userLists = append(userLists, lv.UpdatedBy)
			}
		}
		userLists = uniqueSlice(userLists)
		usrs, err := s.rbac.ListUsers(ctx, &rbupb.ListUsersRequest{
			ID:    userLists,
			OrgID: mw.GetOrgID(ctx),
		})
		if err == nil {
			userInfoLists = usrs.GetUser()
		}
	}
	var sList []DSASvcApplicant
	for _, v := range lSvcR.GetServiceRequst() {
		sts := v.GetStatus().String()
		psts := "ENABLED"
		gpLst, err := s.pf.GetPartnerList(ctx, &partnerlist.GetPartnerListRequest{
			Stype:       strings.TrimSpace(v.GetPartners()),
			Status:      "ENABLED",
			ServiceName: v.GetType().String(),
		})
		if err != nil {
			psts = "DISABLED"
		}
		if gpLst == nil || len(gpLst.GetPartnerList()) == 0 {
			psts = "DISABLED"
		}
		if psts == "DISABLED" && sts == "ACCEPTED" {
			sts = "DISABLED"
		}

		sList = append(sList, DSASvcApplicant{
			CompanyName: v.GetCompanyName(),
			ServiceName: v.GetType().String(),
			PartName:    v.GetPartner(),
			Status:      sts,
			UpatedBy:    userInfoLists[v.GetUpdatedBy()].GetFirstName() + " " + userInfoLists[v.GetUpdatedBy()].GetLastName(),
			Remark:      v.GetRemarks(),
			ID:          v.GetID(),
			OrgID:       v.GetOrgID(),
			Updated:     v.GetUpdated().AsTime(),
		})
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	templateData := DSASvcTemplateData{
		OrgID:            mw.GetOrgID(ctx),
		CSRFFieldValue:   csrf.Token(r),
		CSRFField:        csrf.TemplateField(r),
		UserInfo:         &usrInfo.UserInfo,
		DSASvcApplicants: sList,
		PaginationData:   paginator.Paginator{},
		SearchTerm:       searchTerm,
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
	}
	if len(sList) > 0 {
		templateData.PaginationData = paginator.NewPaginator(int32(convertedPageNumber), perPage, lSvcR.Total, r)
	}
	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postAddRemarkSvcReq(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	var form AddRemarkForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := map[string]string{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.OrgID, validation.Required),
		validation.Field(&form.ServiceName, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["ID"] != nil {
				formErr["ID"] = err["ID"].Error()
			}
		}
	}
	if len(formErr) > 0 {
		log.Error("ID is required to add remark")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if _, err := s.pf.AddRemarkSvcRequest(ctx, &svr.AddRemarkSvcRequestRequest{
		OrgID:     form.OrgID,
		SvcName:   form.ServiceName,
		Remark:    form.Remark,
		UpdatedBy: mw.GetUserID(ctx),
	}); err != nil {
		logging.WithError(err, log).Error("add remark failed in SVC request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, serviceReqPath, http.StatusSeeOther)
}
