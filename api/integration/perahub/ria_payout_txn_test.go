package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRiaPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := RiaPayoutRequest{
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             "1893",
		TrxDate:            "2021-07-13",
		CustomerID:         "7712780",
		CurrencyID:         "1",
		RemcoID:            "12",
		TrxType:            "2",
		IsDomestic:         "7",
		CustomerName:       "cust-name",
		ServiceCharge:      "0",
		RmtLocID:           "9",
		DstAmount:          "0",
		TotalAmount:        "11",
		RmtUserID:          "12",
		RmtIPAddr:          "127.0.0.1",
		PurposeTxn:         "prps",
		SourceFund:         "src",
		Occupation:         "occ",
		RelationTo:         "rel",
		BirthDate:          "birty",
		BirthPlace:         "bplace",
		BirthCountry:       "bcountry",
		IDType:             "id-type",
		IDNumber:           "id-no",
		Address:            "addr",
		Barangay:           "barang",
		City:               "city",
		Province:           "prov",
		ZipCode:            "12345",
		Country:            "country",
		ContactNumber:      "54321",
		RiskScore:          "1",
		PayoutType:         "1",
		SenderName:         "send",
		RcvName:            "rcv",
		PnplAmt:            "14",
		ControlNo:          "ctrl-no",
		RefNo:              "ref-no",
		Natl:               "PH",
		Gender:             "M",
		FormType:           "frm-type",
		FormNumber:         "1234",
		IDIssueDate:        "2021-07-13",
		IDExpDate:          "2021-07-13",
		IDIssBy:            "PH",
		BBAmt:              "1234",
		MCRateID:           "123",
		MCRate:             "123",
		RateCat:            "1",
		OriginatingCountry: "PH",
		DestinationCountry: "PH",
		DeviceID:           "desktop123",
		AgentID:            "84424911",
		AgentCode:          "ITD",
		Currency:           "PHP",
		ClientReferenceNo:  "189528302",
		OrderNo:            "TH1950882455",
		DsaCode:            "TEST_DSA",
		DsaTrxType:         "digital",
	}

	tests := []struct {
		name        string
		in          RiaPayoutRequest
		expectedReq RiaPayoutRequest
		want        *RiaPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &RiaPayoutResponseBody{
				Code:    "200",
				Msg:     "Successful.",
				RemcoID: "12",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			got, err := s.RiaPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("RiaPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq RiaPayoutRequest
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}

			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
