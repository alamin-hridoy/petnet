package handler

import (
	"html/template"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	bpa "brank.as/petnet/gunk/drp/v1/remittance"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	spf "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfsvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"brank.as/petnet/svcutil/random"
	rbupb "brank.as/rbac/gunk/v1/user"
)

type ProvidersTemplateData struct {
	UserInfo         *User
	PresetPermission map[string]map[string]bool
	CSRFFieldValue   string
	ServiceRequest   bool
	CSRFField        template.HTML
	ProviderName     string
	Stype            string
	Provider         string
	Errors           map[string]error
	ProviderList     []CreateProviderForm
	AllProviders     []string
	MessageDetails   string
	MsgType          string
	AllProviderList  []AllPartnerList
	SearchTerm       string
}

type CreateProviderForm struct {
	ProviderName string
	Stype        string
	Provider     string
}

type AllPartnerList struct {
	ID            string
	Stype         string
	Name          string
	Created       time.Time
	Updated       time.Time
	Status        string
	ServiceName   string
	UpdatedBy     string
	DisableReason string
	Platform      string
	IsProvider    bool
}

// getPartnerServices is get action for partner services form to display
func (s *Server) doGetproviders(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(r.Context())
	template := s.templates.Lookup("providers.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	sList := []string{}
	for _, v := range pfsvc.ServiceType_name {
		if v != "EMPTYSERVICETYPE" && (v == pfsvc.ServiceType_REMITTANCE.String() || v == pfsvc.ServiceType_CASHINCASHOUT.String()) {
			sList = append(sList, v)
		}
	}

	queryParams := r.URL.Query()
	searchTerm, err := url.PathUnescape(queryParams.Get("search-term"))
	if err != nil {
		log.Error("unable to decode url type param")
	}

	pl, err := s.pf.GetPartnerList(r.Context(), &spf.GetPartnerListRequest{
		Name:       searchTerm,
		IsProvider: true,
	})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	userLists := []string{}
	for _, ptnr := range pl.GetPartnerList() {
		if ptnr.UpdatedBy != "" {
			userLists = append(userLists, ptnr.UpdatedBy)
		}
	}
	userLists = uniqueSlice(userLists)
	userInfoLists := map[string]*rbupb.User{}
	usrs, err := s.rbac.ListUsers(ctx, &rbupb.ListUsersRequest{
		ID:    userLists,
		OrgID: mw.GetOrgID(ctx),
	})
	if err == nil {
		userInfoLists = usrs.GetUser()
	}
	ptnrs := []AllPartnerList{}
	for _, v := range pl.GetPartnerList() {
		ptnrs = append(ptnrs, AllPartnerList{
			ID:            v.GetID(),
			Stype:         v.GetStype(),
			Name:          v.GetName(),
			Created:       v.GetCreated().AsTime(),
			Updated:       v.GetUpdated().AsTime(),
			Status:        v.GetStatus(),
			ServiceName:   v.GetServiceName(),
			UpdatedBy:     getUserName(userInfoLists, v.UpdatedBy),
			DisableReason: v.GetDisableReason(),
			Platform:      v.GetPlatform(),
			IsProvider:    v.GetIsProvider(),
		})
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(r.Context())
	templateData := ProvidersTemplateData{
		CSRFField:        csrf.TemplateField(r),
		UserInfo:         &usrInfo.UserInfo,
		PresetPermission: etd.PresetPermission,
		ServiceRequest:   etd.ServiceRequests,
		CSRFFieldValue:   csrf.Token(r),
		AllProviders:     sList,
		AllProviderList:  ptnrs,
		SearchTerm:       searchTerm,
	}
	msgType, _ := url.PathUnescape(queryParams.Get("msgType"))

	switch msgType {
	case "success":
		templateData.MessageDetails = "Provider Added successfully."
		templateData.MsgType = "success"
	case "error":
		templateData.MessageDetails = "Something Went Wrong."
		templateData.MsgType = "error"
	case "update":
		templateData.MessageDetails = "Provider Updated successfully."
		templateData.MsgType = "success"
	case "delete":
		templateData.MessageDetails = "Provider Deleted successfully."
		templateData.MsgType = "success"
	}
	templateData.UserInfo.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, templateData); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) createProvider(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	ServiceInfo := map[string]string{
		"CASHINCASHOUT": "cico",
		"REMITTANCE":    "remittance",
	}

	var f CreateProviderForm

	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&f,
		validation.Field(&f.ProviderName, validation.Required),
		validation.Field(&f.Provider, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["ProviderName"] != nil {
				formErr["ProviderName"] = err["ProviderName"]
			}
			if err["Provider"] != nil {
				formErr["Provider"] = err["Provider"]
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}
	var ff ProvidersTemplateData
	ff.Errors = formErr
	ff.CSRFField = csrf.TemplateField(r)
	if len(formErr) > 0 {
		template := s.templates.Lookup("providers.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, providersPath+"?msgType=error", http.StatusSeeOther)
			return
		}
	}
	str := s.generateUniqueCode(r, f.Provider, 0)
	if str == "" {
		log.Error("unable to generate unique code")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	apiReq := &bpa.PartnersCreateRequest{
		PartnerCode: str,
		PartnerName: f.ProviderName,
		Service:     ServiceInfo[f.Provider],
	}
	apiRes, err := s.drpSB.PartnersCreate(ctx, apiReq)
	if err != nil {
		logging.WithError(err, log).Error("create partner list api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if apiRes == nil || apiRes.GetResult() == nil {
		logging.WithError(err, log).Error("get partner list api failed")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if _, err := s.pf.CreatePartnerList(ctx, &spf.CreatePartnerListRequest{
		PartnerList: &spf.PartnerList{
			Stype:            str,
			Name:             f.ProviderName,
			Created:          &timestamppb.Timestamp{},
			Updated:          &timestamppb.Timestamp{},
			Status:           spb.PartnerStatusType_ENABLED.String(),
			ServiceName:      strings.ToUpper(f.Provider),
			UpdatedBy:        mw.GetUserID(ctx),
			IsProvider:       true,
			PerahubPartnerID: strconv.Itoa(int(apiRes.GetResult().ID)),
		},
	}); err != nil {
		logging.WithError(err, log).Error("creating partner list")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, providersPath+"?msgType=success", http.StatusSeeOther)
}

func (s *Server) generateUniqueCode(r *http.Request, partner string, n int) string {
	ctx := r.Context()
	str := partner[:3]
	if n != 0 {
		str = random.InvitationCode(3)
	}
	_, err := s.pf.GetPartnerByStype(ctx, &spf.GetPartnerByStypeRequest{
		Stype: str,
	})
	if status.Code(err) != codes.NotFound {
		str = s.generateUniqueCode(r, partner, n+1)
	}
	return strings.ToUpper(str)
}

func (s Server) postProviderUpdate(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	lid := goji.Param(r, "lid")
	if lid == "" {
		log.Error("missing provider id id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var f ProvidersTemplateData
	if err := s.decoder.Decode(&f, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	ServiceInfo := map[string]string{
		"CASHINCASHOUT": "cico",
		"REMITTANCE":    "remittance",
	}

	if _, err := s.pf.UpdatePartnerList(ctx, &spf.UpdatePartnerListRequest{
		PartnerList: &spf.PartnerList{
			Stype:       lid,
			Name:        f.ProviderName,
			Status:      spb.PartnerStatusType_ENABLED.String(),
			ServiceName: strings.ToUpper(f.Provider),
			IsProvider:  true,
			UpdatedBy:   mw.GetUserID(ctx),
		},
	}); err != nil {
		logging.WithError(err, log).Error("update partner list")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	getPtnr, err := s.pf.GetPartnerByStype(ctx, &spf.GetPartnerByStypeRequest{
		Stype: lid,
	})
	if err != nil {
		logging.WithError(err, log).Error("get partner by stype failed")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if getPtnr == nil || getPtnr.GetPartnerList() == nil {
		logging.WithError(err, log).Error("get partner by stype failed")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	apiReq := &bpa.PartnersUpdateRequest{
		ID:          getPtnr.PartnerList.GetPerahubPartnerID(),
		PartnerCode: lid,
		PartnerName: f.ProviderName,
		Service:     ServiceInfo[getPtnr.PartnerList.ServiceName], // need to get service
	}

	if _, err := s.drpSB.PartnersUpdate(ctx, apiReq); err != nil {
		logging.WithError(err, log).Error("update partner list api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, providersPath+"?msgType=update", http.StatusSeeOther)
}

func (s Server) deleteProvider(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	lid := goji.Param(r, "lid")
	if lid == "" {
		log.Error("missing provider id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	getPtnr, err := s.pf.GetPartnerByStype(ctx, &spf.GetPartnerByStypeRequest{
		Stype: lid,
	})
	if err != nil {
		logging.WithError(err, log).Error("get partner by stype failed")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if getPtnr == nil || getPtnr.GetPartnerList() == nil {
		logging.WithError(err, log).Error("get partner by stype failed")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if _, err := s.pf.DeletePartnerList(ctx, &spf.DeletePartnerListRequest{
		Stype: lid,
	}); err != nil {
		logging.WithError(err, log).Error("updating provider")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	ss := &bpa.PartnersDeleteRequest{
		ID: getPtnr.PartnerList.PerahubPartnerID,
	}

	if _, err := s.drpSB.PartnersDelete(ctx, ss); err != nil {
		logging.WithError(err, log).Error("delete partner list api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, providersPath+"?msgType=delete", http.StatusSeeOther)
}
