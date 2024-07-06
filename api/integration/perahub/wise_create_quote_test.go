package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISECreateQuote(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)

	rs := []WISERequirement{}
	rs = append(rs, WISERequirement{PropName: "accountNumber", Value: "28821822"})
	rs = append(rs, WISERequirement{PropName: "address", Value: map[string]string{"country": "GB", "city": "London", "firstLine": "10 Downing Street", "postCode": "SW1A 2AA"}})

	wantReq := WISECreateQuoteReq{
		Email:          "sender2@brankas.com",
		SourceCurrency: "PHP",
		TargetCurrency: "GBP",
		SourceAmount:   "1500",
	}
	wantResp := &WISECreateQuoteResp{
		Requirements: []WISERequirementsResp{
			{
				Type:      "sort_code",
				Title:     "Local bank account",
				UsageInfo: "usage",
				Fields: []WISEField{
					{
						Name: "Recipient type",
						Group: []WISEGroup{
							{
								Key:                "legalType",
								Name:               "Recipient type",
								Type:               "select",
								RefreshReqOnChange: true,
								Required:           true,
								Example:            "example",
								MinLength:          "1",
								MaxLength:          "2",
								ValidationAsync: WISEValidationAsync{
									URL: "url",
									Params: []WISEParam{
										{
											Key:       "key",
											ParamName: "paramname",
											Required:  true,
										},
									},
								},
								ValuesAllowed: []WISEValueAllowed{
									{
										Key:  "PRIVATE",
										Name: "Person",
									},
									{
										Key:  "BUSINESS",
										Name: "Business",
									},
								},
							},
						},
					},
				},
			},
		},
		QuoteSummary: WISEQuoteInquiryResp{
			SourceCurrency: "PHP",
			TargetCurrency: "GBP",
			SourceAmount:   "1500",
			TargetAmount:   "21.1",
			FeeBreakdown: WISEFeeBreakdown{
				Transferwise: "79.66",
				PayIn:        "0",
				Discount:     "0",
				Total:        "79.66",
				PriceSetID:   "132",
				Partner:      "0",
			},
			TotalFee:       "79.66",
			TransferAmount: "1420.34",
			PayOut:         "BANK_TRANSFER",
			Rate:           "0.0148552",
		},
	}
	gotResp, err := s.WISECreateQuote(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISECreateQuoteReq
	if err := json.Unmarshal(m.GetMockRequest(), &gotReq); err != nil {
		t.Fatal(err)
	}

	if !cmp.Equal(wantReq, gotReq) {
		t.Error(cmp.Diff(wantReq, gotReq))
	}
	if !cmp.Equal(wantResp, gotResp) {
		t.Error(cmp.Diff(wantResp, gotResp))
	}
}
