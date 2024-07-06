package handler

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"strings"

	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"google.golang.org/protobuf/types/known/timestamppb"

	bpb "brank.as/petnet/gunk/dsa/v2/branch"
	fpb "brank.as/petnet/gunk/dsa/v2/file"
	ppb "brank.as/petnet/gunk/dsa/v2/profile"
)

type LocationBranchData struct {
	ID                         string
	Title                      string
	Address                    Address
	PhoneNumber                string
	FaxNumber                  string
	ContactPerson              string
	CSRFField                  template.HTML
	BusinessInfo               *ppb.BusinessInfo
	Branches                   []*bpb.Branch
	OrgID                      string
	Status                     string
	AllDocsSubmitted           bool
	Errors                     map[string]error
	PresetPermission           map[string]map[string]bool
	ServiceRequest             bool
	SearchTerm                 string
	User                       *User
	DsaCode                    string
	TerminalIdOtc              string
	TerminalIdDigital          string
	TransactionTypesForDSACode TransactionTypesForDSACode
}

func (a Address) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address1, validation.Required),
		validation.Field(&a.City, validation.Required),
		validation.Field(&a.PostalCode, validation.Required),
		validation.Field(&a.State, validation.Required),
	)
}

func (s *Server) getLocationConfig(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())

	template := s.templates.Lookup("location-configuration.html")
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

	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}

	fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: oid})
	if err != nil {
		logging.WithError(err, log).Info("getting files")
	}

	bs, err := s.pf.ListBranches(r.Context(), &bpb.ListBranchesRequest{
		OrgID: oid,
		Title: searchTerm,
	})
	if err != nil {
		logging.WithError(err, log).Info("getting profile")
	}
	transactionTypesForDSACode := checkTrnTypeForDSACode(pf.Profile.GetTransactionTypes())
	businessInfo := pf.GetProfile().GetBusinessInfo()
	etd := s.getEnforceTemplateData(r.Context())
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	data := LocationBranchData{
		BusinessInfo:               businessInfo,
		Branches:                   bs.Branches,
		PresetPermission:           etd.PresetPermission,
		ServiceRequest:             etd.ServiceRequests,
		SearchTerm:                 searchTerm,
		User:                       &usrInfo.UserInfo,
		DsaCode:                    pf.Profile.DsaCode,
		TerminalIdOtc:              pf.Profile.TerminalIdOtc,
		TerminalIdDigital:          pf.Profile.TerminalIdDigital,
		TransactionTypesForDSACode: transactionTypesForDSACode,
	}
	data.CSRFField = csrf.TemplateField(r)
	data.OrgID = oid
	status := pf.GetProfile().GetStatus().String()
	if pf.GetProfile().GetStatus() == ppb.Status_UnknownStatus {
		status = "Incomplete"
	}
	if pf.GetProfile().GetStatus() == ppb.Status_PendingDocuments {
		status = "Pending Documents"
	}

	data.Status = status
	data.AllDocsSubmitted = allDocsSubmitted(fs.GetFileUploads())
	data.User.ProfileImage = usrInfo.ProfileImage
	if err := template.Execute(w, data); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s Server) postLocationUpdate(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	uid := mw.GetUserID(ctx)
	oid := mw.GetOrgID(ctx)
	var form BusinessInfoForm
	ts := strings.TrimSpace
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if _, err := s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
		Profile: &ppb.OrgProfile{
			UserID: uid,
			OrgID:  oid,
			BusinessInfo: &ppb.BusinessInfo{
				Address: &ppb.Address{
					Address1:   ts(form.Address1),
					City:       ts(form.City),
					State:      ts(form.State),
					PostalCode: ts(form.PostalCode),
				},
			},
		},
	}); err != nil {
		logging.WithError(err, log).Error("updating businessinfo profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/dashboard/location/"+oid, http.StatusSeeOther)
	return
}

