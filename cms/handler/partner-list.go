package handler

import (
	"context"
	"html/template"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"

	ppf "brank.as/petnet/gunk/dsa/v1/user"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfsvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/gorilla/csrf"
)

type (
	PartnerListApplicant struct {
		ID        string
		Stype     string
		Name      string
		Created   time.Time
		Updated   time.Time
		Status    string
		ImageLink string
	}

	partnerListTemplateData struct {
		PartnerListApplicants         []PartnerListApplicant
		DisabledPartnerListApplicants []PartnerListApplicant
		UserInfo                      *User
		SearchTerm                    string
		PresetPermission              map[string]map[string]bool
		ServiceRequest                bool
		CSRFField                     template.HTML
	}
)

func (s *Server) postPartnerLists(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)
	var form PartnerListApplicant
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	lr, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{
		Stype:       form.Stype,
		ServiceName: pfsvc.ServiceType_REMITTANCE.String(),
	})
	if err != nil {
		if _, err := s.createPartner(ctx, form); err != nil {
			logging.WithError(err, log).Error("creating partner")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, partnerListPath, http.StatusSeeOther)
			return
		}
	}
	if len(lr.GetPartnerList()) > 0 {
		if _, err := s.updatePartner(ctx, form); err != nil {
			logging.WithError(err, log).Error("updating partner")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, partnerListPath, http.StatusSeeOther)
			return
		}
	} else {
		if _, err := s.createPartner(ctx, form); err != nil {
			logging.WithError(err, log).Error("creating partner")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		} else {
			http.Redirect(w, r, partnerListPath, http.StatusSeeOther)
			return
		}
	}
}

func (s *Server) createPartner(ctx context.Context, form PartnerListApplicant) (*spbl.CreatePartnerListResponse, error) {
	return s.pf.CreatePartnerList(ctx, &spbl.CreatePartnerListRequest{
		PartnerList: &spbl.PartnerList{
			ID:     form.ID,
			Stype:  form.Stype,
			Name:   form.Name,
			Status: form.Status,
		},
	})
}

func (s *Server) updatePartner(ctx context.Context, form PartnerListApplicant) (*spbl.UpdatePartnerListResponse, error) {
	return s.pf.UpdatePartnerList(ctx, &spbl.UpdatePartnerListRequest{
		PartnerList: &spbl.PartnerList{
			ID:     form.ID,
			Stype:  form.Stype,
			Name:   form.Name,
			Status: form.Status,
		},
	})
}

func (s *Server) getPartnerLists(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	template := s.templates.Lookup("partner-list.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	queryParams := r.URL.Query()
	searchTerm, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		logging.WithError(err, log).Error("unable to decode url type param")
	}

	newReq := &spbl.GetPartnerListRequest{
		Name:        searchTerm,
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: pfsvc.ServiceType_REMITTANCE.String(),
	}
	svc, err := s.pf.GetPartnerList(r.Context(), newReq)
	var partnerListApplicantList []PartnerListApplicant
	if err != nil {
		log.Error("failed to Get Partner List")
	} else {
		for _, sv := range svc.GetPartnerList() {
			imgLink := path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean("/images/partners/"+sv.GetStype()+".png"), "/")))
			partnerListApplicantList = append(partnerListApplicantList, PartnerListApplicant{
				ID:        sv.GetID(),
				Stype:     sv.GetStype(),
				Name:      sv.GetName(),
				Created:   sv.GetCreated().AsTime(),
				Updated:   sv.GetUpdated().AsTime(),
				Status:    sv.GetStatus(),
				ImageLink: imgLink,
			})
		}
	}
	newReqDis := &spbl.GetPartnerListRequest{
		Name:        searchTerm,
		Status:      spb.PartnerStatusType_DISABLED.String(),
		ServiceName: pfsvc.ServiceType_REMITTANCE.String(),
	}
	svcDis, err := s.pf.GetPartnerList(r.Context(), newReqDis)
	var DisabledPartnerListApplicants []PartnerListApplicant
	if err != nil {
		log.Error("failed to Get Partner List")
	} else {
		for _, sv := range svcDis.GetPartnerList() {
			imgLink := path.Join("/", s.assetFS.HashName(strings.TrimPrefix(path.Clean("/images/partners/"+sv.GetStype()+".png"), "/")))
			DisabledPartnerListApplicants = append(DisabledPartnerListApplicants, PartnerListApplicant{
				ID:        sv.GetID(),
				Stype:     sv.GetStype(),
				Name:      sv.GetName(),
				Created:   sv.GetCreated().AsTime(),
				Updated:   sv.GetUpdated().AsTime(),
				Status:    sv.GetStatus(),
				ImageLink: imgLink,
			})
		}
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(r.Context())
	templateData := partnerListTemplateData{
		PartnerListApplicants:         partnerListApplicantList,
		DisabledPartnerListApplicants: DisabledPartnerListApplicants,
		UserInfo:                      &usrInfo.UserInfo,
		SearchTerm:                    searchTerm,
		PresetPermission:              etd.PresetPermission,
		ServiceRequest:                etd.ServiceRequests,
		CSRFField:                     csrf.TemplateField(r),
	}
	uid := mw.GetUserID(r.Context())
	gp, err := s.pf.GetUserProfile(r.Context(), &ppf.GetUserProfileRequest{
		UserID: uid,
	})
	if err != nil {
		log.Error("failed to get profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	templateData.UserInfo.ProfileImage = gp.GetProfile().ProfilePicture
	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
