package handler

import (
	"fmt"
	"net/http"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
)

const (
	DIGITAL = "DIGITAL"
	OTC     = "OTC"
)

func (s *Server) skipOnboarding(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := logging.FromContext(ctx).WithField("handler", "skipOnboarding")

	sess, err := s.sess.Get(r, sessionCookieName)
	if err != nil {
		log.WithError(err).Error("fetching session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	orgId, ok := sess.Values[sessionOrgID].(string)
	if !ok {
		log.Error("missing org id in session")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	_, err = s.pf.UpsertProfile(ctx, &ppb.UpsertProfileRequest{
		Profile: &ppb.OrgProfile{
			OrgID:            orgId,
			TransactionTypes: fmt.Sprintf("%s,%s", DIGITAL, OTC),
		},
	})
	if err != nil {
		log.WithError(err).Error("error upserting profile")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
