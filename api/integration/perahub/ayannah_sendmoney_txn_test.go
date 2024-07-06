package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestAYANNAHSendMoney(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := AYANNAHSendRequest{
		LocationID:         "191",
		LocationName:       "MALOLOS",
		UserID:             "1893",
		TrxDate:            "2021-07-13",
		CustomerID:         "7712780",
		CurrencyID:         "1",
		RemcoID:            "22",
		TrxType:            "1",
		IsDomestic:         "1",
		CustomerName:       "Cortez, Fernando",
		ServiceCharge:      "1",
		RemoteLocationID:   "191",
		DstAmount:          "0",
		TotalAmount:        "1001.00",
		RemoteUserID:       "1893",
		RemoteIPAddress:    "130.211.2.203",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Salary/Income",
		Occupation:         "Unemployed",
		RelationTo:         "Family",
		BirthDate:          "2000-12-15",
		BirthPlace:         "MALOLOS,BULACAN",
		BirthCountry:       "Philippines",
		IDType:             "PASSPORT",
		IDNumber:           "PRND32200265569P",
		Address:            "18 SITIO PULO",
		Barangay:           "BULIHAN",
		City:               "MALOLOS",
		Province:           "BULACAN",
		ZipCode:            "3000A",
		Country:            "PH",
		ContactNumber:      "09265454935",
		RiskScore:          "1",
		PayoutType:         "1",
		SenderName:         "Mercado, Marites Cueto",
		ReceiverName:       "Cortez, Fernando",
		PrincipalAmount:    "1000.00",
		ClientReferenceNo:  "1",
		ControlNumber:      "CTRL1",
		McRate:             "1",
		McRateID:           "1",
		RateCategory:       "required",
		OriginatingCountry: "Philippines",
		DestinationCountry: "Philippines",
		FormType:           "0",
		FormNumber:         "0",
		IPAddress:          "::1",
		ReferenceNumber:    "REF1",
		BuyBackAmount:      "1",
	}

	tests := []struct {
		name        string
		in          AYANNAHSendRequest
		expectedReq AYANNAHSendRequest
		want        *AYANNAHSendResponseBody
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &AYANNAHSendResponseBody{
				Code:    "200",
				Message: "Success",
				Result: AYANNAHSendResult{
					Message:            "Successfully Sendout.",
					ID:                 "7048",
					LocationID:         "191",
					UserID:             "1893",
					TrxDate:            "2021-07-13",
					CurrencyID:         "1",
					RemcoID:            "22",
					TrxType:            "1",
					IsDomestic:         "1",
					CustomerID:         "7712780",
					CustomerName:       "Cortez, Fernando",
					ControlNumber:      "CTRL1",
					SenderName:         "Mercado, Marites Cueto",
					ReceiverName:       "Cortez, Fernando",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					DstAmount:          "0.00",
					TotalAmount:        "1001.00",
					McRate:             "1.00",
					BuyBackAmount:      "1.00",
					RateCategory:       "required",
					McRateID:           "1",
					OriginatingCountry: "Philippines",
					DestinationCountry: "Philippines",
					PurposeTransaction: "Family Support/Living Expenses",
					SourceFund:         "Salary/Income",
					Occupation:         "Unemployed",
					RelationTo:         "Family",
					BirthDate:          "2000-12-15",
					BirthPlace:         "MALOLOS,BULACAN",
					BirthCountry:       "Philippines",
					IDType:             "PASSPORT",
					IDNumber:           "PRND32200265569P",
					Address:            "18 SITIO PULO",
					Barangay:           "BULIHAN",
					City:               "MALOLOS",
					Province:           "BULACAN",
					Country:            "PH",
					ContactNumber:      "09265454935",
					CurrentAddress: NonexAddress{
						Address1: "Marcos Highway",
						Address2: "null",
						Barangay: "Mayamot",
						City:     "ERMITA",
						Province: "MANILA METROPOLITAN",
						ZipCode:  "1000",
						Country:  "PH",
					},
					PermanentAddress: NonexAddress{
						Address1: "Marcos Highway",
						Address2: "null",
						Barangay: "Mayamot",
						City:     "ERMITA",
						Province: "MANILA METROPOLITAN",
						ZipCode:  "1000",
						Country:  "PH",
					},
					RiskScore:         "1",
					RiskCriteria:      "0",
					ClientReferenceNo: "1",
					FormType:          "0",
					FormNumber:        "0",
					PayoutType:        "1",
					RemoteLocationID:  "191",
					RemoteUserID:      "1893",
					RemoteIPAddress:   "130.211.2.203",
					IPAddress:         "::1",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
					ReferenceNumber:   "REF1",
					ZipCode:           "3000A",
					Status:            "1",
					APIRequest:        "null",
					SapForm:           "null",
					SapFormNumber:     "null",
					SapValidID1:       "null",
					SapValidID2:       "null",
					SapOboLastName:    "null",
					SapOboFirstName:   "null",
					SapOboMiddleName:  "null",
					AyannahStatus:     "NEW",
				},
				RemcoID: "22",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.AYANNAHSendMoney(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("AYANNAHSendMoney() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq AYANNAHSendRequest
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}
			tOps := []cmp.Option{
				cmpopts.IgnoreFields(AYANNAHSendResult{}, "CreatedAt", "UpdatedAt"),
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
