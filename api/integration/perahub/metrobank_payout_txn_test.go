package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMBPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := MBPayoutRequest{
		LocationId:         "1",
		LocationName:       "loc-name",
		UserId:             "1",
		TrxDate:            "trx-date",
		CustomerId:         "1",
		CurrencyId:         "1",
		RemcoId:            "1",
		TrxType:            "trx-tp",
		IsDomestic:         "1",
		CustomerName:       "cust-name",
		ServiceCharge:      "8",
		RemoteLocationId:   "1",
		DstAmount:          "10",
		TotalAmount:        "11",
		RemoteUserId:       "1",
		RemoteIpAddress:    "rep-ip-add",
		PurposeTransaction: "pur-trn",
		SourceFund:         "src",
		Occupation:         "occ",
		RelationTo:         "rel",
		BirthDate:          "birty",
		BirthPlace:         "bplace",
		BirthCountry:       "bcountry",
		IdType:             "ip-tp",
		IdNumber:           "ip-num",
		Address:            "addr",
		Barangay:           "barang",
		City:               "city",
		Province:           "prov",
		ZipCode:            "12345",
		Country:            "country",
		ContactNumber:      "54321",
		RiskScore:          "1",
		PayoutType:         "13",
		SenderName:         "sender name",
		ReceiverName:       "rcv-nam",
		PrincipalAmount:    "1600.00",
		ControlNumber:      "ctrl-num",
		ReferenceNumber:    "ref-no",
		McRate:             "mc-rat",
		BuyBackAmount:      "buy-bank-amnt",
		RateCategory:       "rat-cat",
		McRateId:           "mc-rat-id",
		OriginatingCountry: "orgnl-contr",
		DestinationCountry: "dest-contr",
		ClientReferenceNo:  "cl-ref-num",
		FormType:           "f-type",
		FormNumber:         "f-no",
		Currency:           "PHP",
		DsaCode:            "TEST_DSA",
		DsaTrxType:         "digital",
	}

	tests := []struct {
		name        string
		in          MBPayoutRequest
		expectedReq MBPayoutRequest
		want        *MBPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &MBPayoutResponseBody{
				Code: "200",
				Msg:  "Successful.",
				Result: MBPayResult{
					RefNo:           "REF1",
					ClientRefNo:     "1",
					ControlNo:       "CTRL1",
					StatusText:      "0",
					PrincipalAmount: "1000.00",
					RcvName:         "rcv-nam",
					Address:         "addr",
					ReceiptNo:       "",
					ContactNumber:   "201-20211126-000288",
				},
				RemcoID: "8",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.MBPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("MBPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq MBPayoutRequest
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
