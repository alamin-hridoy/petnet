package handler

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"time"

	ppu "brank.as/petnet/gunk/dsa/v1/user"
	fpb "brank.as/petnet/gunk/dsa/v2/fees"
	fipb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type (
	feeCommissionTemplateData struct {
		CSRFField        template.HTML
		OrgID            string
		Status           string
		AllDocsSubmitted bool
		Fees             []*fpb.Fee
		Commissions      []*fpb.Fee
		BusinessInfo     *ppb.BusinessInfo
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
		User             *User
	}

	FeeCommissionForm struct {
		Rate      *fpb.Rate
		Type      string
		StartDate string
		Updated   string
		EndDate   string
	}
)

func (s *Server) getFeeManagement(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	template := s.templates.Lookup("fee-management.html")
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
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := feeCommissionTemplateData{
		BusinessInfo: loadProfile.GetBusinessInfo(),
		User:         &usrInfo.UserInfo,
	}

	data.CSRFField = csrf.TemplateField(r)
	data.OrgID = oid

	status := loadProfile.GetStatus().String()
	if loadProfile.GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	if loadProfile.GetStatus() == ppb.Status_PendingDocuments {
		status = "Pending Documents"
	}

	fs, err := s.pf.ListFiles(r.Context(), &fipb.ListFilesRequest{OrgID: string(oid)})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	data.Status = status
	data.AllDocsSubmitted = allDocsSubmitted(fs.GetFileUploads())

	rsf, err := s.pf.ListFees(r.Context(), &fpb.ListFeesRequest{
		OrgID: string(oid),
		Type:  fpb.FeeType_TypeFee.String(),
	})
	if err != nil {
		logging.WithError(err, log).Info("getting fees")
	}
	data.Fees = rsf.GetFees()

	rsc, err := s.pf.ListFees(r.Context(), &fpb.ListFeesRequest{
		OrgID: string(oid),
		Type:  fpb.FeeType_TypeCommission.String(),
	})
	if err != nil {
		logging.WithError(err, log).Info("getting fees")
	}
	data.Commissions = rsc.GetFees()
	etd := s.getEnforceTemplateData(r.Context())
	data.PresetPermission = etd.PresetPermission
	data.ServiceRequest = etd.ServiceRequests
	uidd := mw.GetUserID(r.Context())
	gp, err := s.pf.GetUserProfile(r.Context(), &ppu.GetUserProfileRequest{
		UserID: uidd,
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	data.User.ProfileImage = gp.GetProfile().ProfilePicture
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postFeeCommission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	var form FeeCommissionForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.StartDate, validation.Required),
		validation.Field(&form.EndDate, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	layout := "2006-01-02"
	t1, err := time.Parse(layout, form.StartDate)
	if err != nil {
		logging.WithError(err, log).Error("decoding form")
		return
	}
	t2, err := time.Parse(layout, form.EndDate)
	if err != nil {
		logging.WithError(err, log).Error("decoding form")
		return
	}
	t2 = t2.Add(time.Hour*time.Duration(23) +
		time.Minute*time.Duration(59) +
		time.Second*time.Duration(59))

	convertedStartDate := timestamppb.New(t1)
	convertedEndDate := timestamppb.New(t2)
	oid := goji.Param(r, "id")
	tc := convertedStartDate.AsTime().Before(convertedEndDate.AsTime())
	if !tc {
		logging.WithError(err, log).Error("End date must be greater than start date")
		http.Redirect(w, r, "/dashboard/fee-management/"+oid, http.StatusSeeOther)
		return
	}

	feeType := fpb.FeeType_TypeFee
	if form.Type == "2" {
		feeType = fpb.FeeType_TypeCommission
	}

	res, err := s.pf.UpsertFee(ctx, &fpb.UpsertFeeRequest{
		Fee: &fpb.Fee{
			OrgID: oid,
			Type:  feeType,
			Schedule: &fpb.Schedule{
				StartDate: convertedStartDate,
				EndDate:   convertedEndDate,
			},
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("Add Fee and Commissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	fields := r.Form.Get("Feetierlen")

	totField, err := strconv.Atoi(fields)
	if err != nil {
		logging.WithError(err, log).Error("Add Fee and Commissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	for key := 1; key <= totField; key++ {
		ds := fmt.Sprintf("DSAShare%d", key)
		mint := fmt.Sprintf("MinTran%d", key)
		maxt := fmt.Sprintf("MaxTran%d", key)

		dsaShare := r.Form.Get(ds)
		minValue := r.Form.Get(mint)
		maxValue := r.Form.Get(maxt)
		if _, err := s.pf.UpsertRate(ctx, &fpb.UpsertRateRequest{
			Rate: &fpb.Rate{
				FeeComID:  res.ID,
				MinVolume: minValue,
				MaxVolume: maxValue,
				TxnRate:   dsaShare,
			},
			FeeCommissionID: res.ID,
		}); err != nil {
			logging.WithError(err, log).Error("Add Rate under fees")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/fee-management/%s", oid), http.StatusSeeOther)
}

func (s *Server) updateFeeCommission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	var updateForm FeeCommissionForm
	if err := s.decoder.Decode(&updateForm, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&updateForm,
		validation.Field(&updateForm.StartDate, validation.Required),
		validation.Field(&updateForm.EndDate, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	layout := "2006-01-02"
	t1, err := time.Parse(layout, updateForm.StartDate)
	if err != nil {
		logging.WithError(err, log).Error("decoding form")
		return
	}

	t2, err := time.Parse(layout, updateForm.EndDate)
	if err != nil {
		logging.WithError(err, log).Error("decoding form")
		return
	}
	t2 = t2.Add(time.Hour*time.Duration(23) +
		time.Minute*time.Duration(59) +
		time.Second*time.Duration(59))

	feeType := fpb.FeeType_TypeFee
	if updateForm.Type == "2" {
		feeType = fpb.FeeType_TypeCommission
	}

	convertedStartDate := timestamppb.New(t1)
	convertedEndDate := timestamppb.New(t2)
	convertedUpdatedDate := timestamppb.New(time.Now())
	oid := goji.Param(r, "id")
	tc := convertedStartDate.AsTime().Before(convertedEndDate.AsTime())
	if !tc {
		logging.WithError(err, log).Error("End date must be greater than start date")
		http.Redirect(w, r, "/dashboard/fee-management/"+oid, http.StatusSeeOther)
		return
	}

	fid := goji.Param(r, "fid")
	if _, err := s.pf.UpsertFee(ctx, &fpb.UpsertFeeRequest{
		Fee: &fpb.Fee{
			ID:    fid,
			OrgID: oid,
			Type:  feeType,
			Schedule: &fpb.Schedule{
				StartDate: convertedStartDate,
				EndDate:   convertedEndDate,
			},
			Updated: convertedUpdatedDate,
		},
	}); err != nil {
		logging.WithError(err, log).Error("Add Fee and Commissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	fields := r.Form.Get("Feetierlen")

	totField, err := strconv.Atoi(fields)
	if err != nil {
		logging.WithError(err, log).Error("Add Fee and Commissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	for key := 0; key < totField; key++ {
		ds := fmt.Sprintf("DSAShare%d", key)
		mint := fmt.Sprintf("MinTran%d", key)
		maxt := fmt.Sprintf("MaxTran%d", key)
		rateid := fmt.Sprintf("rateid%d", key)
		dsaShare := r.Form.Get(ds)
		minValue := r.Form.Get(mint)
		maxValue := r.Form.Get(maxt)
		rateidValue := r.Form.Get(rateid)
		if _, err := s.pf.UpsertRate(ctx, &fpb.UpsertRateRequest{
			Rate: &fpb.Rate{
				FeeComID:  fid,
				ID:        rateidValue,
				MinVolume: minValue,
				MaxVolume: maxValue,
				TxnRate:   dsaShare,
			},
			FeeCommissionID: fid,
		}); err != nil {
			logging.WithError(err, log).Error("Add Rate under fees")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/fee-management/%s", oid), http.StatusSeeOther)
}

func (s *Server) deleteFeeCommission(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	convertedDeleteDate := timestamppb.New(time.Now())
	oid := goji.Param(r, "id")
	fid := goji.Param(r, "fid")
	if _, err := s.pf.UpsertFee(ctx, &fpb.UpsertFeeRequest{
		Fee: &fpb.Fee{
			ID:      fid,
			OrgID:   oid,
			Deleted: convertedDeleteDate,
		},
	}); err != nil {
		logging.WithError(err, log).Error("Delete Fee and Commissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/fee-management/%s", oid), http.StatusSeeOther)
}
