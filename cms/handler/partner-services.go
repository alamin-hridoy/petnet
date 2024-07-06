package handler

import (
	"html/template"
	"net/http"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"

	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbupb "brank.as/rbac/gunk/v1/user"
)

type PartnerServicesTemplateData struct {
	UserInfo           *User
	PresetPermission   map[string]map[string]bool
	CSRFFieldValue     string
	ServiceRequest     bool
	CSRFField          template.HTML
	ServicePartnerList []ServicePartnerList
}

type ServicePartnerList struct {
	ID             string
	Stype          string
	Name           string
	Status         string
	Updated        time.Time
	ServiceName    string
	UpdatedBy      string
	UpdatedUsrName string
	DisableReason  string
}

// getPartnerServices is get action for partner services form to display
func (s *Server) getPartnerServices(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	template := s.templates.Lookup("partner-services.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := mw.GetOrgID(ctx)
	svc, err := s.pf.GetPartnerList(ctx, &spbl.GetPartnerListRequest{})
	if err != nil {
		log.Error("failed to get partner list")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	userLists := []string{}
	for _, ptnr := range svc.GetPartnerList() {
		if ptnr.UpdatedBy != "" {
			userLists = append(userLists, ptnr.UpdatedBy)
		}
	}

	userLists = uniqueSlice(userLists)
	usrs, err := s.rbac.ListUsers(ctx, &rbupb.ListUsersRequest{
		ID:    userLists,
		OrgID: oid,
	})
	if err != nil {
		log.Error("failed to get updated by user list")
	}
	userInfoLists := usrs.GetUser()

	var serviceTypeList []ServicePartnerList
	SPList := svc.GetPartnerList()
	if len(SPList) > 0 {
		for _, v := range SPList {
			updated := v.Updated.AsTime()
			serviceTypeList = append(serviceTypeList, ServicePartnerList{
				ID:             v.GetID(),
				Stype:          v.GetStype(),
				Name:           v.GetName(),
				Status:         v.GetStatus(),
				Updated:        updated,
				ServiceName:    v.GetServiceName(),
				UpdatedBy:      v.GetUpdatedBy(),
				UpdatedUsrName: getUserName(userInfoLists, v.GetUpdatedBy()),
				DisableReason:  v.GetDisableReason(),
			})
		}
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(r.Context())
	templateData := PartnerServicesTemplateData{
		CSRFField:          csrf.TemplateField(r),
		UserInfo:           &usrInfo.UserInfo,
		PresetPermission:   etd.PresetPermission,
		ServiceRequest:     etd.ServiceRequests,
		CSRFFieldValue:     csrf.Token(r),
		ServicePartnerList: serviceTypeList,
	}
	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func getUserName(usrInfo map[string]*rbupb.User, uid string) string {
	usrName := ""
	if usrInfo[uid] != nil {
		usrName = usrInfo[uid].FirstName + " " + usrInfo[uid].LastName
	}
	return usrName
}

// Update Partner Services Status
func (s *Server) updatePartnerServicesStatus(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}
	var form ServicePartnerList
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Stype, validation.Required),
		validation.Field(&form.Status, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Stype"] != nil {
				log.Error("Stype is Required")
			}
			if err["Status"] != nil {
				log.Error("Status is Required")
			}
		}
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if form.Status == "DISABLED" {
		_, err := s.pf.DisableMultiplePartnerList(ctx, &spbl.DisableMultiplePartnerListRequest{
			Stypes:        strings.Split(form.Stype, ","),
			DisableReason: form.DisableReason,
			UpdatedBy:     mw.GetUserID(ctx),
		})
		if err != nil {
			log.Error("failed to disable partner list")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	if form.Status == "ENABLED" {
		_, err := s.pf.EnableMultiplePartnerList(ctx, &spbl.EnableMultiplePartnerListRequest{
			Stypes:    strings.Split(form.Stype, ","),
			UpdatedBy: mw.GetUserID(ctx),
		})
		if err != nil {
			log.Error("failed to enable partner list")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}
	http.Redirect(w, r, partnerServicesPath, http.StatusSeeOther)
}