func (s Server) postLocationConfig(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form LocationBranchData
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := validation.Errors{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Title, validation.Required),
		validation.Field(&form.PhoneNumber, is.Digit, validation.Required),
		validation.Field(&form.FaxNumber, is.Digit),
		validation.Field(&form.Address, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["PhoneNumber"] != nil {
				formErr["PhoneNumber"] = errors.New("Phone number needs to be digits only")
			}
			if err["FaxNumber"] != nil {
				formErr["FaxNumber"] = errors.New("Fax number should be digit")
			}
		}
		logging.WithError(err, log).Error("invalid request")
	}

	form.Errors = formErr
	form.CSRFField = csrf.TemplateField(r)
	if len(formErr) > 0 {
		template := s.templates.Lookup("location-configuration.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		pf, err := s.pf.GetProfile(r.Context(), &ppb.GetProfileRequest{OrgID: oid})
		if err != nil {
			logging.WithError(err, log).Info("getting profile")
		}

		form.BusinessInfo = pf.GetProfile().GetBusinessInfo()
		status := pf.GetProfile().GetStatus().String()
		if pf.GetProfile().GetStatus() == ppb.Status_UnknownStatus {
			status = "Incomplete"
		}
		if pf.GetProfile().GetStatus() == ppb.Status_PendingDocuments {
			status = "Pending Documents"
		}

		fs, err := s.pf.ListFiles(r.Context(), &fpb.ListFilesRequest{OrgID: oid})
		if err != nil {
			logging.WithError(err, log).Info("getting files")
		}
		form.Status = status
		form.OrgID = oid
		form.AllDocsSubmitted = allDocsSubmitted(fs.GetFileUploads())
		etd := s.getEnforceTemplateData(ctx)
		form.PresetPermission = etd.PresetPermission
		form.ServiceRequest = etd.ServiceRequests
		if err := template.Execute(w, form); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		return
	}

	if _, err := s.pf.UpsertBranch(ctx, &bpb.UpsertBranchRequest{
		Branch: &bpb.Branch{
			ID:    form.ID,
			OrgID: oid,
			Title: form.Title,
			Address: &ppb.Address{
				Address1:   form.Address.Address1,
				City:       form.Address.City,
				State:      form.Address.State,
				PostalCode: form.Address.PostalCode,
			},
			PhoneNumber:   form.PhoneNumber,
			FaxNumber:     form.FaxNumber,
			ContactPerson: form.ContactPerson,
		},
	}); err != nil {
		logging.WithError(err, log).Error("creating location branch")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("/dashboard/location/%s", oid), http.StatusSeeOther)
}

func (s Server) updateLocationConfig(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	lid := goji.Param(r, "lid")
	if lid == "" {
		log.Error("missing localtion branch id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var form LocationBranchData
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Title, validation.Required),
		validation.Field(&form.PhoneNumber, validation.Required),
		validation.Field(&form.Address, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if _, err := s.pf.UpsertBranch(ctx, &bpb.UpsertBranchRequest{
		Branch: &bpb.Branch{
			ID:    lid,
			OrgID: oid,
			Title: form.Title,
			Address: &ppb.Address{
				Address1:   form.Address.Address1,
				City:       form.Address.City,
				State:      form.Address.State,
				PostalCode: form.Address.PostalCode,
			},
			PhoneNumber:   form.PhoneNumber,
			FaxNumber:     form.FaxNumber,
			ContactPerson: form.ContactPerson,
		},
	}); err != nil {
		logging.WithError(err, log).Error("updating location branch")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/location/%s", oid), http.StatusSeeOther)
}

func (s Server) deleteLocationConfig(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	oid := goji.Param(r, "id")
	if oid == "" {
		log.Error("missing org id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	lid := goji.Param(r, "lid")
	if lid == "" {
		log.Error("missing localtion branch id in url param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if _, err := s.pf.UpsertBranch(ctx, &bpb.UpsertBranchRequest{
		Branch: &bpb.Branch{
			ID:      lid,
			OrgID:   oid,
			Deleted: timestamppb.Now(),
		},
	}); err != nil {
		logging.WithError(err, log).Error("updating location branch")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard/location/%s", oid), http.StatusSeeOther)
}
