package perahub

import (
	"context"
	"encoding/json"
)

type FSRequest struct {
	RefNo                   string   `json:"foreign_reference_no"`
	PrincipalAmount         string   `json:"principal_amount"`
	FixedAmountFlag         string   `json:"fixed_amount_flag"`
	DestinationCountryCode  string   `json:"destination_country_code"`
	DestinationCurrencyCode string   `json:"destination_currency_code"`
	TransactionType         string   `json:"transaction_type"`
	PromoCode               string   `json:"promo_code"`
	Message                 []string `json:"message"`
	MessageLineCount        string   `json:"message_line_count"`
	TerminalID              string   `json:"terminal_id"`
	OperatorID              string   `json:"operator_id"`
}

type FSResponseBody struct {
	StagingBuffer string `json:"staging_buffer"`
}

func (s *Svc) FeeSurvey(ctx context.Context, fr FSRequest) (string, error) {
	const mod, modReq = "prereq", "feesurvey"
	req, err := s.newParahubRequest(ctx, mod, modReq, fr)
	if err != nil {
		return "", err
	}

	resp, err := s.post(ctx, s.moduleURL(mod, ""), *req)
	if err != nil {
		return "", err
	}

	var fsRes FSResponseBody
	if err := json.Unmarshal(resp, &fsRes); err != nil {
		return "", err
	}

	return fsRes.StagingBuffer, nil
}
