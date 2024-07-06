package handler

import (
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"brank.as/petnet/cms/storage"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	ppb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type passTempData struct {
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	LoginUserInfo    *User
	CSRFField        template.HTML
	MessageDetails   string
	Errors           map[string]string
	IsPetnetAdmin    bool
	CompanyName      string
}

type PasswordEditForm struct {
	CurrentPassword string
	NewPassword     string
	ConfirmPassword string
}

func (s *Server) geteditPass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(r.Context())
	template := s.templates.Lookup("edit-password.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	etd := s.getEnforceTemplateData(ctx)
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	tempData := passTempData{
		IsPetnetAdmin:    mw.IsPetnetOwner(ctx),
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		LoginUserInfo:    &usrInfo.UserInfo,
		CSRFField:        csrf.TemplateField(r),
		CompanyName:      usrInfo.CompanyName,
	}
	queryParams := r.URL.Query()
	msgType, _ := url.PathUnescape(queryParams.Get("msgType"))
	if msgType == "success" {
		tempData.MessageDetails = "Your password has been changed successfully."
	}
	tempData.LoginUserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, tempData); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
	}
}

func (s *Server) posteditPass(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("User ID not found")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var form PasswordEditForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := map[string]string{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.CurrentPassword, validation.Required, validation.Length(8, 64)),
		validation.Field(&form.NewPassword, validation.Required, validation.Length(8, 64)),
		validation.Field(&form.ConfirmPassword, validation.Required, validation.Length(8, 64)),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {

			if err["CurrentPassword"] != nil {
				formErr["CurrentPassword"] = err["CurrentPassword"].Error()
			}
			if err["NewPassword"] != nil {
				formErr["NewPassword"] = err["NewPassword"].Error()
			}
			if err["ConfirmPassword"] != nil {
				formErr["ConfirmPassword"] = err["ConfirmPassword"].Error()
			}
			if formErr["NewPassword"] == "" && formErr["ConfirmPassword"] == "" && !isPasswordValid(form.NewPassword) {
				formErr["NewPassword"] = "The password entered does not meet policy."
			}
		}
	}
	if form.NewPassword != form.ConfirmPassword {
		formErr["ConfirmPassword"] = "New Password and Confirm Password doesn't match."
	}
	if len(formErr) > 0 {
		template := s.templates.Lookup("edit-password.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		etd := s.getEnforceTemplateData(ctx)
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		tempData := passTempData{
			PresetPermission: etd.PresetPermission,
			ServiceRequest:   etd.ServiceRequests,
			LoginUserInfo:    &usrInfo.UserInfo,
			CSRFField:        csrf.TemplateField(r),
			Errors:           formErr,
		}
		if err := template.Execute(w, tempData); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
		}
	}
	resCP, err := s.rbacUserAuth.ChangePassword(ctx, &ppb.ChangePasswordRequest{
		UserID:      uid,
		OldPassword: form.CurrentPassword,
		NewPassword: form.NewPassword,
	})
	if err != nil {
		if status.Code(err) == codes.PermissionDenied && strings.Contains(err.Error(), "old password is invalid") {
			formErr["CurrentPassword"] = "Invalid credentials. Please try again."
			template := s.templates.Lookup("edit-password.html")
			if template == nil {
				log.Error("unable to load template")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			etd := s.getEnforceTemplateData(ctx)
			usrInfo := s.GetUserInfoFromCookie(w, r, false)

			tempData := passTempData{
				PresetPermission: etd.PresetPermission,
				LoginUserInfo:    &usrInfo.UserInfo,
				CSRFField:        csrf.TemplateField(r),
				Errors:           formErr,
			}
			if err := template.Execute(w, tempData); err != nil {
				log.Infof("error with template execution: %+v", err)
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
			}
		}
		log.Error("unable to connect api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	cr := &ppb.ConfirmUpdateRequest{
		UserID:     uid,
		MFAEventID: resCP.GetMFAEventID(),
		MFAType:    resCP.GetMFAType(),
	}
	d, err := json.Marshal(cr)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	me := &mfaEvent{
		resource: string(storage.ChangePassword),
		action:   tpb.ActionType_Update,
		data:     d,
		eventID:  resCP.MFAEventID,
	}
	if resCP.MFAEventID != "" {
		if err := s.initMFAEvent(w, r, me); err != nil {
			if err != storage.MFANotFound {
				logging.WithError(err, log).Error("initializing mfa event")
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
			http.Redirect(w, r, edPassPath+"?msgType=success", http.StatusSeeOther)
		}
		http.Redirect(w, r, edPassPath+"?show_otp=true", http.StatusSeeOther)
	}
	http.Redirect(w, r, edPassPath+"?msgType=success", http.StatusSeeOther)
}
