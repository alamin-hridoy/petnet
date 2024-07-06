package wu

import (
	"context"
	"encoding/json"
	"strings"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

const (
	terminalID = "PH259AMT001A" // petnet will provide this operator id
	operatorID = "drp"          // petnet will provide this terminal id
)

func (s *Svc) FeeInquiry(ctx context.Context, r core.FeeInquiryReq) (map[string]string, error) {
	log := logging.FromContext(ctx)

	fmFlag := "N"
	if r.DestinationAmount {
		fmFlag = "F"
	}

	usrMsg := strings.Split(r.Message, "\\n")
	fee, err := s.ph.FeeInquiry(ctx, perahub.FIRequest{
		FrgnRefNo:       random.InvitationCode(20),
		PrincipalAmount: json.Number(r.PrincipalAmount.Number()),
		FixedAmountFlag: fmFlag,
		DestCountry:     r.DestCountry,
		DestCurrency:    r.DestCurrency,
		TransactionType: static.WUTxType(r.RemitType.Code),
		PromoCode:       r.Promo,
		Message:         usrMsg,
		// todo(robin): enable dynamic terminalid and operatorid once petnet has set it up for
		// sandbox
		// TerminalID:      phmw.GetTerminalID(ctx),
		// OperatorID:      phmw.GetOperatorID(ctx),
		TerminalID: terminalID,
		OperatorID: operatorID,
	})
	if err != nil {
		logging.WithError(err, log).Error("fee inquiry")
		return nil, err
	}

	return map[string]string{
		"tax_rate":                       fee.Taxes.TaxRate.String(),
		"municipal_tax":                  fee.Taxes.MuniTax.String(),
		"state_tax":                      fee.Taxes.StateTax.String(),
		"county_tax":                     fee.Taxes.CountyTax.String(),
		"tax_worksheet":                  fee.Taxes.TaxWorksheet,
		"originators_principal_amount":   fee.OrigPrincipal.String(),
		"originating_currency_principal": fee.OrigCurrency,
		"destination_principal_amount":   fee.DestPrincipal.String(),
		"exchange_rate":                  fee.ExchangeRate.String(),
		"gross_total_amount":             fee.GrossTotal.String(),
		"plus_charges_amount":            fee.PlusCharges.String(),
		"pay_amount":                     fee.PayAmount.String(),
		"charges":                        fee.Charges.String(),
		"tolls":                          fee.Tolls.String(),
		"canadian_dollar_exchange_fee":   fee.CNDExgFee.String(),
		"message_charge":                 fee.MessageCharge.String(),
		"promo_code_description":         fee.PromoCodeDesc,
		"promo_sequence_no":              fee.PromoSequenceNo,
		"promo_name":                     fee.PromoName,
		"promo_discount_amount":          fee.PromoDiscountAmt.String(),
		"base_message_charge":            fee.BaseMsgCharge,
		"base_message_limit":             fee.BaseMsgLimit,
		"incremental_message_charge":     fee.IncMsgCharge,
		"incremental_message_limit":      fee.IncMsgLimit,
	}, nil
}
