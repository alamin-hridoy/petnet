package handler

import (
	"html/template"
	"net/http"
	"net/url"

	ppf "brank.as/petnet/gunk/dsa/v1/user"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	ppb "brank.as/rbac/gunk/v1/user"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type TempData struct {
	PresetPermission map[string]map[string]bool
	ServiceRequest   bool
	LoginUserInfo    *User
	CSRFField        template.HTML
	Errors           map[string]string
	MessageDetails   string
	IsPetnetAdmin    bool
	CompanyName      string
}

type ProfileEditForm struct {
	FirstName    string
	LastName     string
	ProfileImage string
}

var validedProfileImgFileType = []string{
	"image/jpeg", "image/png", "image/jpg",
}

func (s *Server) getProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(r.Context())
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("User ID not found")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("profile.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	etd := s.getEnforceTemplateData(ctx)
	tempData := TempData{
		IsPetnetAdmin:    mw.IsPetnetOwner(ctx),
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		LoginUserInfo:    &usrInfo.UserInfo,
		CSRFField:        csrf.TemplateField(r),
		CompanyName:      usrInfo.CompanyName,
	}

	tempData.LoginUserInfo.ProfileImage = usrInfo.ProfileImage
	queryParams := r.URL.Query()
	msgType, _ := url.PathUnescape(queryParams.Get("msgType"))
	if msgType == "success" {
		tempData.MessageDetails = "Successfully updated your profile."
	}
	if err := template.Execute(w, tempData); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
	}
}

func (s *Server) postProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	uid := mw.GetUserID(ctx)
	if uid == "" {
		log.Error("User ID not found")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		logging.WithError(err, log).Error("parsing form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var form ProfileEditForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	formErr := map[string]string{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.FirstName, validation.Required),
		validation.Field(&form.LastName, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["FirstName"] != nil {
				formErr["FirstName"] = err["FirstName"].Error()
			}
			if err["LastName"] != nil {
				formErr["LastName"] = err["LastName"].Error()
			}
		}
	}
	ferr := s.validateSingleFileType(r, "ProfileImage", validedProfileImgFileType)
	if ferr != nil {
		formErr["ProfileImage"] = ferr.Error()
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	if len(formErr) > 0 {
		template := s.templates.Lookup("profile.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		etd := s.getEnforceTemplateData(ctx)
		tempData := TempData{
			IsPetnetAdmin:    mw.IsPetnetOwner(ctx),
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

	_, err := s.rbacUserAuth.UpdateUser(ctx, &ppb.UpdateUserRequest{
		UserID:    uid,
		FirstName: form.FirstName,
		LastName:  form.LastName,
	})
	if err != nil {
		log.Error("unable to connect api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	_, _, fuerr := r.FormFile("ProfileImage")
	fid, _, uerr := s.storeSingleToGCS(r, "ProfileImage", uid)
	if uerr != nil && fuerr == nil {
		logging.WithError(uerr, log).Error("storing ProfileImage")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	gp, err := s.pf.GetUserProfile(ctx, &ppf.GetUserProfileRequest{
		UserID: uid,
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	gpProfile := gp.GetProfile()
	uupr := &ppf.UpdateUserProfileRequest{
		Profile: &ppf.Profile{
			ID:     gpProfile.GetID(),
			OrgID:  gpProfile.GetOrgID(),
			UserID: uid,
		},
	}

	if fid != "" {
		uupr.Profile.ProfilePicture = fid
	}

	if fid == "" {
		uupr.Profile.ProfilePicture = usrInfo.ProfileImage
	}
	_, err = s.pf.UpdateUserProfile(ctx, uupr)
	if err != nil {
		log.Error("failed to update profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	s.GetUserInfoFromCookie(w, r, true)
	http.Redirect(w, r, profilePath+"?msgType=success", http.StatusSeeOther)
}
