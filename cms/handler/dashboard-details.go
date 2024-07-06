package handler

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"

	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"

	ppf "brank.as/petnet/gunk/dsa/v1/user"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	rbupb "brank.as/rbac/gunk/v1/user"
)

type (
	DashboardDetails struct {
		CSRFField                  template.HTML
		CSRFFieldValue             string
		ID                         string
		User                       *User
		Userr                      User
		OrgID                      string
		Status                     string
		AllDocsSubmitted           bool
		BusinessInfo               *BusinessInfo
		AccountInfo                *AccountInfo
		DRPInfo                    *DRPInfo
		TransactionTypes           []string
		Files                      Files
		BaseURL                    string
		ServiceRequest             bool
		PresetPermission           map[string]map[string]bool
		DsaCode                    string
		TerminalIdOtc              string
		TerminalIdDigital          string
		TransactionTypesForDSACode TransactionTypesForDSACode
	}

	TransactionTypesForDSACode struct {
		Otc     bool
		Digital bool
	}

	Files struct {
		IDPhoto                 []string
		Picture                 []string
		NBIClearance            []string
		CourtClearance          []string
		IncorporationPapers     []string
		MayorsPermit            []string
		FinancialStatement      []string
		BankStatement           []string
		Questionnaire           []string
		NDA                     []string
		IDPhotoName             map[string]string
		PictureName             map[string]string
		NBIClearanceName        map[string]string
		CourtClearanceName      map[string]string
		IncorporationPapersName map[string]string
		MayorsPermitName        map[string]string
		FinancialStatementName  map[string]string
		BankStatementName       map[string]string
		QuestionnaireName       map[string]string
		NDAName                 map[string]string
	}
	ReqFiles struct {
		CICO_nda        []string
		CICO_sis        []string
		CICO_psf        []string
		CICO_pspp       []string
		CICO_sec        []string
		CICO_gis        []string
		CICO_afs        []string
		CICO_bir        []string
		CICO_bsp        []string
		CICO_aml        []string
		CICO_sccas      []string
		CICO_vgid       []string
		CICO_cp         []string
		CICO_moa        []string
		CICO_amla       []string
		CICO_mtpp       []string
		CICO_is         []string
		CICO_edd        []string
		CICO_ndaName    map[string]string
		CICO_sisName    map[string]string
		CICO_psfName    map[string]string
		CICO_psppName   map[string]string
		CICO_secName    map[string]string
		CICO_gisName    map[string]string
		CICO_afsName    map[string]string
		CICO_birName    map[string]string
		CICO_bspName    map[string]string
		CICO_amlName    map[string]string
		CICO_sccasName  map[string]string
		CICO_vgidName   map[string]string
		CICO_cpName     map[string]string
		CICO_moaName    map[string]string
		CICO_amlaName   map[string]string
		CICO_mtppName   map[string]string
		CICO_isName     map[string]string
		CICO_eddName    map[string]string
		Mbf             []string
		Sec             []string
		Gis             []string
		Afs             []string
		Brs             []string
		Bmp             []string
		Scbr            []string
		Via             []string
		Moa             []string
		WU_cd           []string
		WU_lbp          []string
		WU_sr           []string
		WU_lgis         []string
		WU_dtisspa      []string
		WU_birf         []string
		WU_bspr         []string
		WU_iqa          []string
		WU_bqa          []string
		AYA_dfedr       []string
		AYA_ddf         []string
		AYA_ialaws      []string
		AYA_aialaws     []string
		AYA_mtl         []string
		AYA_cbpr        []string
		AYA_brdcsa      []string
		AYA_gis         []string
		AYA_ccpwi       []string
		AYA_fas         []string
		AYA_am          []string
		AYA_laf         []string
		AYA_birr        []string
		AYA_od          []string
		Bspr            []string
		Cp              []string
		Aml             []string
		Nnda            []string
		Psf             []string
		Psp             []string
		Kddq            []string
		Sis             []string
		MbfName         map[string]string
		SecName         map[string]string
		GisName         map[string]string
		AfsName         map[string]string
		BrsName         map[string]string
		BmpName         map[string]string
		ScbrName        map[string]string
		ViaName         map[string]string
		MoaName         map[string]string
		WU_cdName       map[string]string
		WU_lbpName      map[string]string
		WU_srName       map[string]string
		WU_lgisName     map[string]string
		WU_dtisspaName  map[string]string
		WU_birfName     map[string]string
		WU_bsprName     map[string]string
		WU_iqaName      map[string]string
		WU_bqaName      map[string]string
		AYA_dfedrName   map[string]string
		AYA_ddfName     map[string]string
		AYA_ialawsName  map[string]string
		AYA_aialawsName map[string]string
		AYA_mtlName     map[string]string
		AYA_cbprName    map[string]string
		AYA_brdcsaName  map[string]string
		AYA_gisName     map[string]string
		AYA_ccpwiName   map[string]string
		AYA_fasName     map[string]string
		AYA_amName      map[string]string
		AYA_lafName     map[string]string
		AYA_birrName    map[string]string
		AYA_odName      map[string]string
		BsprName        map[string]string
		CpName          map[string]string
		AmlName         map[string]string
		NndaName        map[string]string
		PsfName         map[string]string
		PspName         map[string]string
		KddqName        map[string]string
		SisName         map[string]string
		MI_mbf          []string
		MI_nda          []string
		MI_sec          []string
		MI_gis          []string
		MI_afs          []string
		MI_bir          []string
		MI_scb          []string
		MI_via          []string
		MI_moa          []string
		MI_mbfName      map[string]string
		MI_ndaName      map[string]string
		MI_secName      map[string]string
		MI_gisName      map[string]string
		MI_afsName      map[string]string
		MI_birName      map[string]string
		MI_scbName      map[string]string
		MI_viaName      map[string]string
		MI_moaName      map[string]string
		BP_nda          []string
		BP_sis          []string
		BP_psf          []string
		BP_psp          []string
		BP_sec          []string
		BP_gis          []string
		BP_lpafs        []string
		BP_bir          []string
		BP_bsp          []string
		BP_aml          []string
		BP_sccas        []string
		BP_vg           []string
		BP_cp           []string
		BP_moa          []string
		BP_amla         []string
		BP_mttp         []string
		BP_aib          []string
		BP_bp           []string
		BP_sci          []string
		BP_edd          []string
		BP_ndaName      map[string]string
		BP_sisName      map[string]string
		BP_psfName      map[string]string
		BP_pspName      map[string]string
		BP_secName      map[string]string
		BP_gisName      map[string]string
		BP_lpafsName    map[string]string
		BP_birName      map[string]string
		BP_bspName      map[string]string
		BP_amlName      map[string]string
		BP_sccasName    map[string]string
		BP_vgName       map[string]string
		BP_cpName       map[string]string
		BP_moaName      map[string]string
		BP_amlaName     map[string]string
		BP_mttpName     map[string]string
		BP_aibName      map[string]string
		BP_bpName       map[string]string
		BP_sciName      map[string]string
		BP_eddName      map[string]string
	}

	BusinessInfo struct {
		CompanyName   string
		StoreName     string
		PhoneNumber   string
		FaxNumber     string
		Website       string
		CompanyEmail  string
		ContactPerson string
		Position      string
		Address       *Address
	}

	Address struct {
		Address1   string
		City       string
		State      string
		PostalCode string
	}

	AccountInfo struct {
		Bank              string
		BankAccountNumber string
		BankAccountHolder string
	}

	DRPInfo struct {
		WUInfo   *spb.WesternUnionPartner
		TFInfo   *spb.TransfastPartner
		IRInfo   *spb.IRemitPartner
		RIAInfo  *spb.RiaPartner
		MBInfo   *spb.MetroBankPartner
		RMInfo   *spb.RemitlyPartner
		BPIInfo  *spb.BPIPartner
		USSCInfo *spb.USSCPartner
		JPRInfo  *spb.JapanRemitPartner
		ICInfo   *spb.InstantCashPartner
		UNTInfo  *spb.UnitellerPartner
	}

	User struct {
		ID           string
		OrgID        string
		OrgName      string
		FirstName    string
		LastName     string
		ProfileImage string
		Email        string
		Aemail       string
		InviteStatus string
		CountryCode  string
		Phone        string
		Created      time.Time
		Updated      time.Time
	}
)

