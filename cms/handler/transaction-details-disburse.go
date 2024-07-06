package handler

import (
	"net/http"

	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
)

func (s *Server) getTransactionDisburseDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()
	id := goji.Param(r, "id")
	if id == "" {
		log.Error("missing id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	oid := goji.Param(r, "oid")
	if oid == "" {
		log.Error("missing org id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	template := s.templates.Lookup("transaction-disburse-detail.html")
	if template == nil {
		log.Error("unable to load template")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	md := metautils.ExtractIncoming(ctx)
	ctx = md.Add("x-forward-dsaorgid", oid).ToOutgoing(ctx)
	res, err := s.drpSB.ListRemit(ctx, &tpb.ListRemitRequest{
		ControlNumbers: []string{id},
	})
	if err != nil {
		logging.WithError(err, log).Error("unable to connect api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
	var tk int
	for rk, rd := range res.GetRemittances() {
		if rd.GetRemitType() == "DISBURSE" {
			tk = rk
		}
	}

	rmtCnt := len(res.GetRemittances())
	if rmtCnt == 0 || rmtCnt <= tk {
		logging.WithError(err, log).Error("don't have any transaction")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	details := transactionData(res.GetRemittances()[tk])
	if rmtCnt > 1 {
		details.CreateRemit = false
		details.LinkedTrans = true
	}
	if details.RemitDetails.RemitType == "SEND" {
		details.CreateRemit = true
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.OrgID = oid
	etd := s.getEnforceTemplateData(ctx)
	details.PresetPermission = etd.PresetPermission
	details.ServiceRequest = etd.ServiceRequests
	if err := template.Execute(w, details); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}
