package handler

import (
	"net/http"

	upb "brank.as/petnet/gunk/dsa/v1/user"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/kenshaw/goji"
)

func (s *Server) getManageEnableUser(w http.ResponseWriter, r *http.Request) {
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
	if _, err := s.pf.EnableUserProfile(ctx, &upb.EnableUserProfileRequest{
		UserID: userid,
	}); err != nil {
		log.WithError(err).Error(err.Error())
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, manageUserListPath, http.StatusSeeOther)
}
