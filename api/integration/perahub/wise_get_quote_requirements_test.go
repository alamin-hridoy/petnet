package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISEGetQuoteRequirements(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)

	rs := []WISERequirement{}
	rs = append(rs, WISERequirement{PropName: "accountNumber", Value: "28821822"})
	rs = append(rs, WISERequirement{PropName: "address", Value: map[string]string{"country": "GB", "city": "London", "firstLine": "10 Downing Street", "postCode": "SW1A 2AA"}})

	wantReq := WISEGetQuoteRequirementsReq{
		SourceCurrency: "PHP",
		TargetCurrency: "GBP",
		SourceAmount:   "1500",
	}
	wantResp := &WISEGetQuoteRequirementsResp{
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
		Quote: WISEQuote{
			SourceCurrency: "PHP",
			TargetCurrency: "GBP",
			SourceAmount:   "1500",
		},
	}
	gotResp, err := s.WISEGetQuoteRequirements(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	var gotReq WISEGetQuoteRequirementsReq
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
