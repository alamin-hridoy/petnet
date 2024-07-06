package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestIESend(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := IESendRequest{
		LocationID:         "192",
		LocationName:       "MALOLOSS",
		UserID:             "1894",
		TrxDate:            "2021-07-14",
		CustomerID:         "7712781",
		CurrencyID:         "1",
		RemcoID:            "25",
		TrxType:            "1",
		IsDomestic:         "1",
		CustomerName:       "Cortez, Fernandos",
		ServiceCharge:      "0",
		RemoteLocationID:   "192",
		DstAmount:          "0",
		TotalAmount:        "3301",
		RemoteUserID:       "1894",
		RemoteIPAddress:    "130.211.2.204",
		PurposeTransaction: "Family Support/Living Expenses",
		SourceFund:         "Salary/Income",
		Occupation:         "Unemployed",
		RelationTo:         "Family",
		BirthDate:          "2000-12-16",
		BirthPlace:         "MALOLOS,BULACANS",
		BirthCountry:       "Philippines",
		IDType:             "Postal ID",
		IDNumber:           "PRND32200265569Q",
		Address:            "18 SITIO PULO",
		Barangay:           "BULIHAN",
		City:               "MALOLOS",
		Province:           "BULACAN",
		ZipCode:            "3000B",
		Country:            "Philippines",
		ContactNumber:      "09265454936",
		RiskScore:          "1",
		PayoutType:         "1",
		SenderName:         "Mercado, Marites Cuetos",
		ReceiverName:       "Cortez, Fernandos",
		PrincipalAmount:    "3301",
		ClientReferenceNo:  "ccea7bc90e207c0016e7",
		ControlNumber:      "78K11786",
		McRate:             "1",
		McRateID:           "1",
		RateCategory:       "required333",
		OriginatingCountry: "Philippines",
		DestinationCountry: "Philippines",
		FormType:           "0",
		FormNumber:         "0",
		IPAddress:          "::1",
		ReferenceNumber:    "0",
		BuyBackAmount:      "0",
		ReceiverIDNumber:   "000010102",
		ReceiverPhone:      "09999999998",
		ReceiverAddress: ReceiverAddress{
			Address1: "Marcos Highway",
			Address2: "null",
			Barangay: "Mayamot",
			City:     "ERMITA",
			Province: "MANILA METROPOLITAN",
			ZipCode:  "1000",
			Country:  "PH",
		},
	}

	tests := []struct {
		name        string
		in          IESendRequest
		expectedReq IESendRequest
		want        *IESendResponse
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &IESendResponse{
				Code:    "200",
				Message: "Good",
				Result: IESendResult{
					ID:                 "7288",
					LocationID:         "192",
					UserID:             "1894",
					TrxDate:            "2021-07-14",
					CurrencyID:         "1",
					RemcoID:            "25",
					TrxType:            "1",
					IsDomestic:         "1",
					CustomerID:         "7712781",
					CustomerName:       "Cortez, Fernandos",
					ControlNumber:      "CTRL1",
					SenderName:         "John, Michael Doe",
					ReceiverName:       "Jane, Emily Doe",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					DstAmount:          "0.00",
					TotalAmount:        "1001.00",
					McRate:             "1.00",
					BuyBackAmount:      "0.00",
					RateCategory:       "required",
					McRateID:           "1",
					OriginatingCountry: "Philippines",
					DestinationCountry: "Philippines",
					PurposeTransaction: "Family Support/Living Expenses",
					SourceFund:         "Salary/Income",
					Occupation:         "Unemployed",
					RelationTo:         "Family",
					BirthDate:          "2000-12-16",
					BirthPlace:         "MALOLOS,BULACANS",
					BirthCountry:       "Philippines",
					IDType:             "Postal ID",
					IDNumber:           "PRND32200265569Q",
					Address:            "18 SITIO PULO",
					Barangay:           "BULIHAN",
					City:               "MALOLOS",
					Province:           "BULIHAN",
					Country:            "Philippines",
					ContactNumber:      "09265454936",
					CurrentAddress:     "null",
					PermanentAddress:   "null",
					RiskScore:          "1",
					RiskCriteria:       "1",
					ClientReferenceNo:  "REF1",
					FormType:           "0",
					FormNumber:         "0",
					PayoutType:         "1",
					RemoteLocationID:   "192",
					RemoteUserID:       "1894",
					RemoteIPAddress:    "130.211.2.204",
					IPAddress:          "::1",
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
					ReferenceNumber:    "1",
					ZipCode:            "3000B",
					Status:             "1",
					APIRequest:         "null",
					SapForm:            "null",
					SapFormNumber:      "null",
					SapValidID1:        "null",
					SapValidID2:        "null",
					SapOboLastName:     "null",
					SapOboFirstName:    "null",
					SapOboMiddleName:   "null",
				},
				RemcoID: "24",
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, m := newTestSvc(t, st)
			got, err := s.IESend(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("IESend() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq IESendRequest
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}
			tOps := []cmp.Option{
				cmpopts.IgnoreFields(IESendResult{}, "CreatedAt", "UpdatedAt"),
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
