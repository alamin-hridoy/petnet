package perahub

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestIEPayout(t *testing.T) {
	t.Parallel()
	st := newTestStorage(t)
	req := IEPayoutRequest{
		LocationID:         "371",
		UserID:             "4092",
		TrxDate:            "2021-10-12",
		CurrencyID:         "1",
		RemcoID:            "24",
		TrxType:            "2",
		IsDomestic:         "0",
		CustomerID:         "7085005",
		CustomerName:       "Bering,Elizbeth",
		ControlNumber:      "15K01231",
		SenderName:         "Hamza,Hanif",
		ReceiverName:       "Bering,Elizbeth",
		PrincipalAmount:    "2340.0000",
		ServiceCharge:      "0",
		DstAmount:          "0",
		TotalAmount:        "2340.0000",
		McRate:             "0",
		RateCategory:       "required",
		BuyBackAmount:      "0",
		McRateID:           "0",
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
		RiskCriteria:      "",
		ClientReferenceNo: "9dd5e7289db5cf334893",
		FormType:          "0",
		FormNumber:        "0",
		PayoutType:        "1",
		RemoteLocationID:  "371",
		RemoteUserID:      "4092",
		RemoteIPAddress:   "::1",
		IPAddress:         "::1",
		ReferenceNumber:   "20211012PGFUDGVTUR",
		ZipCode:           "1000A",
		DsaCode:           "TEST_DSA",
		DsaTrxType:        "digital",
	}

	tests := []struct {
		name        string
		in          IEPayoutRequest
		expectedReq IEPayoutRequest
		want        *IEPayoutResponse
		wantErr     bool
	}{
		{
			name:        "Success",
			in:          req,
			expectedReq: req,
			want: &IEPayoutResponse{
				Code:    "200",
				Message: "Good",
				Result: IEPayoutResult{
					ID:                 "7287",
					LocationID:         "371",
					UserID:             "5500",
					TrxDate:            "2021-06-03",
					CurrencyID:         "1",
					RemcoID:            "24",
					TrxType:            "2",
					IsDomestic:         "1",
					CustomerID:         "6925594",
					CustomerName:       "Levy Robert Sogocio",
					ControlNumber:      "CTRL1",
					SenderName:         "Reyes, Ana,",
					ReceiverName:       "Bayuga, Mary Monica",
					PrincipalAmount:    "1000.00",
					ServiceCharge:      "1.00",
					DstAmount:          "0.00",
					TotalAmount:        "1001.00",
					McRate:             "0.00",
					BuyBackAmount:      "0.00",
					RateCategory:       "0",
					McRateID:           "0",
					OriginatingCountry: "Philippines",
					DestinationCountry: "PH",
					PurposeTransaction: "Family Support/Living Expenses",
					SourceFund:         "Savings",
					Occupation:         "OTH",
					RelationTo:         "Family",
					BirthDate:          "1995-04-14",
					BirthPlace:         "TAGUDIN,ILOCOS SUR",
					BirthCountry:       "PH",
					IDType:             "LICENSE",
					IDNumber:           "B83180608851",
					Address:            "MAIN ST",
					Barangay:           "TALLAOEN",
					City:               "AKLAN CITY",
					Province:           "AKLAN CITY",
					Country:            "PH",
					ContactNumber:      "09516738640",
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
					RiskScore:         "0",
					RiskCriteria:      "0",
					ClientReferenceNo: "7884447474",
					FormType:          "0",
					FormNumber:        "0",
					PayoutType:        "1",
					RemoteLocationID:  "371",
					RemoteUserID:      "5684",
					RemoteIPAddress:   "130.211.2.187",
					IPAddress:         "130.211.2.187",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
					ReferenceNumber:   "848jjfu23u333",
					ZipCode:           "36989",
					Status:            "1",
					APIRequest:        "null",
					SapForm:           "null",
					SapFormNumber:     "null",
					SapValidID1:       "null",
					SapValidID2:       "null",
					SapOboLastName:    "null",
					SapOboFirstName:   "null",
					SapOboMiddleName:  "null",
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
			got, err := s.IEPayout(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Errorf("IEPayout() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			var newReq IEPayoutRequest
			if err := json.Unmarshal(m.GetMockRequest(), &newReq); err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.expectedReq, newReq) {
				t.Error(cmp.Diff(test.expectedReq, newReq))
			}
			tOps := []cmp.Option{
				cmpopts.IgnoreFields(IEPayoutResult{}, "CreatedAt", "UpdatedAt"),
			}
			if !cmp.Equal(test.want, got, tOps...) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
