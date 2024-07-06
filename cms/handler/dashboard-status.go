package handler

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"brank.as/petnet/cms/storage"
	"brank.as/petnet/gunk/drp/v1/dsa"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	epb "brank.as/petnet/gunk/dsa/v1/email"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbupb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/protobuf/types/known/emptypb"
)

type (
	DashboardStatusForm struct {
		CSRFField         template.HTML
		OrgID             string
		Status            string
		DsaCode           string
		TerminalIdOtc     string
		TerminalIdDigital string
	}
)

type iDRP interface {
	revcom.RevenueCommissionServiceClient
	dsa.DSAServiceClient
}

func (s *Server) resolveDRP() iDRP {
	if s.env == "production" {
		return s.drpLV
	}

	return s.drpSB
}

func (s *Server) postDashboardChangeStatus(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form DashboardStatusForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	tOtc := false
	tDigital := false
	ss, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{
		OrgID: form.OrgID,
	})
	if err != nil {
		logging.WithError(err, log).Error("org TransactionTypes not found")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if form.Status == "1" {
		if strings.Contains(ss.Profile.TransactionTypes, "OTC") {
			tOtc = true
		}
		if strings.Contains(ss.Profile.TransactionTypes, "DIGITAL") {
			tDigital = true
		}
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&form.Status, validation.Required, is.Digit),
		validation.Field(&form.DsaCode, validation.When(form.Status == "1", validation.Required, validation.Length(1, 3))),
		validation.Field(&form.TerminalIdOtc, validation.When(tOtc, validation.Required, validation.Length(0, 15))),
		validation.Field(&form.TerminalIdDigital, validation.When(tDigital, validation.Required, validation.Length(0, 15))),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	prf, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{
		OrgID: form.OrgID,
	})
	if err != nil {
		logging.WithError(err, log).Error("getting profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if form.DsaCode != "" {
		// TODO: Hoew to get these values?
		vatable, tin := "1", "1234567890"
		if prf.GetProfile().GetDsaCode() == "" {
			if prf.GetProfile().GetDsaCode() != form.DsaCode {
				_, err := s.resolveDRP().CreateDSA(ctx, &dsa.CreateDSARequest{
					DsaCode:      form.DsaCode,
					DsaName:      prf.GetProfile().GetBusinessInfo().CompanyName,
					EmailAddress: prf.GetProfile().GetBusinessInfo().CompanyEmail,
					Address:      prf.GetProfile().GetBusinessInfo().GetAddress().Address1,
					City:         prf.GetProfile().GetBusinessInfo().GetAddress().City,
					Province:     prf.GetProfile().GetBusinessInfo().GetAddress().State,
					UpdatedBy:    mw.GetUserID(ctx),
					Vatable:      vatable,
					Tin:          tin,
				})
				if err != nil {
					logging.WithError(err, log).Error("creating dsa")
					http.Redirect(w, r, errorPath, http.StatusSeeOther)
					return
				}
			}
		}

		if prf.GetProfile().GetDsaCode() != "" {
			ldsa, err := s.resolveDRP().ListDSA(ctx, &emptypb.Empty{})
			if err != nil {
				logging.WithError(err, log).Error("listing dsa")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			for _, v := range ldsa.GetDSAList() {
				if v.DsaCode == prf.GetProfile().GetDsaCode() {
					_, err := s.resolveDRP().UpdateDSA(ctx, &dsa.UpdateDSARequest{
						DsaCode:       form.DsaCode,
						DsaID:         v.DsaID,
						DsaName:       prf.GetProfile().GetBusinessInfo().CompanyName,
						EmailAddress:  prf.GetProfile().GetBusinessInfo().CompanyEmail,
						Address:       prf.GetProfile().GetBusinessInfo().GetAddress().Address1,
						City:          prf.GetProfile().GetBusinessInfo().GetAddress().City,
						Province:      prf.GetProfile().GetBusinessInfo().GetAddress().State,
						UpdatedBy:     mw.GetUserID(ctx),
						Zipcode:       prf.GetProfile().GetBusinessInfo().GetAddress().PostalCode,
						ContactPerson: prf.GetProfile().GetBusinessInfo().ContactPerson,
						// Not sure how to get these values as we doesn't store president and General Manager in Company Profile
						// But these fields are mared as required in DRP Update DSA API
						President:      prf.GetProfile().GetBusinessInfo().ContactPerson,
						GeneralManager: prf.GetProfile().GetBusinessInfo().ContactPerson,
						Vatable:        vatable,
						Tin:            tin,
					})
					if err != nil {
						logging.WithError(err, log).Error("updating dsa")
						http.Redirect(w, r, errorPath, http.StatusSeeOther)
						return
					}
				}
			}
		}
	}

	status, err := strconv.Atoi(form.Status)
	if err != nil {
		logging.WithError(err, log).Error("convert string to int")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pf := &ppb.UpsertProfileRequest{
		Profile: &ppb.OrgProfile{
			OrgID:             form.OrgID,
			Status:            ppb.Status(status),
			DsaCode:           form.DsaCode,
			TerminalIdOtc:     form.TerminalIdOtc,
			TerminalIdDigital: form.TerminalIdDigital,
			UserID:            ss.Profile.UserID,
		},
	}
	d, err := json.Marshal(pf)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	me := &mfaEvent{
		resource: string(storage.Status),
		action:   tpb.ActionType_Update,
		data:     d,
	}
	if s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if _, err := s.pf.UpsertProfile(ctx, pf); err != nil {
			logging.WithError(err, log).Error("updating status")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		ui, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: ss.Profile.UserID})
		if err != nil {
			logging.WithError(err, log).Info("getting user")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if ui == nil || ui.User == nil {
			logging.WithError(err, log).Error("getting user")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		isEmailSend := "1"
		if _, err := s.pf.SendOnboardingReminder(ctx, &epb.SendOnboardingReminderRequest{
			Email:  ui.User.Email,
			OrgID:  ss.Profile.OrgID,
			UserID: ss.Profile.UserID,
		}); err != nil {
			isEmailSend = "0"
			logging.WithError(err, log).Error("sending reminder")
		}
		http.Redirect(w, r, fmt.Sprintf("/dashboard/dsa-applicant-list/%s?emailsend=%s", form.OrgID, isEmailSend), http.StatusSeeOther)
	}
	http.Redirect(w, r, fmt.Sprintf("/dashboard/dsa-applicant-list/%s?show_otp=true", form.OrgID), http.StatusSeeOther)
}

func (s *Server) getProfileByDsaCode(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	isDsaCodeExist := false

	var form DashboardStatusForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.DsaCode, validation.Required, validation.Length(1, 3)),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["DsaCode"] != nil {
				log.Error("DsaCode is Required")
			}
		}
		return
	}

	pf, err := s.pf.GetProfileByDsaCode(ctx, &ppb.GetProfileByDsaCodeRequest{DsaCode: form.DsaCode})
	if err != nil {
		errMsg := strings.ToUpper(err.Error())
		if strings.Contains(errMsg, "ORG NOT FOUND") {
			isDsaCodeExist = false
		} else {
			logging.WithError(err, log).Info("getting profile")
			return
		}
	}

	if pf != nil {
		isDsaCodeExist = true
	}

	jsn, _ := json.Marshal(isDsaCodeExist)
	w.WriteHeader(200)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsn)
}
