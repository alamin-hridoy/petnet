package handler

import (
	"encoding/json"
	"net/http"

	"brank.as/petnet/cms/storage"
	upb "brank.as/petnet/gunk/dsa/v1/user"
	tpb "brank.as/petnet/gunk/dsa/v2/temp"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	rbipb "brank.as/rbac/gunk/v1/invite"
	rbupb "brank.as/rbac/gunk/v1/user"
	"github.com/kenshaw/goji"

	pps "brank.as/petnet/profile/storage"
)

func (s *Server) getManageDisableUser(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	userid := goji.Param(r, "id")
	cuid := mw.GetUserID(ctx)
	if userid == "" {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	if cuid == userid {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	du := &rbupb.DisableUserRequest{
		UserID: userid,
	}

	ireq := &rbipb.CancelInviteRequest{
		ID: userid,
	}

	d, err := json.Marshal(du)
	if err != nil {
		logging.WithError(err, log).Error("marshal request")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	me := &mfaEvent{
		resource: string(storage.UserDelete),
		action:   tpb.ActionType_Delete,
		data:     d,
	}
	if err := s.initMFAEvent(w, r, me); err != nil {
		if err != storage.MFANotFound {
			logging.WithError(err, log).Error("initializing mfa event")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		ul, er := s.rbac.GetUser(ctx, &rbupb.GetUserRequest{
			ID: userid,
		})
		if er != nil {
			logging.WithError(err, log).Error("getting user failed.")
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
		}
		if ul.User.GetInviteStatus() == pps.InviteSent {
			if _, err := s.rbac.CancelInvite(ctx, ireq); err != nil {
				log.WithError(err).Error(err.Error())
				http.Redirect(w, r, errorPath, http.StatusSeeOther)
				return
			}
		}
		if _, err := s.pf.DeleteUserProfile(ctx, &upb.DeleteUserProfileRequest{
			UserID: du.UserID,
		}); err != nil {
			log.WithError(err).Error(err.Error())
			http.Redirect(w, r, errorPath, http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
	}
	http.Redirect(w, r, manageUserListPath+"?show_otp=true", http.StatusSeeOther)
}
