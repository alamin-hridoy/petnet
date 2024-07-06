package handler

import (
	"encoding/json"
	"net/http"

	"brank.as/petnet/cms/storage"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	ppb "brank.as/rbac/gunk/v1/permissions"

	"github.com/kenshaw/goji"
)

func (s *Server) getManageDeleteRole(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	rid := goji.Param(r, "id")
	if rid == "" {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	res, err := s.rbac.ListRole(ctx, &ppb.ListRoleRequest{
		OrgID: mw.GetOrgID(ctx),
		ID:    []string{rid},
	})
	if err != nil {
		logging.WithError(err, log).Error(err.Error())
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	rs := res.GetRoles()
	if len(rs) < 1 {
		log.Error("no roles found")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	// todo(codemen): make error be shown on the frontend when role has members and can't be deleted
	if len(rs[0].GetMembers()) > 0 {
		log.Error("can't delete role which has to be assigned to user")
		http.Redirect(w, r, manageRoleListPath+"?errorMsg=true", http.StatusSeeOther)
		return
	}

	dr := &ppb.DeleteRoleRequest{
		ID: rid,
	}
	d, err := json.Marshal(dr)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	me := &mfaEvent{
		resource: string(storage.RoleDelete),
		action:   tpb.ActionType_Delete,
		data:     d,
	}
	if err := s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		if _, err := s.rbac.DeleteRole(ctx, dr); err != nil {
			log.WithError(err).Error(err.Error())
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, manageRoleListPath, http.StatusSeeOther)
	}
	http.Redirect(w, r, manageRoleListPath+"?show_otp=true", http.StatusSeeOther)
}
