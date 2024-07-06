package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"brank.as/petnet/api/core/static"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FIRequest struct {
	FrgnRefNo       string          `json:"foreign_reference_no"`
	TransactionType static.WUTxType `json:"transaction_type"`
	PrincipalAmount json.Number     `json:"principal_amount"`
	FixedAmountFlag string          `json:"fixed_amount_flag"`
	DestCountry     string          `json:"destination_country_code"`
	DestCurrency    string          `json:"destination_currency_code"`
	PromoCode       string          `json:"promo_code"`
	Message         []string        `json:"message"`
	MessageLen      string          `json:"message_line_count"`
	TerminalID      string          `json:"terminal_id"`
	OperatorID      string          `json:"operator_id"`
	UserCode        string          `json:"-"`
}

type FIResponseBody struct {
	Taxes         FITaxes     `json:"taxes"`
	OrigPrincipal json.Number `json:"originators_principal_amount"`
	OrigCurrency  string      `json:"originating_currency_principal"`
	DestPrincipal json.Number `json:"destination_principal_amount"`

	ExchangeRate json.Number `json:"exchange_rate"`
	GrossTotal   json.Number `json:"gross_total_amount"`
	PlusCharges  json.Number `json:"plus_charges_amount"`
	PayAmount    json.Number `json:"pay_amount"`
	Charges      json.Number `json:"charges"`
	Tolls        json.Number `json:"tolls"`
	CNDExgFee    json.Number `json:"canadian_dollar_exchange_fee"`

	MessageCharge   json.Number `json:"message_charge"`
	PromoCodeDesc   string      `json:"promo_code_description"`
	PromoSequenceNo string      `json:"promo_sequence_no"`
	PromoName       string      `json:"promo_name"`

	PromoDiscountAmt json.Number `json:"promo_discount_amount"`
	BaseMsgCharge    string      `json:"base_message_charge"`
	BaseMsgLimit     string      `json:"base_message_limit"`
	IncMsgCharge     string      `json:"incremental_message_charge"`
	IncMsgLimit      string      `json:"incremental_message_limit"`
}

type FITaxes struct {
	TaxRate      json.Number `json:"tax_rate"`
	MuniTax      json.Number `json:"municipal_tax"`
	StateTax     json.Number `json:"state_tax"`
	CountyTax    json.Number `json:"county_tax"`
	TaxWorksheet string      `json:"tax_worksheet"`
}

func (s *Svc) FeeInquiry(ctx context.Context, fr FIRequest) (*FIResponseBody, error) {
	if err := validation.ValidateStruct(&fr,
		validation.Field(&fr.TransactionType, validation.Required,
			validation.By(func(interface{}) error {
				switch fr.TransactionType {
				case static.WUSendMoney, static.WUDirectBank, static.WUMobileTransfer:
					return nil
				}
				return fmt.Errorf("invalid transaction type")
			})),
	); err != nil {
		return nil, err
	}
	ln := len(fr.Message)
	if ln == 1 && len(fr.Message[0]) == 0 {
		// empty, non-nil, slice is zero length
		ln = 0
		fr.Message = fr.Message[:0]
	}
	fr.MessageLen = strconv.Itoa(ln)
	const mod, modReq = "wuso", "feeinquiry"
	req, err := s.newParahubRequest(ctx,
		mod, modReq, fr,
		WithUserCode(json.Number(fr.UserCode)),
		WithLocationCode(fr.UserCode))
	if err != nil {
		return nil, err
	}

	resp, err := s.post(ctx, s.moduleURL(mod, modReq), *req)
	if err != nil {
		return nil, err
	}

	var fiRes FIResponseBody
	if err := json.Unmarshal(resp, &fiRes); err != nil {
		return nil, err
	}

	return &fiRes, nil
}