func (s *Server) getDashboardDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("dashboard-details.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
		return
	}
	fs, err := s.pf.ListFiles(ctx, &fpb.ListFilesRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("listing files")
		return
	}

	u, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: pf.Profile.UserID})
	if err != nil {
		logging.WithError(err, log).Info("getting user")
		return
	}

	svc, err := s.pf.GetPartners(ctx, &spb.GetPartnersRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting services")
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(ctx)
	details := dashboardData(pf.GetProfile(), fs.GetFileUploads(), svc.GetPartners())
	details.Userr = userData(u.GetUser())
	details.User = &usrInfo.UserInfo
	details.OrgID = oid
	details.CSRFFieldValue = csrf.Token(r)
	details.BaseURL = s.urls.Base + "/u/files/"
	details.CSRFField = csrf.TemplateField(r)
	details.PresetPermission = etd.PresetPermission
	details.ServiceRequest = etd.ServiceRequests

	uid := mw.GetUserID(ctx)
	gp, err := s.pf.GetUserProfile(ctx, &ppf.GetUserProfileRequest{
		UserID: uid,
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	details.User.ProfileImage = usrInfo.ProfileImage
	details.User.Aemail = gp.GetProfile().Email
	if err := template.Execute(w, details); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func userData(u *rbupb.User) User {
	return User{
		ID:           u.ID,
		OrgID:        u.OrgID,
		OrgName:      u.OrgName,
		FirstName:    u.FirstName,
		LastName:     u.LastName,
		Email:        u.Email,
		InviteStatus: u.InviteStatus,
		CountryCode:  u.CountryCode,
		Phone:        u.Phone,
		Created:      u.GetCreated().AsTime(),
		Updated:      u.GetUpdated().AsTime(),
	}
}

func dashboardData(bip *ppb.OrgProfile, fs []*fpb.FileUpload, svc *spb.Partners) DashboardDetails {
	status := bip.GetStatus().String()
	if bip.GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	transactionTypes := strings.Split(bip.GetTransactionTypes(), ",")
	transactionTypesForDSACode := checkTrnTypeForDSACode(bip.GetTransactionTypes())

	fm := DashboardDetails{
		OrgID:            bip.OrgID,
		Status:           status,
		AllDocsSubmitted: allDocsSubmitted(fs),
		TransactionTypes: transactionTypes,
		BusinessInfo: &BusinessInfo{
			CompanyName:   bip.GetBusinessInfo().GetCompanyName(),
			StoreName:     bip.GetBusinessInfo().GetStoreName(),
			PhoneNumber:   bip.GetBusinessInfo().GetPhoneNumber(),
			FaxNumber:     bip.GetBusinessInfo().GetFaxNumber(),
			Website:       bip.GetBusinessInfo().GetWebsite(),
			CompanyEmail:  bip.GetBusinessInfo().GetCompanyEmail(),
			ContactPerson: bip.GetBusinessInfo().GetContactPerson(),
			Position:      bip.GetBusinessInfo().GetPosition(),
			Address: &Address{
				Address1:   bip.GetBusinessInfo().GetAddress().GetAddress1(),
				City:       bip.GetBusinessInfo().GetAddress().GetCity(),
				State:      bip.GetBusinessInfo().GetAddress().GetState(),
				PostalCode: bip.GetBusinessInfo().GetAddress().GetPostalCode(),
			},
		},
		AccountInfo: &AccountInfo{
			Bank:              bip.GetAccountInfo().GetBank(),
			BankAccountNumber: bip.GetAccountInfo().GetBankAccountNumber(),
			BankAccountHolder: bip.GetAccountInfo().GetBankAccountHolder(),
		},
		DRPInfo: &DRPInfo{
			WUInfo:   svc.GetWesternUnionPartner(),
			TFInfo:   svc.GetTransfastPartner(),
			IRInfo:   svc.GetIRemitPartner(),
			RIAInfo:  svc.GetRiaPartner(),
			MBInfo:   svc.GetMetroBankPartner(),
			RMInfo:   svc.GetRemitlyPartner(),
			BPIInfo:  svc.GetBPIPartner(),
			USSCInfo: svc.GetUSSCPartner(),
			JPRInfo:  svc.GetJapanRemitPartner(),
			ICInfo:   svc.GetInstantCashPartner(),
			UNTInfo:  svc.GetUnitellerPartner(),
		},
		DsaCode:                    bip.DsaCode,
		TerminalIdOtc:              bip.TerminalIdOtc,
		TerminalIdDigital:          bip.TerminalIdDigital,
		TransactionTypesForDSACode: transactionTypesForDSACode,
	}
	for _, f := range fs {
		switch f.Type {
		case fpb.UploadType_IDPhoto:
			fm.Files.IDPhoto = f.GetFileNames()
			fm.Files.IDPhotoName = f.GetFileName()
		case fpb.UploadType_Picture:
			fm.Files.Picture = f.GetFileNames()
			fm.Files.PictureName = f.GetFileName()
		case fpb.UploadType_NBIClearance:
			fm.Files.NBIClearance = f.GetFileNames()
			fm.Files.NBIClearanceName = f.GetFileName()
		case fpb.UploadType_CourtClearance:
			fm.Files.CourtClearance = f.GetFileNames()
			fm.Files.CourtClearanceName = f.GetFileName()
		case fpb.UploadType_IncorporationPapers:
			fm.Files.IncorporationPapers = f.GetFileNames()
			fm.Files.IncorporationPapersName = f.GetFileName()
		case fpb.UploadType_MayorsPermit:
			fm.Files.MayorsPermit = f.GetFileNames()
			fm.Files.MayorsPermitName = f.GetFileName()
		case fpb.UploadType_FinancialStatement:
			fm.Files.FinancialStatement = f.GetFileNames()
			fm.Files.FinancialStatementName = f.GetFileName()
		case fpb.UploadType_BankStatement:
			fm.Files.BankStatement = f.GetFileNames()
			fm.Files.BankStatementName = f.GetFileName()
		case fpb.UploadType_Questionnaire:
			fm.Files.Questionnaire = f.GetFileNames()
			fm.Files.QuestionnaireName = f.GetFileName()
		case fpb.UploadType_NDA:
			fm.Files.NDA = f.GetFileNames()
			fm.Files.NDAName = f.GetFileName()
		}
	}
	return fm
}

func checkTrnTypeForDSACode(str string) TransactionTypesForDSACode {
	otc, digital := false, false
	if strings.Contains(str, "OTC") {
		otc = true
	}
	if strings.Contains(str, "DIGITAL") {
		digital = true
	}
	transactionTypesForDSACode := TransactionTypesForDSACode{
		Otc:     otc,
		Digital: digital,
	}

	return transactionTypesForDSACode
}
