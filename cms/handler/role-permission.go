package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"brank.as/petnet/cms/storage"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	pm "brank.as/petnet/svcutil/permission"
	ppb "brank.as/rbac/gunk/v1/permissions"
	"github.com/gorilla/csrf"
	"github.com/kenshaw/goji"
	"github.com/sirupsen/logrus"
)

type rolePermissionForm struct {
	RoleID                    string
	CSRFField                 template.HTML
	DSAListAndDetails         []string
	TransactionListAndDetails []string
	DocumentChecklist         []string
	FeeAndCommission          []string
	RiskAssessment            []string
	Currency                  []string
	LocationAndBranches       []string
	ServiceCatalog            []string
	RolePresetPermission      map[string]map[string]string
	PresetPermission          map[string]map[string]bool
	ServiceRequest            bool
	ListPermission            *ppb.ListPermissionResponse
	ListRoles                 *ppb.ListRoleResponse
	LoginUserInfo             *User
	ErrorMsg                  string
}

type PermissionData struct {
	ServiceName string
	Name        string
	Description string
	Resource    string
}

func (s *Server) getRolePermission(w http.ResponseWriter, r *http.Request) {
	s.permissionForm(w, r, false)
}

func (s *Server) permissionForm(w http.ResponseWriter, r *http.Request, hasError bool) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	rid := goji.Param(r, "id")
	oid := mw.GetOrgID(ctx)
	if rid == "" {
		log.Error("missing role id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("edit-role.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	pl, err := s.rbac.ListPermission(ctx, &ppb.ListPermissionRequest{
		OrgID: mw.GetOrgID(ctx),
	})
	if err != nil {
		log.Error("list permissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	res, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{ID: []string{rid}, OrgID: oid})
	if err != nil {
		log.Error("list Role")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	roles := res.GetRoles()
	presetPermission := make(map[string]map[string]string)
	for _, pp := range pl.GetPermissions() {
		if pp.Name != "" {
			presetPermission[pp.Name] = make(map[string]string)
		}
	}

	for _, pp := range pl.GetPermissions() {
		if pp.Action != "" {
			presetPermission[pp.Name][pp.Action] = ""
		}
	}

	for _, v := range roles {
		for _, r := range v.Permissions {
			for _, pp := range pl.GetPermissions() {
				if pp.Action != "" && pp.Name != "" {
					if pp.ID == r {
						presetPermission[pp.Name][pp.Action] = r
					}
				}
			}
		}
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	etd := s.getEnforceTemplateData(ctx)
	form := rolePermissionForm{
		RoleID:               rid,
		CSRFField:            csrf.TemplateField(r),
		ListPermission:       pl,
		ListRoles:            res,
		RolePresetPermission: presetPermission,
		PresetPermission:     etd.PresetPermission,
		ServiceRequest:       etd.ServiceRequests,
		LoginUserInfo:        &usrInfo.UserInfo,
	}

	if hasError {
		form.ErrorMsg = "You have to checked at least one permission"
	}

	form.LoginUserInfo.ProfileImage = usrInfo.ProfileImage

	if err := template.Execute(w, form); err != nil {
		log.Infof("error with template execution: %+v", err)
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func (s *Server) postRolePermission(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	if err := r.ParseForm(); err != nil {
		errMsg := "parsing form"
		log.WithError(err).Error(errMsg)
		http.Error(w, errMsg, http.StatusBadRequest)
		return
	}

	var form rolePermissionForm
	var permissionList []string
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	res, err := s.rbac.ListPermission(ctx, &ppb.ListPermissionRequest{
		OrgID: mw.GetOrgID(ctx),
	})
	if err != nil {
		log.Error("list permissions")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if form.DSAListAndDetails != nil {
		for _, d := range form.DSAListAndDetails {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["dsaListAndDetails"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.TransactionListAndDetails != nil {
		for _, d := range form.TransactionListAndDetails {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["transactionListAndDetails"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.DocumentChecklist != nil {
		for _, d := range form.DocumentChecklist {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["documentChecklist"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.FeeAndCommission != nil {
		for _, d := range form.FeeAndCommission {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["feeAndCommission"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.RiskAssessment != nil {
		for _, d := range form.RiskAssessment {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["riskAssessment"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.Currency != nil {
		for _, d := range form.Currency {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["currency"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.LocationAndBranches != nil {
		for _, d := range form.LocationAndBranches {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["locationAndBranches"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if form.ServiceCatalog != nil {
		for _, d := range form.ServiceCatalog {
			for _, p := range res.GetPermissions() {
				if p.Resource == pm.Permissions["serviceCatalog"].Resource &&
					p.Action == d {
					permissionList = append(permissionList, p.GetID())
				}
			}
		}
	}

	if len(permissionList) < 1 {
		s.permissionForm(w, r, true)
		return
	}

	urp := UpdateRolePermissions{
		ID:          form.RoleID,
		Permissions: permissionList,
	}
	d, err := json.Marshal(urp)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	me := &mfaEvent{
		resource: string(storage.RolePermission),
		action:   tpb.ActionType_Update,
		data:     d,
	}
	if err := s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if err := s.updateRolePermissions(ctx, log, urp); err != nil {
			log.Error("updating role permission")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
	}

	http.Redirect(w, r, "/dashboard/manage-role/edit/"+form.RoleID+"?show_otp=true", http.StatusSeeOther)
}

func (s *Server) updateRolePermissions(ctx context.Context, log *logrus.Entry, urp UpdateRolePermissions) error {
	res, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{ID: []string{urp.ID}})
	if err != nil {
		logging.WithError(err, log).Error("getting role")
		return err
	}
	if len(res.GetRoles()) < 1 {
		log.Error("listing role")
		return fmt.Errorf("listing role")
	}
	r := res.GetRoles()[0]

	for _, p := range urp.Permissions {
		_, err := s.rbac.AssignRolePermission(ctx, &ppb.AssignRolePermissionRequest{
			RoleID:     urp.ID,
			Permission: p,
		})
		if err != nil {
			logging.WithError(err, log).Error("assigning permissions to role")
			return err
		}
	}

	if _, err := s.rbac.UpdateRole(ctx, &ppb.UpdateRoleRequest{
		ID:          urp.ID,
		Name:        r.GetName(),
		Description: r.GetDescription(),
	}); err != nil {
		logging.WithError(err, log).Error("UpdateRole failed")
		return err
	}

	revps := getRevokePermissions(urp.Permissions, r.GetPermissions())
	for _, p := range revps {
		_, err := s.rbac.RevokeRolePermission(ctx, &ppb.RevokeRolePermissionRequest{
			RoleID:     urp.ID,
			Permission: p,
		})
		if err != nil {
			logging.WithError(err, log).Error("revoking permissions from role")
			return err
		}
	}
	return nil
}

func getRevokePermissions(nps, ops []string) []string {
	revps := []string{}
	for _, op := range ops {
		if !contains(nps, op) {
			revps = append(revps, op)
		}
	}
	return revps
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type UpdateRolePermissions struct {
	ID          string
	Permissions []string
}
