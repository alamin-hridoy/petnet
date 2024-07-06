package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISERefreshRecipient(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)

	rs := []WISERequirement{}
	rs = append(rs, WISERequirement{PropName: "accountNumber", Value: "28821822"})
	rs = append(rs, WISERequirement{PropName: "address", Value: map[string]string{"country": "GB", "city": "London", "firstLine": "10 Downing Street", "postCode": "SW1A 2AA"}})

	wantReq := WISERefreshRecipientReq{
		Email:             "sender2@brankas.com",
		Currency:          "GBP",
		Type:              "sort_code",
		OwnedByCustomer:   false,
		AccountHolderName: "Brankas Receiver",
		Requirements:      rs,
	}
	wantResp := &WISERefreshRecipientResp{
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
	}
	gotResp, err := s.WISERefreshRecipient(context.Background(), wantReq)
	if err != nil {
		t.Fatal(err)
	}

	gotReq := m.GetMockRequest()
	if jsonReq != string(gotReq) {
		t.Errorf("not equal, want:\n%v\n, got: \n%v\n", jsonReq, string(gotReq))
	}
	if !cmp.Equal(wantResp, gotResp) {
		t.Error(cmp.Diff(wantResp, gotResp))
	}
}
