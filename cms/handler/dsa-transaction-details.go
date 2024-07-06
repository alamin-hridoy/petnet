package handler

import (
	"net/http"

	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/mw"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/kenshaw/goji"
)

type (
	DSATransactionDetails struct {
		ID               string
		UserInfo         *User
		RemitDetails     *tpb.Remittance
		PresetPermission map[string]map[string]bool
		ServiceRequest   bool
		CreateRemit      bool
		LinkedTrans      bool
	}
)

func (s *Server) getDSATransactionDetails(w http.ResponseWriter, r *http.Request) {
	log := logging.FromContext(r.Context())
	ctx := r.Context()

	refID := goji.Param(r, "id")
	if refID == "" {
		log.Error("missing ref id query param")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	template := s.templates.Lookup("dsa-transaction-detail.html")
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
		log.Error("unable to connect api")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}

	var tk int
	for rk, rd := range res.GetRemittances() {
		if rd.GetRemitType() == "SEND" {
			tk = rk
		}
	}
	details := DSATransactionData(res.GetRemittances()[tk])
	details.CreateRemit = false
	details.LinkedTrans = false
	if len(res.GetRemittances()) > 1 {
		details.LinkedTrans = true
	}
	if details.RemitDetails.RemitType == "SEND" {
		details.CreateRemit = true
	}
	usrInfo := s.GetUserInfoFromCookie(w, r, false)

	details.UserInfo = &usrInfo.UserInfo
	details.ID = refID
	etd := s.getEnforceTemplateData(ctx)
	details.PresetPermission = etd.PresetPermission
	details.ServiceRequest = etd.ServiceRequests
	if err := template.Execute(w, details); err != nil {
		logging.WithError(err, log).Error("error with template execution")
		http.Redirect(w, r, errorPath, http.StatusSeeOther)
		return
	}
}

func DSATransactionData(res *tpb.Remittance) DSATransactionDetails {
	partner := res.RemitPartner
	switch partner {
	case "WU":
		partner = "Western Union"
	case "IR":
		partner = "IRemit"
	case "TF":
		partner = "Transfast"
	}

	status := res.RemitType
	switch status {
	case "CREATE":
		status = "Create-Remit"
	case "DISBURSE":
		status = "Disburse-Remit"
	}
	return DSATransactionDetails{
		RemitDetails: &tpb.Remittance{
			ControlNumber:            res.GetControlNumber(),
			RemitType:                res.GetRemitType(),
			GrossAmount:              res.GetGrossAmount(),
			RemitAmount:              res.GetRemitAmount(),
			Remitter:                 res.GetRemitter(),
			Receiver:                 res.GetReceiver(),
			RemitPartner:             partner,
			TransactionStagedTime:    res.TransactionStagedTime,
			TransactionCompletedTime: res.TransactionCompletedTime,
		},
	}
}
