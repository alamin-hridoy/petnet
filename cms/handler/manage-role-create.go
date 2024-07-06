package handler

import (
	"encoding/json"
	"html/template"
	"net/http"

	"brank.as/petnet/cms/storage"
	rbupb "brank.as/rbac/gunk/v1/user"

	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	ppb "brank.as/rbac/gunk/v1/permissions"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/gorilla/csrf"
)

type (
	RoleCreateForm struct {
		Name        string
		Description string
	}

	CreateRoleListTempData struct {
		CSRFField        template.HTML
		Roles            []RoleDetails
		Errors           map[string]string
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
		LoginUserInfo    *User
	}
)

func (s *Server) postManageRoleCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx)

	var form RoleCreateForm
	if err := s.decoder.Decode(&form, r.PostForm); err != nil {
		logging.WithError(err, log).Error("decoding form")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	formErr := map[string]string{}
	if err := validation.ValidateStruct(&form,
		validation.Field(&form.Name, validation.Required),
		validation.Field(&form.Description, validation.Required),
	); err != nil {
		if err, ok := (err).(validation.Errors); ok {
			if err["Name"] != nil {
				formErr["Name"] = err["Name"].Error()
			}
			if err["Description"] != nil {
				formErr["Description"] = err["Description"].Error()
			}
		}
	}
	lr, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{
		Name: form.Name,
	})
	if err != nil {
		log.Error("unable to connect api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	if len(lr.GetRoles()) > 0 {
		formErr["Name"] = "Role Already Exists. Please try another Name"
	}

	if len(formErr) > 0 {
		res, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{
			OrgID: mw.GetOrgID(ctx),
		})
		if err != nil {
			log.Error("unable to connect api")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		var roleList []RoleDetails
		for _, role := range res.GetRoles() {
			CreatedBy := ""
			UpdatedBy := ""
			cu, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: role.CreateUID})
			if err == nil {
				CreatedBy = cu.User.FirstName + " " + cu.User.LastName
			}

			// todo(codemen): &rbupb.GetUserRequest{ID: role.UpdateUID}
			du, err := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{ID: role.CreateUID})
			if err == nil {
				UpdatedBy = du.User.FirstName + " " + du.User.LastName
			}

			// todo(codemen): Updated:     role.Updated,
			roleList = append(roleList, RoleDetails{
				ID:          role.ID,
				Name:        role.Name,
				Description: role.Description,
				CreatedBy:   CreatedBy,
				UpdatedBy:   UpdatedBy,
				Updated:     role.GetCreated().AsTime(),
			})
		}
		etd := s.getEnforceTemplateData(ctx)
		usrInfo := s.GetUserInfoFromCookie(w, r, false)

		tempData := CreateRoleListTempData{
			CSRFField:        csrf.TemplateField(r),
			Roles:            roleList,
			PresetPermission: etd.PresetPermission,
			ServiceRequest:   etd.ServiceRequests,
			LoginUserInfo:    &usrInfo.UserInfo,
		}
		tempData.Errors = formErr
		tempData.CSRFField = csrf.TemplateField(r)
		template := s.templates.Lookup("role-list.html")
		if template == nil {
			log.Error("unable to load template")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if err := template.Execute(w, tempData); err != nil {
			log.Infof("error with template execution: %+v", err)
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		return
	}

	cr := &ppb.CreateRoleRequest{
		Name:        form.Name,
		Description: form.Description,
	}
	d, err := json.Marshal(cr)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	me := &mfaEvent{
		resource: string(storage.CreateRole),
		action:   tpb.ActionType_Create,
		data:     d,
	}
	if err := s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if _, err := s.rbac.CreateRole(ctx, cr); err != nil {
			logging.WithError(err, log).Error("Create Role")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}

		http.Redirect(w, r, manageRoleListPath, http.StatusSeeOther)
	}

	http.Redirect(w, r, manageRoleListPath+"?show_otp=true", http.StatusSeeOther)
}
