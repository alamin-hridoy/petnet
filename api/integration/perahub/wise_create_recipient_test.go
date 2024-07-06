package perahub

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWISECreateRecipient(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	s, m := newTestSvc(t, st)

	rs := []WISERequirement{}
	rs = append(rs, WISERequirement{PropName: "accountNumber", Value: "28821822"})
	rs = append(rs, WISERequirement{PropName: "address", Value: map[string]string{"country": "GB", "city": "London", "firstLine": "10 Downing Street", "postCode": "SW1A 2AA"}})

	wantReq := WISECreateRecipientReq{
		Email:             "sender2@brankas.com",
		Currency:          "GBP",
		Type:              "sort_code",
		OwnedByCustomer:   false,
		AccountHolderName: "Brankas Receiver",
		Requirements:      rs,
	}
	wantResp := &WISECreateRecipientResp{
		RecipientID: "148142936",
		Details: WISECRDetails{
			Address: WISECRAddress{
				Country:     "GB",
				CountryCode: "GB",
				FirstLine:   "10 Downing Street",
				PostCode:    "SW1A 2AA",
				City:        "London",
			},
			LegalType:     "PRIVATE",
			AccountNumber: "28821822",
			SortCode:      "231470",
		},
		AccountHolderName: "Brankas Receiver",
		Currency:          "GBP",
		OwnedByCustomer:   false,
		Country:           "GB",
		Msg:               "Success! Recipient Account Created",
	}
	gotResp, err := s.WISECreateRecipient(context.Background(), wantReq)
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

var jsonReq = `{"email":"sender2@brankas.com","currency":"GBP","type":"sort_code","ownedByCustomer":false,"accountHolderName":"Brankas Receiver","requirements":[{"propName":"accountNumber","value":"28821822"},{"propName":"address","value":{"city":"London","country":"GB","firstLine":"10 Downing Street","postCode":"SW1A 2AA"}}]}`
