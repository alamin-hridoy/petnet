package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"brank.as/petnet/cms/paginator"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	psvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
)

const limitPerPage int32 = 10

type transactionTempData struct {
	Remittances              []*tpb.Remittance
	UserInfo                 *User
	CompanyName              string
	PaginationData           paginator.Paginator
	SearchTerms              string
	OrgId                    string
	Environment              string
	HasLiveAccess            bool
	PresetPermission         map[string]map[string]bool
	ServiceRequest           bool
	PartnerListApplicantList []PartnerListApplicant
}

func (s *Server) getTransactionListSandbox(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	md := metautils.ExtractIncoming(ctx)
	oid := goji.Param(r, "id")

	if oid == "" {
		log.Error("missing id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("transaction-list.html")
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
	apiEnvVal, _ := url.PathUnescape(queryParams.Get(apiEnv))
	apiEnvType := "sandbox"
	if apiEnvVal == "production" {
		apiEnvType = "production"
	}
	var offset int32 = 0
	convertedPageNumber, _ := strconv.Atoi(pageNumber)
	if convertedPageNumber <= 0 {
		convertedPageNumber = 1
	} else {
		offset = limitPerPage*int32(convertedPageNumber) - limitPerPage
	}

	searchTerms, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	services, err := url.PathUnescape(queryParams.Get("services"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	types, err := url.PathUnescape(queryParams.Get("types"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	sbv := tpb.SortOrder_DESC
	if sb == "asc" {
		sbv = tpb.SortOrder_ASC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	currentPage, _ := strconv.Atoi(queryParams.Get("page"))
	if currentPage == 0 {
		currentPage = 1
	}
	sbcv := tpb.SortByColumn_TransactionCompletedTime
	if sbc == "remitted" {
		sbcv = tpb.SortByColumn_RemittedTo
	} else if sbc == "remittedamount" {
		sbcv = tpb.SortByColumn_TotalRemittedAmount
	} else if sbc == "dateprocessed" {
		sbcv = tpb.SortByColumn_TransactionCompletedTime
	}
	prfl, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
		return
	}
	etd := s.getEnforceTemplateData(ctx)
	profile := prfl.GetProfile()
	ctx = md.Add("x-forward-dsaorgid", profile.OrgID).ToOutgoing(ctx)
	lr := &tpb.ListRemitRequest{
		Limit:          limitPerPage,
		Offset:         offset,
		SortOrder:      sbv,
		SortByColumn:   sbcv,
		ExcludePartner: services,
		ExcludeType:    types,
	}
	if searchTerms != "" {
		lr.ControlNumbers = []string{searchTerms}
	}
	var res *tpb.ListRemitResponse
	if apiEnvType == "production" {
		res, err = s.drpLV.ListRemit(ctx, lr)
	} else {
		res, err = s.drpSB.ListRemit(ctx, lr)
	}
	if err != nil {
		logging.WithError(err, log).Error("listing remits")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	newReq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: psvc.ServiceType_REMITTANCE.String(),
	}
	svc, err := s.pf.GetPartnerList(r.Context(), newReq)
	var partnerListApplicantList []PartnerListApplicant
	if err != nil {
		log.Error("failed to Get Partner List")
	} else {
		for _, sv := range svc.GetPartnerList() {
			partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
				Stype: sv.GetStype(),
				Name:  sv.GetName(),
			})
		}
	}

	for i := range res.Remittances {
		partner := res.Remittances[i].RemitPartner
		for _, v := range partnerListApplicantList {
			if partner == v.Stype {
				partner = v.Name
			}
		}
		res.Remittances[i].RemitPartner = partner
		status := res.Remittances[i].RemitType
		switch res.Remittances[i].RemitType {
		case "SEND":
			status = "Create-Remit"
		case "DISBURSE":
			status = "Disburse-Remit"
		}
		res.Remittances[i].RemitType = status
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	tempData := transactionTempData{
		Remittances:              res.GetRemittances(),
		UserInfo:                 &usrInfo.UserInfo,
		SearchTerms:              searchTerms,
		OrgId:                    oid,
		Environment:              apiEnvType,
		HasLiveAccess:            s.hasLiveAccess(r.Context(), oid),
		CompanyName:              profile.GetBusinessInfo().GetCompanyName(),
		PresetPermission:         etd.PresetPermission,
		ServiceRequest:           etd.ServiceRequests,
		PartnerListApplicantList: partnerListApplicantList,
	}
	if len(tempData.Remittances) > 0 {
		tempData.PaginationData = paginator.NewPaginator(int32(currentPage), limitPerPage, res.Total, r)
	}

	tempData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
