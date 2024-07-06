package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestJPRPayout(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	req := JPRPayoutRequest{
		LocationID:         "371",
		UserID:             4092,
		TrxDate:            "2021-10-12",
		CurrencyID:         "1",
		RemcoID:            "16",
		TrxType:            "2",
		IsDomestic:         "0",
		CustomerID:         "7085005",
		CustomerName:       "TEST,MERLINDA",
		ControlNumber:      "800118975",
		SenderName:         "TEST,NOMAN",
		ReceiverName:       "TEST,MERLINDA",
		PrincipalAmount:    "4959.37",
		ServiceCharge:      "0",
		DstAmount:          "0",
		TotalAmount:        "4959.37",
		McRate:             "0",
		BuyBackAmount:      "0",
		McRateId:           "0",
		OriginatingCountry: "Japan",
		DestinationCountry: "PH",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Salary/Income",
		Occupation:         "IT and Tech Professional",
		RelationTo:         "Family",
		BirthDate:          "2020-07-05",
		BirthPlace:         "ERMITA,MANILA",
		BirthCountry:       "Philippines",
		IDType:             "Voter's ID",
		IDNumber:           "01291904",
		Address:            "#32 Griffin",
		Barangay:           "Pinagbuhatan",
		City:               "ERMITA",
		Province:           "MANILA",
		Country:            "Philippines",
		ContactNumber:      "09187568452",
		CurrentAddress: NonexAddress{
			Address1: "#32 Griffin",
			Address2: "",
			Barangay: "Pinagbuhatan",
			City:     "ERMITA",
			Province: "MANILA",
			ZipCode:  "1000A",
			Country:  "Philippines",
		},
		PermanentAddress: NonexAddress{
			Address1: "#32 Griffin",
			Address2: "",
			Barangay: "Pinagbuhatan",
			City:     "ERMITA",
			Province: "MANILA",
			ZipCode:  "1000A",
			Country:  "Philippines",
		},
		RiskScore:         "1",
		RiskCriteria:      "{}",
		ClientReferenceNo: "9dd5e7289db5cf334893",
		FormType:          "OAR",
		FormNumber:        "0",
		PayoutType:        "1",
		RemoteLocationID:  "371",
		RemoteUserID:      4092,
		RemoteIPAddress:   "::1",
		IPAddress:         "::1",
		ReferenceNumber:   "20211224MACW7QZIEO",
		ZipCode:           "1000A",
		DsaCode:           "TEST_DSA",
		DsaTrxType:        "digital",
	}

	tests := []struct {
		name        string
		in          JPRPayoutRequest
		expectedReq JPRPayoutRequest
		want        *JPRPayoutResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &JPRPayoutResponseBody{
				Code:    "1",
				Message: "Success",
				Result: JPRPayoutResult{
					ControlNumber: "CTRL1",
					Status:        "paid",
				},
				RemcoID: "17",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.JPRPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("JPRPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq JPRPayoutRequest
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
