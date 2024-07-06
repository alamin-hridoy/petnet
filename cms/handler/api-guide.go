package handler

import (
	"net/http"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
)

type (
	apiKeyGuideTempData struct {
		UserInfo                   *User
		Environment                string
		CompanyName                string
		OnboardingIncompleteStatus bool
		HasLiveAccess              bool
		ServiceRequest             bool
		PresetPermission           map[string]map[string]bool
	}
)

func (s *Server) getApiGuide(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	oid := mw.GetOrgID(r.Context())
	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	template := s.templates.Lookup("api-key-guide.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	usrInfo := s.GetUserInfoFromCookie(w, r, false)
	etd := s.getEnforceTemplateData(r.Context())
	data := apiKeyGuideTempData{
		UserInfo:                   &usrInfo.UserInfo,
		Environment:                sandEnv,
		HasLiveAccess:              s.hasLiveAccess(r.Context(), oid),
		CompanyName:                usrInfo.CompanyName,
		OnboardingIncompleteStatus: pf.GetProfile().Status == ppb.Status_Incomplete,
		PresetPermission:           etd.PresetPermission,
		ServiceRequest:             etd.ServiceRequests,
	}

	data.UserInfo.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
