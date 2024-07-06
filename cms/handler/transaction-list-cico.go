package handler

import (
	"net/http"
	"net/url"
	"strconv"

	"brank.as/petnet/cms/paginator"
	cpb "brank.as/petnet/gunk/drp/v1/cashincashout"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	psvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
)

type transactionCICOTempData struct {
	CICOTransacts            []*cpb.CICOTransact
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

func (s *Server) getTransactionListCICOSandbox(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	md := metautils.ExtractIncoming(ctx)
	oid := goji.Param(r, "id")

	if oid == "" {
		log.Error("missing id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("transaction-list-cico.html")
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

	sb, err := url.PathUnescape(queryParams.Get("sort"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	sbv := cpb.SortOrder_DESC
	if sb == "asc" {
		sbv = cpb.SortOrder_ASC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	currentPage, _ := strconv.Atoi(queryParams.Get("page"))
	if currentPage == 0 {
		currentPage = 1
	}
	sbcv := cpb.SortByColumn_TransactionCompletedTime
	switch sbc {
	case "provider":
		sbcv = cpb.SortByColumn_Provider
	case "dateprocessed":
		sbcv = cpb.SortByColumn_TransactionCompletedTime
	case "cicoamount":
		sbcv = cpb.SortByColumn_TotalAmount
	case "fee":
		sbcv = cpb.SortByColumn_Fee
	case "commission":
		sbcv = cpb.SortByColumn_Commission
	}
	prfl, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	etd := s.getEnforceTemplateData(ctx)
	profile := prfl.GetProfile()
	ctx = md.Add("x-forward-dsaorgid", profile.OrgID).ToOutgoing(ctx)
	lr := &cpb.CICOTransactListRequest{
		Limit:           limitPerPage,
		Offset:          offset,
		SortOrder:       sbv,
		SortByColumn:    sbcv,
		OrgID:           oid,
		ReferenceNumber: searchTerms,
	}
	var res *cpb.CICOTransactListResponse
	if apiEnvType == "production" {
		res, err = s.drpLV.CICOTransactList(ctx, lr)
	} else {
		res, err = s.drpSB.CICOTransactList(ctx, lr)
	}

	if err != nil {
		logging.WithError(err, log).Error("listing transactions cico")
	}
	newReq := &spbl.GetPartnerListRequest{
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: psvc.ServiceType_CASHINCASHOUT.String(),
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

	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	var total int32
	tempData := transactionCICOTempData{
		CICOTransacts:            []*cpb.CICOTransact{},
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
	if res != nil && res.CICOTransacts != nil {
		tempData.CICOTransacts = res.GetCICOTransacts()
		total = res.Total
	}

	if len(tempData.CICOTransacts) > 0 {
		tempData.PaginationData = paginator.NewPaginator(int32(currentPage), limitPerPage, total, r)
	}

	tempData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
