package handler

import (
	"net/http"
	"net/url"
	"strconv"
	"time"

	"brank.as/petnet/cms/paginator"
	mipb "brank.as/petnet/gunk/drp/v1/microinsurance"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
)

type transactionMITempData struct {
	MiTransacts      []*mipb.InsuranceTransaction
	UserInfo         *User
	CompanyName      string
	PaginationData   paginator.Paginator
	SearchTerms      string
	OrgId            string
	Environment      string
	HasLiveAccess    bool
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
}

func (s *Server) getMITransactionList(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	md := metautils.ExtractIncoming(ctx)
	oid := goji.Param(r, "id")

	if oid == "" {
		log.Error("missing id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("transaction-list-mi.html")
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

	sbv := mipb.SortOrder_DESC
	if sb == "asc" {
		sbv = mipb.SortOrder_ASC
	}

	sbc, err := url.PathUnescape(queryParams.Get("sort_column"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	currentPage, _ := strconv.Atoi(queryParams.Get("page"))
	if currentPage == 0 {
		currentPage = 1
	}
	sbcv := mipb.SortByColumn_TransactionCompletedTime
	if sbc == "miamount" {
		sbcv = mipb.SortByColumn_TotalAmount
	} else if sbc == "fee" {
		sbcv = mipb.SortByColumn_Fee
	} else if sbc == "commission" {
		sbcv = mipb.SortByColumn_Commission
	} else if sbc == "dateprocessed" {
		sbcv = mipb.SortByColumn_TransactionCompletedTime
	}
	prfl, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	etd := s.getEnforceTemplateData(ctx)
	profile := prfl.GetProfile()
	ctx = md.Add("x-forward-dsaorgid", profile.OrgID).ToOutgoing(ctx)
	lr := &mipb.GetTransactionListRequest{
		DateFrom:     "2006-01-02",
		DateTo:       time.Now().Format("2006-01-02"),
		Limit:        limitPerPage,
		Offset:       offset,
		SortOrder:    sbv,
		SortByColumn: sbcv,
		OrgID:        oid,
	}
	var res *mipb.TransactionListResult
	if apiEnvType == "production" {
		res, err = s.drpLV.GetTransactionList(ctx, lr)
	} else {
		res, err = s.drpSB.GetTransactionList(ctx, lr)
	}

	if err != nil {
		logging.WithError(err, log).Error("listing transactions microinsurance")
	}

	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	var total int32
	tempData := transactionMITempData{
		MiTransacts:      []*mipb.InsuranceTransaction{},
		UserInfo:         &usrInfo.UserInfo,
		SearchTerms:      searchTerms,
		OrgId:            oid,
		Environment:      apiEnvType,
		HasLiveAccess:    s.hasLiveAccess(r.Context(), oid),
		CompanyName:      profile.GetBusinessInfo().GetCompanyName(),
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
	}
	if res != nil && res.Transactions != nil {
		tempData.MiTransacts = res.GetTransactions()
		total = res.Total
	}

	if len(tempData.MiTransacts) > 0 {
		tempData.PaginationData = paginator.NewPaginator(int32(currentPage), limitPerPage, total, r)
	}

	tempData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
