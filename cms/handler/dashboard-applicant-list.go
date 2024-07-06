package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"brank.as/petnet/cms/paginator"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	rbupb "brank.as/rbac/gunk/v1/user"
	"github.com/gorilla/csrf"
)

type (
	DSAApplicant struct {
		CSRFField            template.HTML
		ID                   string
		OrgID                string
		CompanyName          string
		DateApplied          time.Time
		RiskScore            string
		Status               string
		ReminderSent         bool
		IsDocumentsSubmitted bool
		User                 User
	}

	dsaTemplateData struct {
		CSRFFieldValue   string
		DSAApplicants    []DSAApplicant
		PaginationData   paginator.Paginator
		UserInfo         *User
		SearchTerm       string
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
	}
)

const perPage int32 = 10

func (s *Server) getDSAApplicantList(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	template := s.templates.Lookup("dashboard-applicant-list.html")
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

	newReq := &ppb.ListProfilesRequest{
		Limit:       perPage,
		Offset:      offset,
		CompanyName: searchTerm,
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

	rScores, err := url.PathUnescape(queryParams.Get("risk_score"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if rScores != "" {
		rsStrArr := strings.Split(rScores, ",")
		rsArr := []ppb.RiskScore{}
		for _, rs := range rsStrArr {
			switch rs {
			case "low":
				rsArr = append(rsArr, ppb.RiskScore_Low)
			case "medium":
				rsArr = append(rsArr, ppb.RiskScore_Medium)
			case "high":
				rsArr = append(rsArr, ppb.RiskScore_High)
			}
		}

		newReq.RiskScore = rsArr
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

	subDoc, err := url.PathUnescape(queryParams.Get("submitted_document"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	if subDoc != "" {
		if subDoc == "submitted" {
			newReq.SubmittedDocument = ppb.SubmittedDocument_Submitted
		} else if subDoc == "not-submitted" {
			newReq.SubmittedDocument = ppb.SubmittedDocument_NotSubmitted
		}
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	newReq.OrgType = ppb.OrgType_DSA
	pf, err := s.pf.ListProfiles(r.Context(), newReq)
	etd := s.getEnforceTemplateData(ctx)
	tData := dsaTemplateData{
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		CSRFFieldValue:   csrf.Token(r),
		DSAApplicants:    []DSAApplicant{},
		UserInfo:         &usrInfo.UserInfo,
		SearchTerm:       searchTerm,
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
	var dsaAplicantList []DSAApplicant
	for _, applicant := range pf.Profiles {
		u, err := s.rbac.GetUser(r.Context(), &rbupb.GetUserRequest{ID: applicant.UserID})
		if err != nil {
			logging.WithError(err, log).Info("getting user")
		} else {
			status := applicant.GetStatus().String()
			if applicant.GetStatus() == ppb.Status_UnknownStatus {
				status = "Incomplete"
			}
			convertedTime := applicant.GetDateApplied().AsTime()
			if !applicant.GetDateApplied().IsValid() {
				convertedTime = applicant.Created.AsTime()
			}
			// Check all documentation submitted
			docSubmitted := true

			fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{
				OrgID: applicant.GetOrgID(),
			})
			if err != nil {
				logging.WithError(err, log).Info("listing files")
				return
			}
			sub := make(map[string]bool)
			for _, f := range fs.FileUploads {
				sub[f.Type.String()] = f.GetSubmitted().String() == "True"
			}
			for _, b := range []bool{
				sub[fpb.UploadType_IDPhoto.String()],
				sub[fpb.UploadType_Picture.String()],
				sub[fpb.UploadType_NBIClearance.String()],
				sub[fpb.UploadType_CourtClearance.String()],
				sub[fpb.UploadType_IncorporationPapers.String()],
				sub[fpb.UploadType_MayorsPermit.String()],
				sub[fpb.UploadType_FinancialStatement.String()],
				sub[fpb.UploadType_BankStatement.String()],
				sub[fpb.UploadType_Questionnaire.String()],
			} {
				switch b {
				case false:
					docSubmitted = false
				}
			}

			dsaAplicantList = append(dsaAplicantList, DSAApplicant{
				CSRFField:            csrf.TemplateField(r),
				ID:                   applicant.ID,
				OrgID:                applicant.OrgID,
				CompanyName:          applicant.BusinessInfo.CompanyName,
				DateApplied:          convertedTime,
				RiskScore:            applicant.RiskScore.String(),
				Status:               status,
				ReminderSent:         applicant.GetReminderSent() == ppb.Boolean_True,
				IsDocumentsSubmitted: docSubmitted,
				User:                 userData(u.GetUser()),
			})
		}
	}

	templateData := dsaTemplateData{
		CSRFFieldValue:   csrf.Token(r),
		DSAApplicants:    dsaAplicantList,
		UserInfo:         &usrInfo.UserInfo,
		SearchTerm:       searchTerm,
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
	}
	if len(dsaAplicantList) > 0 {
		templateData.PaginationData = paginator.NewPaginator(int32(currentPage), perPage, pf.Total, r)
	}

	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, templateData); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
