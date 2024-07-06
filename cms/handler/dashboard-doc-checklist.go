package handler

import (
	"fmt"
	"html/template"
	"net/http"

	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	BusinessInform struct {
		CompanyName                  string
		IsIDPhotoSubmitted           bool
		IsPictureSubmitted           bool
		IsNBISubmitted               bool
		IsCourtClearanceSubmitted    bool
		IsIncorpPapersSubmitted      bool
		IsMayorsPermitSubmitted      bool
		IDPhotoDateCheck             string
		PictureDateCheck             string
		NBIClearanceDateCheck        string
		CourtClearanceDateCheck      string
		IncorporationPapersDateCheck string
		MayorsPermitDateCheck        string
	}

	FinancialInform struct {
		IsFinancialSubmitted        bool
		IsBankSubmitted             bool
		FinancialStatementDateCheck string
		BankStatementDateCheck      string
	}

	DRPInfoDocForm struct {
		QuestionnaireSubmitted bool
		QuestionnaireDateCheck string
	}

	Submittion struct {
		CSRFField                  template.HTML
		OrgID                      string
		Status                     string
		AllDocsSubmitted           bool
		BusinessInfo               BusinessInform
		FinancialInfo              FinancialInform
		DRPInfo                    DRPInfoDocForm
		PresetPermission           map[string]map[string]bool
		User                       *User
		ServiceRequest             bool
		DsaCode                    string
		TerminalIdOtc              string
		TerminalIdDigital          string
		TransactionTypesForDSACode TransactionTypesForDSACode
	}

	DocSubmittionForm struct {
		DocFieldName      string
		DocTimeFieldName  string
		DocFieldIsChecked string
	}
)

func (s *Server) getDashboardCheckoutList(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("dashboard-doc-checklist.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("listing files")
	}

	sub := make(map[string]bool)
	chk := make(map[string]string)
	for _, f := range fs.FileUploads {
		sub[f.Type.String()] = f.GetSubmitted().String() == "True"
		chk[f.Type.String()] = f.GetDateChecked().AsTime().Format("02 January, 2006")
	}

	profile := pf.GetProfile()
	status := profile.GetStatus().String()
	if profile.GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	if profile.GetStatus() == ppb.Status_PendingDocuments {
		status = "Pending Documents"
	}
	transactionTypesForDSACode := checkTrnTypeForDSACode(profile.GetTransactionTypes())
	etd := s.getEnforceTemplateData(r.Context())
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := Submittion{
		OrgID:            oid,
		Status:           status,
		AllDocsSubmitted: allDocsSubmitted(fs.GetFileUploads()),
		CSRFField:        csrf.TemplateField(r),
		BusinessInfo: BusinessInform{
			CompanyName:                  profile.BusinessInfo.CompanyName,
			IsIDPhotoSubmitted:           sub[fpb.UploadType_IDPhoto.String()],
			IsPictureSubmitted:           sub[fpb.UploadType_Picture.String()],
			IsNBISubmitted:               sub[fpb.UploadType_NBIClearance.String()],
			IsCourtClearanceSubmitted:    sub[fpb.UploadType_CourtClearance.String()],
			IsIncorpPapersSubmitted:      sub[fpb.UploadType_IncorporationPapers.String()],
			IsMayorsPermitSubmitted:      sub[fpb.UploadType_MayorsPermit.String()],
			IDPhotoDateCheck:             chk[fpb.UploadType_IDPhoto.String()],
			PictureDateCheck:             chk[fpb.UploadType_Picture.String()],
			NBIClearanceDateCheck:        chk[fpb.UploadType_NBIClearance.String()],
			CourtClearanceDateCheck:      chk[fpb.UploadType_CourtClearance.String()],
			IncorporationPapersDateCheck: chk[fpb.UploadType_IncorporationPapers.String()],
			MayorsPermitDateCheck:        chk[fpb.UploadType_MayorsPermit.String()],
		},
		FinancialInfo: FinancialInform{
			IsFinancialSubmitted:        sub[fpb.UploadType_FinancialStatement.String()],
			IsBankSubmitted:             sub[fpb.UploadType_BankStatement.String()],
			FinancialStatementDateCheck: chk[fpb.UploadType_FinancialStatement.String()],
			BankStatementDateCheck:      chk[fpb.UploadType_BankStatement.String()],
		},
		DRPInfo: DRPInfoDocForm{
			QuestionnaireSubmitted: sub[fpb.UploadType_Questionnaire.String()],
			QuestionnaireDateCheck: chk[fpb.UploadType_Questionnaire.String()],
		},
		PresetPermission:           etd.PresetPermission,
		ServiceRequest:             etd.ServiceRequests,
		User:                       &usrInfo.UserInfo,
		DsaCode:                    profile.DsaCode,
		TerminalIdOtc:              profile.TerminalIdOtc,
		TerminalIdDigital:          profile.TerminalIdDigital,
		TransactionTypesForDSACode: transactionTypesForDSACode,
	}
	data.User.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postDashboardDocument(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	var form DocSubmittionForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid := goji.Param(r, "id")
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	profileData := getDynamicStructField(form.DocFieldName, form.DocFieldIsChecked)
	profileData.OrgID = oid
	profileData.UserID = pf.GetProfile().GetUserID()

	if _, err := s.pf.UpsertFiles(ctx, &fpb.UpsertFilesRequest{
		FileUploads: []*fpb.FileUpload{
			profileData,
		},
	}); err != nil {
		logging.WithError(err, log).Error("Add Fee and Commissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/checkout-list/%s", oid), http.StatusSeeOther)
}

func getDynamicStructField(fieldType string, IsChecked string) *fpb.FileUpload {
	ic := ppb.Boolean_True
	if IsChecked == "unchecked" {
		ic = ppb.Boolean_False
	}
	switch fieldType {
	case "IDPhoto":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_IDPhoto,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Picture":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_Picture,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Nbi":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_NBIClearance,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Court":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_CourtClearance,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Incorp":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_IncorporationPapers,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Mayor":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_MayorsPermit,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Finance":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_FinancialStatement,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Bank":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_BankStatement,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	case "Que":
		return &fpb.FileUpload{
			Type:        fpb.UploadType_Questionnaire,
			Submitted:   ic,
			DateChecked: timestamppb.Now(),
		}
	default:
		return &fpb.FileUpload{}
	}
}

func allDocsSubmitted(fs []*fpb.FileUpload) bool {
	sub := make(map[string]bool)
	for _, f := range fs {
		sub[f.Type.String()] = f.GetSubmitted().String() == "True"
	}
	if sub[fpb.UploadType_IDPhoto.String()] &&
		sub[fpb.UploadType_Picture.String()] &&
		sub[fpb.UploadType_NBIClearance.String()] &&
		sub[fpb.UploadType_CourtClearance.String()] &&
		sub[fpb.UploadType_IncorporationPapers.String()] &&
		sub[fpb.UploadType_MayorsPermit.String()] &&
		sub[fpb.UploadType_FinancialStatement.String()] &&
		sub[fpb.UploadType_BankStatement.String()] &&
		sub[fpb.UploadType_Questionnaire.String()] {
		return true
	}
	return false
}
