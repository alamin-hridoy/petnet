package ussc

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
	"github.com/bojanz/currency"
)

func (s *Svc) FeeInquiry(ctx context.Context, r core.FeeInquiryReq) (map[string]string, error) {
	log := logging.FromContext(ctx)

	fee, err := s.ph.USSCFeeInquiry(ctx, perahub.USSCFeeInquiryRequest{
		// don't change, should be statically empty
		Panalokard: "",
		Amount:     r.PrincipalAmount.Amount.Number(),
		// don't change, should be statically empty
		USSCPromo:  "",
		BranchCode: "branch1", // todo: will change, gotten from petnet
	})
	if err != nil {
		logging.WithError(err, log).Error("fee inquiry")
		return nil, err
	}

	pcplamt, err := currency.NewAmount(fee.Result.PnplAmount, "PHP")
	if err != nil {
		logging.WithError(err, log).Error("creating new prinipal amount")
		return nil, err
	}
	svcchg, err := currency.NewAmount(fee.Result.ServiceCharge, "PHP")
	if err != nil {
		logging.WithError(err, log).Error("creating new service charge")
		return nil, err
	}
	ttlamt, err := currency.NewAmount(fee.Result.TotAmount, "PHP")
	if err != nil {
		logging.WithError(err, log).Error("creating new total amount")
		return nil, err
	}

	return map[string]string{
		"principal_amount": currency.ToMinor(pcplamt.Round()).Number(),
		"service_charge":   currency.ToMinor(svcchg.Round()).Number(),
		"message":          fee.Result.Msg,
		"code":             fee.Result.Code,
		"new_screen":       fee.Result.NewScreen,
		"journal_no":       fee.Result.JournalNo,
		"process_date":     fee.Result.ProcessDate,
		"reference_number": fee.Result.RefNo,
		"total_amount":     currency.ToMinor(ttlamt.Round()).Number(),
		"send_otp":         fee.Result.SendOTP,
	}, nil
}
