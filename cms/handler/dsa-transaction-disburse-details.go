package handler

import (
	"net/http"

	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
)

func (s *Server) getDSATransactionDisburseDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	refID := goji.Param(r, "id")
	if refID == "" {
		log.Error("missing ref id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("dsa-transaction-disburse-detail.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	oid := mw.GetOrgID(ctx)

	md := metautils.ExtractIncoming(ctx)
	ctx = md.Add("x-forward-dsaorgid", oid).ToOutgoing(ctx)
	res, err := s.drpSB.ListRemit(ctx, &tpb.ListRemitRequest{
		ControlNumbers: []string{refID},
	})
	if err != nil {
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var tk int
	for rk, rd := range res.GetRemittances() {
		if rd.GetRemitType() == "DISBURSE" {
			tk = rk
		}
	}
	details := DSATransactionData(res.GetRemittances()[tk])
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.ID = refID
	if len(res.GetRemittances()) > 1 {
		details.CreateRemit = false
		details.LinkedTrans = true
	}
	if details.RemitDetails.RemitType == "SENT" {
		details.CreateRemit = true
	}
	etd := s.getEnforceTemplateData(ctx)
	details.PresetPermission = etd.PresetPermission
	details.ServiceRequest = etd.ServiceRequests
	if err := template.Execute(w, details); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
